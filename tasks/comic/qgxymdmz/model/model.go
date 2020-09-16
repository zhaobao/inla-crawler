package model

type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Genres
type GenresResponse struct {
	BaseResponse
	Data []GenreData `json:"data"`
}

type GenreData struct {
	Id           int64  `json:"id"`
	GenreId      string `json:"genre_id"`
	Genre        string `json:"genre"`
	Count        int    `json:"count"`
	SortIndex    int    `json:"sortIndex"`
	ShowInPortal bool   `json:"showInPortal"`
}

// GenreItem
type GenreItemsResponse struct {
	BaseResponse
	Data struct {
		PageNo     int             `json:"pageNo"`
		PageSize   int             `json:"pageSize"`
		PageCount  int             `json:"pageCount"`
		TotalCount int             `json:"totalCount"`
		ModelList  []GenreItemData `json:"modelList"`
	} `json:"data"`
}

type GenreItemData struct {
	Id             int64       `json:"id"`
	BookId         string      `json:"book_id"`
	Identification string      `json:"identification"`
	ComicTitle     string      `json:"comicTitle"`
	Alternative    string      `json:"alternative"`
	GenreId        string      `json:"genre_id"`
	Genre          string      `json:"genre"`
	Author         string      `json:"author"`
	ReleaseTime    string      `json:"releaseTime"`
	ComicStatus    string      `json:"comicStatus"`
	Illustration   string      `json:"illustration"`
	Introduction   string      `json:"introduction"`
	MainColor      string      `json:"mainColor"`
	ChaptersCount  int         `json:"chaptersCount"`
	Chapters       interface{} `json:"chapters"`
	Level          int         `json:"level"`
	Online         bool        `json:"online"`
}

// Chapter
type ChapterResponse struct {
	BaseResponse
	Data struct {
		PageNo     int           `json:"pageNo"`
		PageSize   int           `json:"pageSize"`
		PageCount  int           `json:"pageCount"`
		TotalCount int           `json:"totalCount"`
		ModelList  []ChapterItem `json:"modelList"`
	} `json:"data"`
}

type ChapterItem struct {
	Id                int64       `json:"id"`
	ComicChapterTitle string      `json:"comicChapterTitle"`
	ChapterIndex      float64     `json:"chapterIndex"`
	ChapterPartIndex  int         `json:"chapterPartIndex"`
	ComicId           string      `json:"comicId"`
	ImgLink           string      `json:"imgLink"`
	LastChapter       bool        `json:"lastChapter"`
	ShowIndex         interface{} `json:"showIndex"`
	Token             string      `json:"token"`

	Clean int `json:"clean"`
}
