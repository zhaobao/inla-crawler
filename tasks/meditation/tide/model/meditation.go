package model

type MeditationObj struct {
	TodayMeditation Album      `json:"today_meditation"`
	Albums          []Album    `json:"albums"`
	MiniMeditations []Album    `json:"mini_meditations"`
	AdBanners       []AdBanner `json:"ad_banners"`
	AllTags         []Tag      `json:"all_tags"`
}

type Album struct {
	Id                  string                 `json:"id"`
	Name                map[string]string      `json:"name"`
	Image               string                 `json:"image"`
	PrimaryColor        string                 `json:"primary_color"`
	Description         map[string]string      `json:"description"`
	DurationDescription map[string]string      `json:"duration_description"`
	Sections            []Section              `json:"sections"`
	UpdatedAt           int64                  `json:"updated_at"`
	CreatedAt           int64                  `json:"created_at"`
	SortKey             int                    `json:"sort_key"`
	Stats               map[string]interface{} `json:"stats"`
	TagsV2              []Tag                  `json:"tags_v2"`
}

type Section struct {
	Id                  string            `json:"id"`
	Name                map[string]string `json:"name"`
	Description         map[string]string `json:"description"`
	DemoSoundUrlMp3     map[string]string `json:"demo_sound_url_mp3"`
	DurationDescription map[string]string `json:"duration_description"`
	Resources           []Resource        `json:"resources"`
}

type Resource struct {
	Duration  int      `json:"duration"`
	Languages []string `json:"languages"`
	Hash      string   `json:"hash"`
	HashKey   string   `json:"hash_key"`
	Name      string   `json:"name"`
	Speaker   Speaker  `json:"speaker"`
}

type Speaker struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}
