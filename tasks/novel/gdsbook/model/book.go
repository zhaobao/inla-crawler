package model

type BookItem struct {
	Url           string `json:"url"`
	Title         string `json:"title"`
	Illustration  string `json:"illustration"`
	Introduction  string `json:"introduction"`
	Tag           string `json:"tag"`
	Author        string `json:"author"`
	Src           string `json:"src"`
	ChaptersCount string `json:"chaptersCount"`
	Category      string `json:"category"`
}
