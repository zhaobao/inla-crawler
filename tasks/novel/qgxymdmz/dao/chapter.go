package dao

import (
	"database/sql"
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/libs/str"
	"inla/inla-crawler/tasks/novel/qgxymdmz/model"
	"log"
)

type ChapterService struct {
}

func NewChapter() *ChapterService {
	return new(ChapterService)
}

func (d *ChapterService) Add(item model.NovelChapter) (model.NovelChapter, error) {
	item.Token = str.RandStr(4)
	stmt, err := database.GetInstance().Prepare(
		`select id, token from novel_chapter where book_id = ? and chapter_index = ?`)
	if err != nil {
		log.Fatal("novel.prepare.add", err)
	}
	defer func() { _ = stmt.Close() }()
	var id int64
	var token string
	err = stmt.QueryRow(item.BookId, item.ChapterIndex).Scan(&id, &token)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal("novel.dao.add", err)
		}
	}
	if id > 0 {
		item.Id = id
		item.Token = token
		return item, nil
	}
	stmt, err = database.GetInstance().Prepare(
		`insert into novel_chapter(
                          book_id, chapter_index, token) values (?,?,?)`)
	if err != nil {
		log.Fatal("novel.prepare.add", err)
	}
	defer func() { _ = stmt.Close() }()
	ret, err := stmt.Exec(item.BookId, item.ChapterIndex, item.Token)
	if err != nil {
		log.Fatal("novel.exec.add", err)
	}
	lastId, err := ret.LastInsertId()
	if err != nil {
		log.Fatal("novel.last.id.add", err)
	}
	item.Id = lastId
	return item, nil
}
