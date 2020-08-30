package dao

import (
	"database/sql"
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/libs/encode"
	"log"
)

type genreImp struct{}

func NewGenre() *genreImp { return new(genreImp) }

func (d *genreImp) FindIdByName(name string) (string, error) {
	var id string
	err := database.GetInstance().QueryRow(`select genre_id from novel_genre where name = ?`, name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			gid := encode.CrcEncode(name)
			_, err := database.GetInstance().Exec(`insert into novel_genre(name, genre_id) values(?,?)`, name, gid)
			if err == nil {
				return gid, nil
			}
		} else {
			log.Fatal("genre", "FindIdByName", err.Error())
		}
	}
	return id, err
}
