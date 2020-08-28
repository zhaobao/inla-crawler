package main

import (
	"encoding/json"
	"fmt"
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/libs/downloader"
	"inla/inla-crawler/libs/net"
	"inla/inla-crawler/libs/shell"
	"inla/inla-crawler/tasks/comic/qgxymdmz/dao"
	"inla/inla-crawler/tasks/comic/qgxymdmz/model"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

const workerDir = "/Volumes/extend/crawler/comic"

var headers = map[string]string{
	"Origin":     "http://wap.qgxymdmz.com",
	"Referer":    "http://wap.qgxymdmz.com/pinkmanga/l/mobile/manga-genres.html?scene=genres",
	"User-Agent": "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Mobile Safari/537.36",
}

const (
	linkGenres        = `http://content.mobgkt.com/api/comic/genres?showInPortal=true`
	linkGenreItems    = `http://content.mobgkt.com/api/comic/list?pageNo=1&pageSize=500&online=true&sortType=1&genre=%s&level=2`
	linkChapterDetail = `http://content.mobgkt.com/api/comic/chapter/list?comicId=%s&&pageNo=1&pageSize=100&chapterIndex=%d`
)

var workerPool = downloader.New(32)
var downloadedCount int32
var total int32
var d = dao.New()

func main() {
	database.Connect("tasks/qgxymdmz/db.sqlite")
	//defaultDao := dao.New()
	//items, err := defaultDao.QueryChapters()
	//if err != nil {
	//	log.Fatal(err.Error())
	//}
	//for i, item := range items {
	//	id, err := defaultDao.BuildBookChapter(item.ComicId, item.ComicChapterTitle, int64(item.ChapterIndex))
	//	fmt.Println(i, id, err)
	//}
	//fmt.Println("done")

	workerPool.Start()
	upload1chapter()
	select {}
}

// 封面
func uploadBookCover() {
	defaultDao := dao.New()
	items, err := defaultDao.QueryBooks()
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, item := range items {
		cover := filepath.Join(workerDir, "imgs", item.BookId, "cover.png")
		err, output := shell.Pipe("sh", "deploy.sh", cover, filepath.Join(item.BookId, "cover.png"))
		fmt.Println(err, output)
	}
}

// 每一本书的第一章上传
func upload1chapter() {
	defaultDao := dao.New()
	items, err := defaultDao.QueryChapters()
	if err != nil {
		log.Fatal(err.Error())
	}

	memoryCache := make(map[string]bool)
	count := 0
	for _, item := range items {
		if item.ChapterIndex >= 19 && item.ChapterIndex <= 20 {
			fname := fmt.Sprintf("%v_%d_%s.jpg", item.ChapterIndex, item.ChapterPartIndex, item.Token)
			if _, ok := memoryCache[fname]; !ok {
				target := fmt.Sprintf("%s/imgs/%s/%s", workerDir, item.ComicId, fname)
				destination := fmt.Sprintf("%s/%s", item.ComicId, fname)
				count++
				workerPool.PutTask(createUploadTask(target, destination))
				fmt.Println("progress", count, "cost", float64(count)*0.005*6.99/1000, "CNY")
				memoryCache[fname] = true
			}
		}
	}
	fmt.Println("uploaded", count)
}

func createUploadTask(target, dest string) func() {
	return func() {
		err, output := shell.Pipe("sh", "deploy.sh", target, dest)
		atomic.AddInt32(&total, 1)
		fmt.Println(err, output, "finished", total)
	}
}

func createCoverWorker(id, link string) func() {
	return func() {
		downloadCover(id, link)
	}
}

func createChapterWorker(bookId, id string, count int) func() {
	return func() {
		chapterResponse := fetchChapterDetail(id, count)
		atomic.AddInt32(&total, int32(len(chapterResponse.Data.ModelList)))
		fmt.Println("TT", total)
		for _, vvv := range chapterResponse.Data.ModelList {
			vvv.ComicId = bookId
			retvvv, _ := d.SaveChapter(&vvv)
			vvv.Token = retvvv.Token
			workerPool.PutTask(createImageWorker(id, vvv))
		}
	}
}

func createImageWorker(id string, item model.ChapterItem) func() {
	return func() {
		downloadPage(id, item.ImgLink, fmt.Sprintf("%v_%d_%s", item.ChapterIndex, item.ChapterPartIndex, item.Token), item.ComicId)
	}
}

// genre/:genre/:id/:chapterIndex/:chapterPartIndex
func downloadPage(id, link, name, nid string) {
	atomic.AddInt32(&total, -1)
	if len(name) == 0 {
		log.Fatal("name.empty", id, link)
	}
	uri, _ := url.Parse(link)
	diskFile := filepath.Join(workerDir, "imgs", id, name+filepath.Ext(uri.Path))
	_ = os.MkdirAll(filepath.Dir(diskFile), 0755)
	if _, err := os.Stat(diskFile); err == nil {
		fmt.Println("<<< download.page", diskFile, "on.disk")
		return
	}
	buf, err := net.FetchResponse(http.MethodGet, link, nil, headers, 0)
	if err != nil {
		fmt.Println("downloadPage", err)
		return
	}
	if len(buf) > 0 {
		_ = ioutil.WriteFile(diskFile, buf, 0644)
		fmt.Println("<<< download.page", total, time.Now().Format("2006/01/02 13:14:15"), diskFile, nid, "save.disk")
	}
	atomic.AddInt32(&downloadedCount, 1)
}

// genre/:genre/:id/cover.png
func downloadCover(id, link string) {
	uri, _ := url.Parse(link)
	diskFile := filepath.Join(workerDir, "covers", id+filepath.Ext(uri.Path))
	if _, err := os.Stat(diskFile); err == nil {
		fmt.Println("downloadCover", "exists", err)
		return
	}
	buf, err := net.FetchResponse(http.MethodGet, link, nil, headers, 0)
	if err != nil {
		fmt.Println("downloadCover", "fetch", err)
		log.Fatal(err)
	}
	_ = ioutil.WriteFile(diskFile, buf, 0644)
	fmt.Println("download.cover", diskFile, "save.disk")
}

//curl 'http://content.mobgkt.com/api/comic/chapter/list?comicId=4daca2a0-2355-4a4a-b743-fbcf97217f91&&pageNo=1&pageSize=100&chapterIndex=206' \
//  -H 'Connection: keep-alive' \
//  -H 'Accept: application/json, text/javascript, */*; q=0.01' \
//  -H 'User-Agent: Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Mobile Safari/537.36' \
//  -H 'Origin: http://wap.qgxymdmz.com' \
//  -H 'Referer: http://wap.qgxymdmz.com/pinkmanga/l/mobile/manga-inside.html?Id=4daca2a0-2355-4a4a-b743-fbcf97217f91&chapter=2&scene=genre' \
//  -H 'Accept-Language: zh-CN,zh;q=0.9,en-SG;q=0.8,en;q=0.7' \
//  --compressed \
//  --insecure
func fetchChapterDetail(comicId string, count int) model.ChapterResponse {

	var dt model.ChapterResponse
	dt.Data.ModelList = make([]model.ChapterItem, 0)
	diskFile := filepath.Join(workerDir, "books", comicId, "data.json")
	diskDir := filepath.Dir(diskFile)
	_ = os.MkdirAll(diskDir, 0755)
	if _, err := os.Stat(diskFile); err == nil {
		buf, _ := ioutil.ReadFile(diskFile)
		_ = json.Unmarshal(buf, &dt)
	} else {
		for i := 1; i < count; i++ {
			fmt.Println("chapter", comicId, i, count)
			targetLink := fmt.Sprintf(linkChapterDetail, comicId, i)
			buf, err := net.FetchResponse(http.MethodGet, targetLink, nil, headers, 0)
			if err != nil {
				log.Fatal("fetchChapterDetail.http", err)
			}
			if len(buf) == 0 {
				log.Fatal("fetchChapterDetail.buf.empty")
			}
			var chapterResponse model.ChapterResponse
			err = json.Unmarshal(buf, &chapterResponse)
			if err != nil {
				log.Fatal("fetchChapterDetail.json", err, string(buf))
			}
			for _, v := range chapterResponse.Data.ModelList {
				dt.Data.ModelList = append(dt.Data.ModelList, v)
			}
		}
		if len(dt.Data.ModelList) > 0 {
			fmt.Println("save.books.json")
			buf, _ := json.Marshal(dt)
			_ = ioutil.WriteFile(diskFile, buf, 0644)
		}
	}
	return dt
}

// genre/:genre/:id
//func createItemDir(genre, id string) {
//	genre = strings.Replace(genre, " ", "_", -1)
//	dir := filepath.Join(workerDir, "genre", genre, id)
//	if _, err := os.Stat(dir); err != nil {
//		_ = os.MkdirAll(dir, 0755)
//	}
//}

//curl 'http://content.mobgkt.com/api/comic/list?pageNo=1&pageSize=20&online=true&sortType=1&genre=Action&level=2' \
//  -H 'Connection: keep-alive' \
//  -H 'Accept: application/json, text/javascript, */*; q=0.01' \
//  -H 'User-Agent: Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Mobile Safari/537.36' \
//  -H 'Origin: http://wap.qgxymdmz.com' \
//  -H 'Referer: http://wap.qgxymdmz.com/pinkmanga/l/mobile/manga-genres.html?scene=genres' \
//  -H 'Accept-Language: zh-CN,zh;q=0.9,en-SG;q=0.8,en;q=0.7' \
//  --compressed \
//  --insecure
func fetchGenreItems(genre string) model.GenreItemsResponse {
	var write bool
	var buf []byte
	var err error
	diskFile := filepath.Join(workerDir, "genre", fmt.Sprintf("fetchGenreItems.%s.json", genre))
	if _, err := os.Stat(diskFile); err == nil {
		fmt.Println(linkGenres, "disk.file.ok")
		buf, _ = ioutil.ReadFile(diskFile)
	} else {
		write = true
		targetLink := fmt.Sprintf(linkGenreItems, genre)
		buf, err = net.FetchResponse(http.MethodGet, targetLink, nil, headers, 0)
		if err != nil {
			log.Fatal("fetchGenreItems.http", err)
		}
		if len(buf) == 0 {
			log.Fatal("fetchGenreItems.buf.empty")
		}
	}
	var genreItemResponse model.GenreItemsResponse
	err = json.Unmarshal(buf, &genreItemResponse)
	if err != nil {
		log.Fatal("fetchGenreItems.json", err)
	}
	if write {
		_ = ioutil.WriteFile(diskFile, buf, 0644)
	}
	return genreItemResponse

}

//curl 'http://content.mobgkt.com/api/comic/genres?showInPortal=true' \
//  -H 'Connection: keep-alive' \
//  -H 'Accept: application/json, text/javascript, */*; q=0.01' \
//  -H 'User-Agent: Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Mobile Safari/537.36' \
//  -H 'Origin: http://wap.qgxymdmz.com' \
//  -H 'Referer: http://wap.qgxymdmz.com/pinkmanga/l/mobile/manga-genres.html?scene=genres' \
//  -H 'Accept-Language: zh-CN,zh;q=0.9,en-SG;q=0.8,en;q=0.7' \
//  --compressed \
//  --insecure
func fetchGenres() model.GenresResponse {
	var write bool
	var buf []byte
	var err error
	diskFile := filepath.Join(workerDir, "genre", "fetchGenres.json")
	if _, err := os.Stat(diskFile); err == nil {
		buf, err = ioutil.ReadFile(diskFile)
		if err != nil {
			log.Fatal("fetchGenres.read", err)
		}
	} else {
		write = true
		buf, err = net.FetchResponse(http.MethodGet, linkGenres, nil, headers, 0)
		if err != nil {
			log.Fatal("fetchGenres.http", err)
		}
		if len(buf) == 0 {
			log.Fatal("fetchGenres.buf.empty")
		}
	}
	var genresResponse model.GenresResponse
	err = json.Unmarshal(buf, &genresResponse)
	if err != nil {
		log.Fatal("fetchGenres.json", err)
	}
	if write {
		_ = ioutil.WriteFile(diskFile, buf, 0644)
	}
	return genresResponse
}
