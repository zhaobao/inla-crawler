package dao

import (
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/tasks/video/wizlimited/model"
	"log"
)

type videoImp struct {
}

func NewVideo() *videoImp { return new(videoImp) }

func (v *videoImp) Save(item *model.VideoRow) (int64, error) {
	stmt, err := database.GetInstance().Prepare(`insert into assets_video_2(res_id, res_index, res_title, 
                         res_link, group_id, group_index, group_title, video_width, video_height, video_size, 
                         video_duration, type_id, type_name, source, count_play, count_love, count_down, 
                         ct_time) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = stmt.Close() }()
	ret, err := stmt.Exec(item.ResId, item.ResIndex, item.ResTitle,
		item.ResLink, item.GroupId, item.GroupIndex, item.GroupTitle, item.VideoWidth, item.VideoHeight, item.VideoSize,
		item.VideoDuration, item.TypeId, item.TypeName, item.Source, item.CountPlay, item.CountLove, item.CountDown,
		item.CtTime)
	if err != nil {
		log.Fatal(err)
	}
	return ret.LastInsertId()
}
