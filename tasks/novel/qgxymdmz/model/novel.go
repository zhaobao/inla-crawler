package model

type NovelBook struct {
	Title         string `json:"title"`
	Illustration  string `json:"illustration"`
	Cover         string `json:"cover"`
	Introduction  string `json:"introduction"`
	Tag           string `json:"tag"`
	Author        string `json:"author"`
	Src           string `json:"src"`
	Finished      bool   `json:"finished"`
	LastUpdate    int64  `json:"lastUpdate"`
	ChaptersCount int    `json:"chaptersCount"`
	Category      string `json:"category"`
	GenreId       string `json:"genre_id"`
	BookId        string `json:"book_id"`
}

type NovelChapter struct {
	Id           int64  `json:"id"`
	Token        string `json:"token"`
	BookId       string `json:"book_id"`
	ChapterIndex int    `json:"chapter_index"`
}

type NovelGenre struct {
	GenreId string `json:"genre_id"`
	Name    string `json:"name"`
	Count   int    `json:"count"`
}
