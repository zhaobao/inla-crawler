package model

type MusicItem struct {
	Scenes    []MusicScene `json:"scenes"`
	AdBanners []AdBanner   `json:"ad_banners"`
	AllTags   []Tag        `json:"all_tags"`
}

type AdBanner struct {
	Id           string                   `json:"id"`
	Content      map[string]AdContentItem `json:"content"`
	Style        string                   `json:"style"`
	ActivityName string                   `json:"activity_name"`
	RowPosition  int                      `json:"row_position"`
	Regions      []string                 `json:"regions"`
	Type         string                   `json:"type"`
	Status       string                   `json:"status"`
	UpdatedAt    int64                    `json:"updated_at"`
	CreatedAt    int64                    `json:"created_at"`
}

type AdContentItem struct {
	Layers []AdLayer `json:"layers"`
	Scheme string    `json:"scheme"`
}

type AdLayer struct {
	Image string `json:"image"`
	Align string `json:"align"`
}

type MusicScene struct {
	Id                string                 `json:"id"`
	Status            string                 `json:"status"`
	Name              map[string]string      `json:"name"`
	Tags              []string               `json:"tags"`
	DemoSoundUrl      string                 `json:"demo_sound_url"`
	DemoSoundUrlMp3   string                 `json:"demo_sound_url_mp3"`
	CoverUrl          string                 `json:"cover_url"`
	Thumbnail         string                 `json:"thumbnail"`
	IconUrl           string                 `json:"icon_url"`
	Types             []string               `json:"types"`
	Preset            bool                   `json:"preset"`
	Duration          int                    `json:"duration"`
	Background        map[string]string      `json:"background"`
	Description       map[string]string      `json:"description"`
	DescriptionAuthor map[string]string      `json:"description_author"`
	SubTitle          map[string]string      `json:"sub_title"`
	PrimaryColor      string                 `json:"primary_color"`
	SecondaryColor    string                 `json:"secondary_color"`
	SortKey           int                    `json:"sort_key"`
	UpdatedAt         int64                  `json:"updated_at"`
	CreatedAt         int64                  `json:"created_at"`
	Stats             map[string]interface{} `json:"stats"`
	TagsV2            []Tag                  `json:"tags_v2"`
	PlayList          []PlayListItem         `json:"playlist"`
}

type PlayListItem struct {
	Id       string            `json:"id"`
	Name     map[string]string `json:"name"`
	Duration int               `json:"duration"`
	Sounds   []SoundItem       `json:"sounds"`
}

type SoundItem struct {
	Id        string `json:"id"`
	Hash      string `json:"hash"`
	Key       string `json:"key"`
	UpdatedAt int64  `json:"updated_at"`
	CreatedAt int64  `json:"created_at"`
}
