package dao

import (
	"database/sql"
	"fmt"
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/tasks/novel/readnovelfull/model"
	"log"
	"strings"
)

type BookService struct {
}

func NewBook() *BookService {
	return new(BookService)
}

func (d *BookService) Add(item model.Book) (int64, error) {
	stmt, err := database.GetInstance().Prepare(
		`select id from novel_book where book_id = ?`)
	if err != nil {
		log.Fatal("novel.prepare.add", err)
	}
	defer func() { _ = stmt.Close() }()
	var id int64
	err = stmt.QueryRow(item.Id).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal("novel.dao.add", err)
		}
	}
	if id > 0 {
		fmt.Println("duplicate.book", item.Id)
		return id, nil
	}
	stmt, err = database.GetInstance().Prepare(
		`insert into novel_book(
                       name, cover, brief, genre_id, author, chapter_count, 
                       name_alter, book_id, finished, is_hot, is_new, source, src_link) 
                       values (?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		log.Fatal("novel.prepare.add", err)
	}
	defer func() { _ = stmt.Close() }()

	var isNew int
	if item.IsNew {
		isNew = 1
	}
	var isHot int
	if item.IsHot {
		isHot = 1
	}
	var finished int
	if strings.ToLower(item.Status) != "ongoing" {
		finished = 1
	}

	ret, err := stmt.Exec(item.Title, item.Cover, item.Brief,
		item.GenreId, item.Author, item.ChaptersCount, item.NameAlter, item.Id,
		finished, isHot, isNew, item.Source, item.Link)
	if err != nil {
		log.Fatal("novel.exec.add", err)
	}
	lastId, err := ret.LastInsertId()
	if err != nil {
		log.Fatal("novel.last.id.add", err)
	}
	return lastId, nil
}

func (d *BookService) UpdateGenreId(booId, genre string) {
	_, err := database.GetInstance().Exec(`update novel_book set genre_id = ? where book_id = ?`, genre, booId)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *BookService) UpdateIsHotIsNew(booId string, isHot, isNew bool) {
	var h int
	if isHot {
		h = 1
	}
	var n int
	if isNew {
		n = 1
	}
	_, err := database.GetInstance().Exec(`update novel_book set is_hot = ?, is_new = ? where book_id = ?`, h, n, booId)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *BookService) UpdateBookColorById(bookId, color string) {
	_, err := database.GetInstance().Exec(`update novel_book set primary_color = ? where book_id = ?`, color, bookId)
	if err != nil {
		log.Fatal(err)
	}
}
