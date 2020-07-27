package dao

import (
	"errors"
	"fmt"
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/libs/encode"
	"inla/inla-crawler/libs/str"
	"inla/inla-crawler/tasks/qgxymdmz/model"
	"strings"
	"time"
)

type Service interface {
	SaveGenre(input *model.GenreData) (*model.GenreData, error)
	SaveBook(input *model.GenreItemData) (*model.GenreItemData, error)
	SaveChapter(input *model.ChapterItem) (*model.ChapterItem, error)
	QueryChapters() ([]*model.ChapterItem, error)
	UpdateChapter(id int64, keys []string, values []interface{}) error
	QueryBooks() ([]*model.GenreItemData, error)
	BuildBookChapter(comicId, chapterTitle string , chapterIndex int64) (int64, error)
}

type dao struct {
}

func New() Service { return new(dao) }

func (d *dao) BuildBookChapter(comicId, chapterTitle string, chapterIndex int64) (int64, error) {
	stmt, err := database.GetInstance().Prepare(
		`select id from comic_book_chapter where comic_id = ? and chapter_index = ?`)
	if err != nil {
		return 0, err
	}
	defer func() { _ = stmt.Close() }()
	var id int64
	_ = stmt.QueryRow(comicId, chapterIndex).Scan(&id)
	if id > 0 {
		return id, errors.New("duplicate.book")
	}
	stmt, err = database.GetInstance().Prepare(
		`insert into comic_book_chapter(comic_id, title, chapter_index) values(?,?,?)`)
	if err != nil {
		return 0, err
	}
	defer func() { _ = stmt.Close() }()
	ret, err := stmt.Exec(comicId, chapterTitle, chapterIndex)
	if err != nil {
		return 0, err
	}
	insertId, err := ret.LastInsertId()
	if err != nil {
		return 0, err
	}
	return insertId, err
}

func (d *dao) SaveGenre(input *model.GenreData) (*model.GenreData, error) {
	input.GenreId = encode.CrcEncode(input.Genre)
	stmt, err := database.GetInstance().Prepare(
		`select id from comic_genre where name = ?`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = stmt.Close() }()
	var id int64
	_ = stmt.QueryRow(input.Genre).Scan(&id)
	if id > 0 {
		input.Id = id
		return input, errors.New("duplicate.genre")
	}
	stmt, err = database.GetInstance().Prepare(
		`insert into comic_genre(name, book_count, sort_index, genre_id, ct_time) values(?,?,?,?,?)`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = stmt.Close() }()
	ret, err := stmt.Exec(input.Genre, input.Count, input.SortIndex, input.GenreId, time.Now().Unix())
	if err != nil {
		return nil, err
	}
	insertId, err := ret.LastInsertId()
	if err != nil {
		return nil, err
	}
	input.Id = insertId
	return input, err
}

func (d *dao) SaveBook(input *model.GenreItemData) (*model.GenreItemData, error) {
	if len(input.GenreId) == 0 {
		return nil, errors.New("save.book.genre_id.empty")
	}
	input.BookId = encode.CrcEncode(input.Identification)
	stmt, err := database.GetInstance().Prepare(
		`select id from comic_book where name = ?`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = stmt.Close() }()
	var id int64
	_ = stmt.QueryRow(input.ComicTitle).Scan(&id)
	if id > 0 {
		input.Id = id
		return input, errors.New("duplicate.book")
	}
	stmt, err = database.GetInstance().Prepare(
		`insert into comic_book(
                       name, comic_id, chapter_count, genre_id, 
                       update_status, cover, primary_color, release_time, 
                       author, name_alter, brief) 
				values(?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = stmt.Close() }()
	ret, err := stmt.Exec(input.ComicTitle, input.BookId, input.ChaptersCount, input.GenreId,
		input.ComicStatus, input.Illustration, input.MainColor, input.ReleaseTime,
		input.Author, input.Alternative, input.Introduction)
	if err != nil {
		return nil, err
	}
	insertId, err := ret.LastInsertId()
	if err != nil {
		return nil, err
	}
	input.Id = insertId
	return input, err
}

func (d *dao) SaveChapter(input *model.ChapterItem) (*model.ChapterItem, error) {
	stmt, err := database.GetInstance().Prepare(
		`select id, token from comic_chapter where comic_id = ? and chapter_index = ? and chapter_part_index = ?`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = stmt.Close() }()
	var id int64
	var token string
	_ = stmt.QueryRow(input.ComicId, input.ChapterIndex, input.ChapterPartIndex).Scan(&id, &token)
	if id > 0 {
		input.Id = id
		input.Token = token
		return input, errors.New("duplicate.chapter")
	}
	stmt, err = database.GetInstance().Prepare(
		`insert into comic_chapter(
                          comic_id, title, chapter_index, 
                          chapter_part_index, token) 
				values(?,?,?,?,?)`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = stmt.Close() }()
	token = str.RandStr(4)
	ret, err := stmt.Exec(input.ComicId, input.ComicChapterTitle, input.ChapterIndex,
		input.ChapterPartIndex, input, token)
	if err != nil {
		return nil, err
	}
	lastId, err := ret.LastInsertId()
	if err != nil {
		return nil, err
	}
	input.Id = lastId
	input.Token = token
	return input, err
}

func (d *dao) QueryBooks() ([]*model.GenreItemData, error) {
	stmt, err := database.GetInstance().Prepare(`select id, name, 
       comic_id, chapter_count, genre_id, update_status, cover, 
       primary_color, release_time, author, name_alter, brief from comic_book`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = stmt.Close() }()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	var dt []*model.GenreItemData
	for rows.Next() {
		var d model.GenreItemData
		err = rows.Scan(&d.Id, &d.ComicTitle, &d.BookId, &d.ChaptersCount, &d.GenreId,
			&d.ComicStatus, &d.Illustration, &d.MainColor, &d.ReleaseTime, &d.Author, &d.Alternative,
			&d.Introduction)
		if err == nil {
			dt = append(dt, &d)
		}
	}
	return dt, nil
}

func (d *dao) QueryChapters() ([]*model.ChapterItem, error) {
	stmt, err := database.GetInstance().Prepare(`select id, comic_id, chapter_index, 
       chapter_part_index, token from comic_chapter`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = stmt.Close() }()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	var dt []*model.ChapterItem
	for rows.Next() {
		var d model.ChapterItem
		err = rows.Scan(&d.Id, &d.ComicId, &d.ChapterIndex,
			&d.ChapterPartIndex, &d.Token)
		if err == nil {
			dt = append(dt, &d)
		}
	}
	return dt, nil
}

func (d *dao) UpdateChapter(id int64, keys []string, values []interface{}) error {
	var updates []string
	for _, k := range keys {
		updates = append(updates, fmt.Sprintf("%s = ?", k))
	}
	values = append(values, id)
	updateSql := fmt.Sprintf(`update comic_chapter set %s where id = ?`, strings.Join(updates, ","))
	stmt, err := database.GetInstance().Prepare(updateSql)
	if err != nil {
		return err
	}
	defer func() { _ = stmt.Close() }()
	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}
	return nil
}
