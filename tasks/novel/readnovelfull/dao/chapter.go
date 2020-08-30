package dao

import (
	"database/sql"
	"fmt"
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/libs/str"
	"inla/inla-crawler/tasks/novel/readnovelfull/model"
	"log"
)

type chapterImp struct {
}

func NewChapter() *chapterImp { return new(chapterImp) }

func (c *chapterImp) MakeAsDone(id int64) (int64, error) {
	stmt, err := database.GetInstance().Prepare(`update novel_chapter set done = 1 where id = ?`)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() { _ = stmt.Close() }()
	ret, err := stmt.Exec(id)
	if err != nil {
		log.Fatal(err.Error())
	}
	return ret.RowsAffected()
}

func (c *chapterImp) Save(ch model.Chapter) {
	stmt, err := database.GetInstance().Prepare(
		`select id from novel_chapter where book_id = ? and chapter_index = ?`)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() { _ = stmt.Close() }()
	var id int64
	err = stmt.QueryRow(ch.BookId, ch.Index).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err.Error())
		}
	}
	if id > 0 {
		fmt.Println("duplicate.chapter", ch.BookId, ch.Index)
		return
	}
	stmt, err = database.GetInstance().Prepare(`insert into novel_chapter(book_id, chapter_index, token, 
                          done, name, src_link) values(?,?,?,?,?,?)`)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() { _ = stmt.Close() }()
	_, err = stmt.Exec(ch.BookId, ch.Index, str.RandStr(5), 0, ch.Name, ch.Link)
	if err != nil {
		log.Fatal(err.Error())
	}
}
