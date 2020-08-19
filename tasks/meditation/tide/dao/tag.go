package dao

import (
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/tasks/meditation/tide/model"
)

type TagDao struct {
}

var TagService = new(TagDao)

func (t *TagDao) IdExists(id string) (int64, error) {
	querySql := `select id from tide_tag where tag_id = ?`
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

func (t *TagDao) Save(row model.TagRow) (int64, error) {
	pid, _ := t.IdExists(row.TagId)
	if pid > 0 {
		return pid, nil
	}
	insertSql := `insert into tide_tag(tag_id, sort_key, key, type, name) values(?,?,?,?,?)`
	stmt, err := database.GetInstance().Prepare(insertSql)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = stmt.Close()
	}()
	ret, err := stmt.Exec(row.TagId, row.SortKey, row.Key, row.Type, row.Name)
	if err != nil {
		return 0, err
	}
	return ret.LastInsertId()
}
