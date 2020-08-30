package model

type Book struct {
	Id            string    `json:"id"`
	Cover         string    `json:"cover"`
	Title         string    `json:"title"`
	NameAlter     string    `json:"name_alter"`
	Genre         string    `json:"genre"`
	Status        string    `json:"status"`
	Author        string    `json:"author"`
	IsHot         bool      `json:"is_hot"`
	IsNew         bool      `json:"is_new"`
	Link          string    `json:"link"`
	Brief         string    `json:"brief"`
	Chapters      []Chapter `json:"chapters"`
	ChaptersCount int       `json:"chapters_count"`
	GenreId       string    `json:"genre_id"`
	Source        string    "readnovelfull"
}

type Chapter struct {
	Id     int64  `json:"id"`
	BookId string `json:"book_id"`
	Link   string `json:"link"`
	Name   string `json:"name"`
	Index  int    `json:"index"`
	Token  string `json:"token"`
}
