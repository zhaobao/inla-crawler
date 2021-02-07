package main

import (
	"encoding/json"
	"fmt"
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/libs/encode"
	"inla/inla-crawler/libs/fs"
	"inla/inla-crawler/libs/shell"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const rootDir = `tasks/novel/spolishnovel`
const resDir = `/Volumes/extend/crawler/novel/spolish`

func main() {
	database.Connect(fmt.Sprintf("%s/db.sqlite", rootDir))
	saveBooksToDb()
}

type Book struct {
	Url           string `json:"url"`
	Title         string `json:"title"`
	Illustration  string `json:"illustration"`
	Introduction  string `json:"introduction"`
	Tag           string `json:"tag"`
	Author        string `json:"author"`
	Src           string `json:"src"`
	LastUpdate    string `json:"last_update"`
	ChaptersCount string `json:"chaptersCount"`
	Category      string `json:"category"`
	IsFinish      string `json:"is_finish"`
	IsHot         string `json:"is_hot"`
}

func saveBooksToDb() {
	bookFile := filepath.Join(resDir, "list_spolish_info.json")
	if !fs.FileExists(bookFile) {
		return
	}
	buf, err := ioutil.ReadFile(bookFile)
	if err != nil {
		log.Fatal(err)
	}
	var books []Book
	err = json.Unmarshal(buf, &books)
	if err != nil {
		log.Fatal(err)
	}
	for _, book := range books {
		bookId := encode.CrcEncode(filepath.Join("spolishnovel", book.Src))
		destFile := filepath.Join(rootDir, "output", bookId, "book.txt")
		destDir := filepath.Dir(destFile)
		if !fs.FileExists(destDir) {
			_ = os.MkdirAll(destDir, 0755)
		}
		if !fs.FileExists(destFile) {
			cpFile(filepath.Join(resDir, book.Src), destFile)
			_, _ = Save(book.Title, "", "", bookId, "spolishnovel", "", "")
		}
	}
}

func Save(name, cover, gid, bid, source, srcLink, color string) (string, error) {
	sqlExec := `
insert into book(name, cover, genre_id, book_id, source, src_link, primary_color)
values (?, ?, ?, ?, ?, ?, ?)
`
	_, err := database.GetInstance().Exec(sqlExec, name, cover, gid, bid, source, srcLink, color)
	if err == nil {
		log.Println("success .....", bid)
		return bid, nil
	}
	return bid, err
}

func cpFile(from, to string) {
	if fs.FileExists(to) {
		log.Println("to.exists")
		return
	}
	if !fs.FileExists(from) {
		log.Fatal("from.not.exists", from)
	}
	log.Println("cp", from, to)
	_, _ = shell.Pipe("cp", from, to)
}
