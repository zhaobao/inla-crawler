package model

type GenreItem struct {
	GenreId string `json:"genre_id"`
	Name    string `json:"name"`
	Count   int    `json:"count"`
}
