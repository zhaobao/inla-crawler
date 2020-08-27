package dao

import (
	"database/sql"
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/tasks/music/cozy/model"
)

type daoImp struct {
}

var Service = new(daoImp)

func (d *daoImp) SaveOrUpdate(row model.Cozy) (int64, error) {
	stmt, err := database.GetInstance().Prepare(`select id from music_cozy where res_id = ?`)
	if err != nil {
		return 0, err
	}
	defer func() { _ = stmt.Close() }()
	var id int64
	err = stmt.QueryRow(row.CozyId).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, err
		}
	}
	stmt, err = database.GetInstance().Prepare(`insert into music_cozy(res_id, title, sub_title, cover_link, 
                       res_link, category) values(?,?,?,?,?,?)`)
	if err != nil {
		return 0, err
	}
	defer func() { _ = stmt.Close() }()
	ret, err := stmt.Exec(row.CozyId, row.Headline, row.Subtitle, row.CoverLink, row.ResLink, row.Category)
	if err != nil {
		return 0, err
	}
	return ret.LastInsertId()
}

func (d *daoImp) UpdateColorByCover(cover, color string) {
	stmt, err := database.GetInstance().Prepare(`update music_cozy set primary_color = ? where cover_link = ?`)
	if err != nil {
		return
	}
	defer func() { _ = stmt.Close() }()
	_, _ = stmt.Exec(color, cover)
}

func (d *daoImp) UpdateDurationByRes(res string, duration float64) {
	stmt, err := database.GetInstance().Prepare(`update music_cozy set duration = ? where res_link = ?`)
	if err != nil {
		return
	}
	defer func() { _ = stmt.Close() }()
	_, _ = stmt.Exec(duration, res)
}
