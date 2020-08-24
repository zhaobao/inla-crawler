package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/tasks/meditation/tide/model"
)

type MedDao struct {
}

var MedService = new(MedDao)

const (
	TypeAlbum    = 1
	TypeSection  = 2
	TypeResource = 3
)

var typeIdMapping = map[int]string{
	TypeAlbum:    "med_id",
	TypeSection:  "section_id",
	TypeResource: "res_id",
}

var tableMapping = map[int]string{
	TypeAlbum:    "tide_med",
	TypeSection:  "tide_med_section",
	TypeResource: "tide_med_res",
}

func (m *MedDao) IdExists(id string, t int) (int64, error) {
	if len(id) == 0 {
		return 0, errors.New("id.empty")
	}
	querySql := fmt.Sprintf(`select id from %s where %s = ?`, tableMapping[t], typeIdMapping[t])
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

func (m *MedDao) SaveMed(item model.MedRow) error {
	id, err := m.IdExists(item.MedId, item.Type)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}
	if id > 0 {
		return nil
	}
	insertSql := fmt.Sprintf(`insert into %s(med_id, tag_ids, name, description, 
                     primary_color, created_at, updated_at, sort_key, cover_link) values(?,?,?,?,?,?,?,?,?)`, tableMapping[item.Type])
	stmt, err := database.GetInstance().Prepare(insertSql)
	if err != nil {
		return err
	}
	defer func() { _ = stmt.Close() }()
	_, err = stmt.Exec(item.MedId, item.TagIds, item.Name, item.Description,
		item.PrimaryColor, item.CreatedAt, item.UpdatedAt, item.SortKey, item.CoverLink)
	if err != nil {
		return err
	}
	return nil
}

func (m *MedDao) SaveSection(item model.MedRow) error {
	id, err := m.IdExists(item.SectionId, item.Type)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}
	if id > 0 {
		return nil
	}
	insertSql := fmt.Sprintf(`insert into %s(med_id, section_id, name, demo_link) values(?,?,?,?)`, tableMapping[item.Type])
	stmt, err := database.GetInstance().Prepare(insertSql)
	if err != nil {
		return err
	}
	defer func() { _ = stmt.Close() }()
	_, err = stmt.Exec(item.MedId, item.SectionId, item.Name, item.DemoLink)
	if err != nil {
		return err
	}
	return nil
}

func (m *MedDao) SaveRes(item model.MedRow) error {
	id, err := m.IdExists(item.ResId, item.Type)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}
	if id > 0 {
		return nil
	}
	insertSql := fmt.Sprintf(`insert into %s(med_id, section_id, res_id, name, 
                     hash, duration, res_link) values(?,?,?,?,?,?,?)`, tableMapping[item.Type])
	stmt, err := database.GetInstance().Prepare(insertSql)
	if err != nil {
		return err
	}
	defer func() { _ = stmt.Close() }()
	_, err = stmt.Exec(item.MedId, item.SectionId, item.ResId, item.Name, item.Hash, item.Duration, item.ResLink)
	if err != nil {
		return err
	}
	return nil
}

func (m *MedDao) UpdateMedGroup(id, group string) {
	stmt, err := database.GetInstance().Prepare("update tide_med set `group` = `group` || ',' || ? where med_id = ?")
	if err != nil {
		return
	}
	defer func() { _ = stmt.Close() }()
	_, _ = stmt.Exec(group, id)
}

func (m *MedDao) UpdateMedTagIds(id, group string) {
	stmt, err := database.GetInstance().Prepare("update tide_med set `tag_ids` = ? where med_id = ?")
	if err != nil {
		return
	}
	defer func() { _ = stmt.Close() }()
	_, _ = stmt.Exec(group, id)
}
