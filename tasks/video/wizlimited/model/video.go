package model

//{
//    "id": 1,
//    "s_id": 1557306089,
//    "update_time": 1564563401,
//    "width": 720,
//    "height": 364,
//    "size": 2865964,
//    "duration": 38,
//    "link": "http://d2eshl90wojc4s.cloudfront.net/video/2019/05/08/yk76vl.mp4",
//    "link_shorten": "http://bit.ly/2OS0hpE",
//    "link_cover": "",
//    "link_src": "",
//    "title": "Dumbbell Zercher Squat 2",
//    "visit": 0,
//    "star": 0,
//    "download": 0,
//    "hash": "",
//    "vip": 0,
//    "quality": 0,
//    "s_index": 4,
//    "g_id": 8,
//    "g_title": "Day 8",
//    "g_index": 0,
//    "country": "",
//    "watermark": 2,
//    "category": "",
//    "type": 100
//  }

type VideoItem struct {
	SId       int64   `json:"s_id"`
	SIndex    int     `json:"s_index"`
	GId       int64   `json:"g_id"`
	GIndex    int64   `json:"g_index"`
	GTitle    string  `json:"g_title"`
	Watermark int     `json:"watermark"`
	Type      int     `json:"type"`
	Width     int     `json:"width"`
	Size      int     `json:"size"`
	Height    int     `json:"height"`
	Duration  float64 `json:"duration"`
	Link      string  `json:"link"`
	Title     string  `json:"title"`
}

type VideoRow struct {
	ResId         string  `json:"res_id"`
	ResIndex      int     `json:"res_index"`
	ResTitle      string  `json:"res_title"`
	ResLink       string  `json:"res_link"`
	GroupId       string  `json:"group_id"`
	GroupIndex    int64   `json:"group_index"`
	GroupTitle    string  `json:"group_title"`
	VideoWidth    int     `json:"video_width"`
	VideoHeight   int     `json:"video_height"`
	VideoSize     int     `json:"video_size"`
	VideoDuration float64 `json:"video_duration"`
	TypeId        string  `json:"type_id"`
	TypeName      string  `json:"type_name"`
	Source        string  `json:"source"`
	CountPlay     int     `json:"count_play"`
	CountLove     int     `json:"count_love"`
	CountDown     int     `json:"count_down"`
	CtTime        int64   `json:"ct_time"`
	RawLink       string  `json:"raw_link"`
}
