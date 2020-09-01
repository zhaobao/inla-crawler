package dao

import (
	"database/sql"
	"inla/inla-crawler/libs/database"
	"log"
)

type bookImp struct{}

func NewBook() *bookImp { return new(bookImp) }

func (d *bookImp) Save(name, cover, gid, bid, source, srcLink, color string) (string, error) {
	var id string
	err := database.GetInstance().QueryRow(`select book_id from novel_book where book_id = ?`, name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err := database.GetInstance().Exec(`insert into novel_book(name, cover, genre_id, book_id, source, src_link, primary_color) values(?,?,?,?,?,?,?)`, name, color, gid, bid, source, srcLink, color)
			if err == nil {
				return gid, nil
			}
		} else {
			log.Fatal("book", "Save", err.Error())
		}
	}
	return id, err
}
