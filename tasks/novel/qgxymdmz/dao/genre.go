package dao

import (
	"database/sql"
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/libs/encode"
	"inla/inla-crawler/tasks/novel/qgxymdmz/model"
	"log"
)

type GenreService struct {
}

func NewGenre() *GenreService {
	return new(GenreService)
}

func (d *GenreService) Add(item model.NovelGenre) (int64, error) {
	genreId := encode.CrcEncode(item.Name)
	item.GenreId = genreId
	stmt, err := database.GetInstance().Prepare(
		`select id from novel_genre where genre_id = ?`)
	if err != nil {
		log.Fatal("novel.prepare.add", err)
	}
	defer func() { _ = stmt.Close() }()
	var id int64
	err = stmt.QueryRow(genreId).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal("novel.dao.add", err)
		}
	}
	if id > 0 {
		return id, nil
	}
	stmt, err = database.GetInstance().Prepare(
		`insert into novel_genre(name, genre_id, count) values (?,?,?)`)
	if err != nil {
		log.Fatal("novel.prepare.add", err)
	}
	defer func() { _ = stmt.Close() }()
	ret, err := stmt.Exec(item.Name, item.GenreId, item.Count)
	if err != nil {
		log.Fatal("novel.exec.add", err)
	}
	lastId, err := ret.LastInsertId()
	if err != nil {
		log.Fatal("novel.last.id.add", err)
	}
	return lastId, nil
}
