package dao

import (
	"database/sql"
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/libs/encode"
	"inla/inla-crawler/tasks/novel/qgxymdmz/model"
	"log"
)

type BookService struct {
}

func NewBook() *BookService {
	return new(BookService)
}

func (d *BookService) Add(item model.NovelBook) (int64, error) {
	novelId := encode.CrcEncode(item.Title)
	item.BookId = novelId
	stmt, err := database.GetInstance().Prepare(
		`select id from novel_book where book_id = ?`)
	if err != nil {
		log.Fatal("novel.prepare.add", err)
	}
	defer func() { _ = stmt.Close() }()
	var id int64
	err = stmt.QueryRow(novelId).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal("novel.dao.add", err)
		}
	}
	if id > 0 {
		return id, nil
	}
	stmt, err = database.GetInstance().Prepare(
		`insert into novel_book(
                       name, cover, brief, genre_id, author, chapter_count, 
                       name_alter, book_id, finished) 
                       values (?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		log.Fatal("novel.prepare.add", err)
	}
	defer func() { _ = stmt.Close() }()
	ret, err := stmt.Exec(item.Title, item.Cover, item.Introduction,
		item.GenreId, item.Author, item.ChaptersCount, item.Tag, item.BookId,
		item.Finished)
	if err != nil {
		log.Fatal("novel.exec.add", err)
	}
	lastId, err := ret.LastInsertId()
	if err != nil {
		log.Fatal("novel.last.id.add", err)
	}
	return lastId, nil
}
