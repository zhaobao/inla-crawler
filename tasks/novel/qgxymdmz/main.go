package main

import (
	"encoding/json"
	"fmt"
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/libs/downloader"
	"inla/inla-crawler/libs/encode"
	"inla/inla-crawler/libs/shell"
	"inla/inla-crawler/tasks/novel/qgxymdmz/dao"
	"inla/inla-crawler/tasks/novel/qgxymdmz/model"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
)

/*
整理小说内容到数据库
1. 整理info到数据库格式
2. 上传S3
*/

const workerDir = "/Volumes/extend/crawler/novel/qgxymdmz"

var uploaderPool *downloader.Client
var taskTotalCount int32
var taskDoneCount int32

func setup() {
	// 转换cover成jpg格式
	_ = filepath.Walk(workerDir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".png") {
			to := strings.ReplaceAll(path, ".png", ".jpg")
			output, err := shell.Pipe("sh", "cover.sh", path, to)
			fmt.Println(output, err, "sh", "convert.sh", path, to)
		}
		return nil
	})
}

func main() {
	do()
}

func do() {
	uploaderPool = downloader.New(16)
	uploaderPool.Start()
	done := make(chan bool)

	database.Connect("tasks/novel/qgxymdmz/db.sqlite")
	contents, err := loadFile(filepath.Join(workerDir, "nov_infos.json"))
	if err != nil {
		log.Fatal("readFile", err)
	}

	var items []model.NovelBook
	err = json.Unmarshal(contents, &items)
	if err != nil {
		log.Fatal("json.Unmarshal", err)
	}

	categories := make(map[string]int)
	for _, item := range items {
		if _, ok := categories[item.Category]; !ok {
			categories[item.Category] = 0
		}
		categories[item.Category] += 1
		item.GenreId = encode.CrcEncode(item.Category)

		cover := filepath.Join(workerDir, item.Category, item.Title, "illustration.jpg")
		if _, err := os.Stat(cover); err != nil {
			log.Fatal("cover.not.exists", err.Error())
		}
		item.Cover = "cover.jpg"
		saveBook(item)
		// cover 已经上传过
		item.BookId = encode.CrcEncode(item.Title)
		saveChapter(done, item)
	}
	saveGenres(categories)
	fmt.Println(len(items))

	<-done
}

func uploadTask(done chan bool, localFile, remoteFile string) func() {
	return func() {
		err, output := shell.Pipe("sh", "deploy.sh", localFile, remoteFile)
		atomic.AddInt32(&taskDoneCount, 1)
		fmt.Println(err, output)
		fmt.Println("--- progress", taskTotalCount, taskDoneCount, "cost", float64(taskDoneCount)*0.005*6.99/1000, "CNY")
		if taskDoneCount == taskTotalCount {
			done <- true
		}
	}
}

func saveChapter(done chan bool, book model.NovelBook) {
	d := dao.NewChapter()
	for i := 1; i <= book.ChaptersCount; i++ {
		item, err := d.Add(model.NovelChapter{ChapterIndex: i, BookId: book.BookId})
		if err != nil {
			log.Fatal("saveChapter.err", err.Error())
		}
		chapter := filepath.Join(workerDir, book.Category, book.Title, fmt.Sprintf("%d.txt", i))
		if _, err := os.Stat(chapter); err != nil {
			log.Fatal("chapter.not.exists", err.Error())
		}
		chapterRemote := filepath.Join(item.BookId, fmt.Sprintf("%d_%s.txt", item.ChapterIndex, item.Token))
		if i > 200 {
			atomic.AddInt32(&taskTotalCount, 1)
			uploaderPool.PutTask(uploadTask(done, chapter, chapterRemote))
			fmt.Println("+++", taskTotalCount, chapter, chapterRemote)
		}
	}
}

func saveBook(book model.NovelBook) () {
	id, err := dao.NewBook().Add(book)
	if err != nil {
		log.Fatal("saveBook.err", err.Error())
	}
	fmt.Println("saveBook", id, book.BookId)
}

func saveGenres(items map[string]int) {
	d := dao.NewGenre()
	for name, count := range items {
		id, err := d.Add(model.NovelGenre{Name: name, Count: count})
		if err != nil {
			log.Fatal("saveGenres.err", err.Error())
		}
		fmt.Println("saveGenres", id, name)
	}
}

func loadFile(input string) ([]byte, error) {
	return ioutil.ReadFile(input)
}
