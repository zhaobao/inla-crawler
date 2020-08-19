package dao

import (
	"fmt"
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/tasks/meditation/tide/model"
)

type MusicDao struct {
}

var MusicService = new(MusicDao)

func (m *MusicDao) IdExists(id string) (int64, error) {
	querySql := `select id from tide_music where res_id = ?`
	stmt, err := database.GetInstance().Prepare(querySql)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = stmt.Close()
	}()
	var pid int64
	err = stmt.QueryRow(id).Scan(&pid)
	return pid, err
}

func (m *MusicDao) Save(row model.MusicRow) (int64, error) {
	pid, _ := m.IdExists(row.Id)
	if pid > 0 {
		return pid, nil
	}
	insertSql := `insert into tide_music(name, demo_link, cover_link, duration, 
                       description, sub_title, primary_color, secondary_color, 
                       updated_at, created_at, res_id, tag_ids, loved, hash) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	stmt, err := database.GetInstance().Prepare(insertSql)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = stmt.Close()
	}()
	ret, err := stmt.Exec(row.Name, row.DemoLink, row.CoverLink, row.Duration,
		row.Description, row.SubTitle, row.PrimaryColor, row.SecondaryColor,
		row.UpdatedAt, row.CreatedAt, row.Id, row.TagIds, row.Loved, row.Hash)
	if err != nil {
		return 0, err
	}
	return ret.LastInsertId()
}

func (m *MusicDao) UpdateGroup(id, group string) {
	stmt, err := database.GetInstance().Prepare("update tide_music set `group` = ? where res_id = ?")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() { _ = stmt.Close() }()
	_, err = stmt.Exec(group, id)
	if err != nil {
		fmt.Println(err)
	}
}
