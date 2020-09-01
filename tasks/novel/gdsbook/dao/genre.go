package dao

import (
	"database/sql"
	"inla/inla-crawler/libs/database"
	"log"
)

type genreImp struct{}

func NewGenre() *genreImp { return new(genreImp) }

func (d *genreImp) Save(name, gid string, count int) (string, error) {
	var id string
	err := database.GetInstance().QueryRow(`select genre_id from novel_genre where name = ?`, name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err := database.GetInstance().Exec(`insert into novel_genre(name, genre_id, count) values(?,?,?)`, name, gid, count)
			if err == nil {
				return gid, nil
			}
		} else {
			log.Fatal("genre", "Save", err.Error())
		}
	}
	return id, err
}

func (d *genreImp) Find(name string) (string, error) {
	var id string
	err := database.GetInstance().QueryRow(`select genre_id from novel_genre where name = ?`, name).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	return id, err
}
