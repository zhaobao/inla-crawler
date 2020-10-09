package dao

type MusicBo struct {
	Id       int64    `json:"id"`
	Name     string   `json:"name"`
	Duration int      `json:"duration"`
	Cover    string   `json:"cover"`
	Album    string   `json:"album"`
	Src      string   `json:"src"`
	LyricSrc string   `json:"lyric_src"`
	Artists  []string `json:"artists"`

	ResId        string `json:"res_id"`
	ResLink      string `json:"res_link"`
	CoverLink    string `json:"cover_link"`
	LyricLink    string `json:"lyric_link"`
	PrimaryColor string `json:"primary_color"`
}
