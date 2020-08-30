package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/libs/downloader"
	"inla/inla-crawler/libs/encode"
	"inla/inla-crawler/libs/fs"
	"inla/inla-crawler/libs/net"
	"inla/inla-crawler/libs/shell"
	"inla/inla-crawler/libs/str"
	"inla/inla-crawler/tasks/novel/readnovelfull/dao"
	"inla/inla-crawler/tasks/novel/readnovelfull/model"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const rootDir = "tasks/novel/readnovelfull"
const bookLink = "https://readnovelfull.com/latest-release-novel?page=%d"
const srcHost = "https://readnovelfull.com"
const chapterTpl = `<!DOCTYPE html><html lang="en-US"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, height=device-height, initial-scale=1.0, user-scalable=0, minimum-scale=1.0, maximum-scale=1.0"><title>Read: {{title}}</title><style>body{font-size:16px;font-family:arial, tahoma, verdana, sans-serif}</style></head><body>{{body}}</body></html>`

var putCount uint32
var getCount uint32
var markCh = make(chan int64, 0)

func main() {
	database.Connect(fmt.Sprintf("%s/db.sqlite", rootDir))
	client := downloader.New(128)
	client.Start()
	//saveGenres()
	//saveBooks()
	//populateGenres()

	//populateIsNot()

	//go func() {
	//	var lastCount uint32
	//	ticker := time.NewTicker(time.Minute)
	//	for {
	//		select {
	//		case <-ticker.C:
	//			average := getCount - lastCount
	//			fmt.Println("过去一分钟完成:", average)
	//			if average > 0 {
	//				fmt.Println("预计需要", putCount/average, "分钟")
	//			}
	//			lastCount = getCount
	//		default:
	//
	//		}
	//	}
	//}()
	//go func() {
	//	chapterDao := dao.NewChapter()
	//	for id := range markCh {
	//		rid, err := chapterDao.MakeAsDone(id)
	//		fmt.Println("get.task", id, rid, err)
	//
	//	}
	//}()
	//downloadBooks(client)

	//select {}
	//serve()
	createPrimaryColor()
}

func serve() {
	http.Handle("/", http.FileServer(http.Dir(fmt.Sprintf("%s/output", rootDir))))
	_ = http.ListenAndServe(":8809", nil)
}

func createPrimaryColor() {
	_ = filepath.Walk(fmt.Sprintf("%s/output", rootDir), func(path string, info os.FileInfo, err error) error {
		if info.Name() != "cover.jpg" {
			return nil
		}
		err, output := shell.Pipe("sh", "primary_color.sh", path)
		if err == nil {
			if strings.Index(output, "srgb") >= 0 {
				ret := strings.TrimSpace(output)
				bookId := filepath.Base(filepath.Dir(path))
				color := parseRgba(ret)
				fmt.Println(bookId, color)
				dao.NewBook().UpdateBookColorById(bookId, color)
			}
		}
		return nil
	})
}

func parseRgba(input string) string {
	if len(input) == 0 {
		return ""
	}
	input = strings.Trim(input, "srgb")
	input = strings.Trim(input, "(")
	input = strings.Trim(input, ")")
	parts := strings.Split(input, ",")
	var rgba []string
	for _, p := range parts {
		pure := strings.Trim(p, "%")
		score, _ := strconv.ParseFloat(pure, 10)
		rgba = append(rgba, fmt.Sprintf("%.f", score*256/100))
	}
	if len(rgba) == 3 {
		rgba = append(rgba, "1")
	}
	return "rgba(" + strings.Join(rgba, ",") + ")"
}

func populateIsNot() {
	bookFile := rootDir + "/books.json"
	if fs.FileExists(bookFile) {
		return
	}

	var bookList []model.Book
	for i := 1; i <= 55; i++ {
		targetLink := fmt.Sprintf(bookLink, i)
		pageCotent, err := net.FetchResponse(http.MethodGet, targetLink, nil, map[string]string{
			"referer":    "https://readnovelfull.com/latest-release-novel",
			"user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36",
		}, 3)
		if err != nil {
			log.Fatal(err)
		}
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageCotent))
		if err != nil {
			log.Fatal(err)
		}
		doc.Find(`div.list-novel`).Find(`div.row`).Each(func(i int, selection *goquery.Selection) {
			imageLink, ok := selection.Find(`img`).Attr(`src`)
			if !ok {
				return
			}
			detailLink, _ := selection.Find(`h3.novel-title a`).Attr(`href`)
			detailLink = strings.TrimSpace(detailLink)
			if len(imageLink) == 0 || len(detailLink) == 0 {
				log.Fatal("imageLink.title.detailLink.empty.author")
			}
			var isHot bool
			if selection.Find(`span.label-hot`).Length() > 0 {
				isHot = true
			}
			var isNew bool
			if selection.Find(`span.label-new`).Length() > 0 {
				isNew = true
			}
			id := encode.CrcEncode(detailLink)
			dao.NewBook().UpdateIsHotIsNew(id, isHot, isNew)
			fmt.Println("update.hot.new", id, isHot, isNew)
		})
	}
	buf, _ := json.Marshal(bookList)
	_ = ioutil.WriteFile(bookFile, buf, 0644)
}

func saveGenres() {
	buf, err := ioutil.ReadFile(fmt.Sprintf("%s/output/genre.html", rootDir))
	if err != nil {
		log.Fatal(err)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(buf))
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(`li`).Each(func(i int, selection *goquery.Selection) {
		name := strings.TrimSpace(selection.Text())
		stmt, err := database.GetInstance().Prepare(`insert into novel_genre(name, genre_id) values(?,?)`)
		if err != nil {
			log.Fatal(err)
		}
		defer func() { _ = stmt.Close() }()
		ret, err := stmt.Exec(name, encode.CrcEncode(name))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(ret.LastInsertId())
	})
}

func downloadBooks(client *downloader.Client) {
	stmt, err := database.GetInstance().Prepare(`select count(*) from novel_chapter 
		where chapter_index <= 200 and done = 0`)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = stmt.Close() }()
	var count int64
	_ = stmt.QueryRow().Scan(&count)
	fmt.Println("TOTAL.COUNT", count)

	stmt, err = database.GetInstance().Prepare(`select id, book_id, token, src_link, 
       name, chapter_index from novel_chapter where chapter_index <= 200 and done = 0`)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = stmt.Close() }()
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}

	items := make([]model.Chapter, 0)
	for rows.Next() {
		var item model.Chapter
		err = rows.Scan(&item.Id, &item.BookId, &item.Token, &item.Link, &item.Name, &item.Index)
		if err == nil {
			items = append(items, item)
		}
	}
	for _, item := range items {
		atomic.AddUint32(&putCount, 1)
		downloadBookTask(client, item, false)
		fmt.Println("progress", count, putCount, getCount)
	}
}

func downloadBookTask(client *downloader.Client, item model.Chapter, override bool) {

	client.PutTask(func() {
		target := fmt.Sprintf("%s/output/%s/chs/%d_%s.html", rootDir, item.BookId, item.Index, item.Token)
		if fs.FileExists(target) && !override {
			markCh <- item.Id
			atomic.AddUint32(&getCount, 1)
			return
		}
		dir := filepath.Dir(target)
		if !fs.FileExists(dir) {
			_ = os.MkdirAll(dir, 0755)
		}
		page, err := net.FetchResponse(http.MethodGet, srcHost+item.Link, nil, map[string]string{}, 3)
		if err != nil {
			return
		}
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(page))
		if err != nil {
			return
		}
		ret, err := doc.Find(`#chr-content`).Html()
		if err != nil {
			return
		}

		content := strings.ReplaceAll(chapterTpl, "{{title}}", item.Name)
		content = strings.ReplaceAll(content, "{{body}}", ret)
		doc1, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(content)))
		if err != nil {
			return
		}
		doc1.Find(`div.ads-holder`).Remove()
		ret1, err := doc1.Html()
		if err != nil {
			fmt.Println("remove", err.Error())
			return
		}
		_ = ioutil.WriteFile(target, []byte(str.CleanText(ret1)), 0644)
		fmt.Println("write.file", target)
		atomic.AddUint32(&getCount, 1)
		markCh <- item.Id
	})

}

func populateGenres() {

	stmt, err := database.GetInstance().Prepare(`select book_id, src_link from novel_book where genre_id = ''`)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = stmt.Close() }()
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}

	items := make([]model.Book, 0)
	updates := make(map[string]string)
	for rows.Next() {
		var booId, link string
		err = rows.Scan(&booId, &link)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, model.Book{Id: booId, Link: link})
	}

	for i, book := range items {
		absDetailLink := srcHost + book.Link
		fmt.Println("start.book", absDetailLink)
		genre, _, _, _, _, _ := detail(absDetailLink)
		fmt.Println("ready.book", absDetailLink)
		if len(genre) > 0 {
			var genreIds []string
			genreDao := dao.NewGenre()
			genres := strings.Split(genre, ",")
			for _, name := range genres {
				name = strings.TrimSpace(name)
				gid, err := genreDao.FindIdByName(name)
				if err == nil {
					genreIds = append(genreIds, gid)
				} else {
					fmt.Println(err.Error())
				}
			}
			fmt.Println(genres, genreIds)
			if len(genreIds) > 0 {
				updates[book.Id] = strings.Join(genreIds, ",")
				dao.NewBook().UpdateGenreId(book.Id, updates[book.Id])
				fmt.Println("---------- DONE", len(items), i, book.Id, updates[book.Id])
			}
		}
	}
}

func saveBooks() {
	bookFile := rootDir + "/output/books.json"
	if !fs.FileExists(bookFile) {
		return
	}
	buf, err := ioutil.ReadFile(bookFile)
	if err != nil {
		log.Fatal(err)
	}
	var books []model.Book
	err = json.Unmarshal(buf, &books)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("TOTAL", len(books))
	for i, book := range books {
		if i < 1002 {
			continue
		}
		absDetailLink := srcHost + book.Link
		fmt.Println("start.book", absDetailLink)
		genre, status, alter, cover, desc, chs := detail(absDetailLink)
		fmt.Println("ready.book", absDetailLink)
		book.Genre = genre

		book.Status = status
		book.NameAlter = alter
		book.Brief = desc

		// 下载小图
		fmt.Println("start.download.thumb")
		thumbExt := filepath.Ext(book.Cover)
		target := fmt.Sprintf("%s/output/%s/thumb%s", rootDir, book.Id, thumbExt)
		_ = os.MkdirAll(filepath.Dir(target), 0755)
		downloadCover(book.Cover, target, false)
		fmt.Println("done.download.thumb")

		// 下载cover
		fmt.Println("start.download.cover")
		book.Cover = cover
		coverExt := filepath.Ext(cover)
		target = fmt.Sprintf("%s/output/%s/cover%s", rootDir, book.Id, coverExt)
		_ = os.MkdirAll(filepath.Dir(target), 0755)
		downloadCover(book.Cover, target, false)
		fmt.Println("done.download.cover")

		book.Chapters = chs
		book.ChaptersCount = len(chs)
		book.Cover = "cover" + coverExt
		_, err := dao.NewBook().Add(book)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("save.book", book.Id, book.Title)
		for _, ch := range chs {
			ch.BookId = book.Id
			dao.NewChapter().Save(ch)
			fmt.Println("save.chapter", ch.BookId, ch.Index)
			// TODO download
		}

		fmt.Println("DONE", i, "of", len(books))
	}
	fmt.Println("DONE")
}

func downloadCover(from, to string, override bool) {
	if fs.FileExists(to) && !override {
		return
	}
	image, err := net.FetchResponse(http.MethodGet, from, nil, map[string]string{
		"referer":    "https://readnovelfull.com/latest-release-novel",
		"user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36",
	}, 3)
	if err != nil {
		log.Fatal(err)
	}
	_ = ioutil.WriteFile(to, image, 0644)
}

func book(link, to string) {
	output, err := net.FetchResponse(http.MethodGet, link, nil, map[string]string{}, 3)
	if err != nil {
		log.Fatal(err)
	}
	content := str.CleanText(string(output))
	err = ioutil.WriteFile(to, []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func detail(link string) (string, string, string, string, string, []model.Chapter) {
	page, err := net.FetchResponse(http.MethodGet, link, nil, map[string]string{
		"referer":    "https://readnovelfull.com/latest-release-novel",
		"user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36",
	}, 3)
	if err != nil {
		log.Fatal(err)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(page))
	if err != nil {
		log.Fatal(err)
	}
	novelId, ok := doc.Find(`div#rating`).Attr(`data-novel-id`)
	if !ok {
		log.Fatal("novel.id.not.ok")
	}
	coverLink, ok := doc.Find(`div.book img`).Attr(`src`)
	if !ok {
		log.Fatal("cover.not.ok")
	}
	kv := make(map[string]string)
	doc.Find(`ul.info-meta`).Find(`li`).Each(func(i int, selection *goquery.Selection) {
		key := strings.TrimSpace(selection.Find(`h3`).Text())
		value := strings.Trim(strings.Replace(strings.TrimSpace(selection.Text()), key, "", 1), `:`)
		kv[strings.ToLower(key)] = strings.TrimSpace(value)
	})
	var genre string
	if v, ok := kv["genre:"]; ok {
		genre = v
	}
	var status string
	if v, ok := kv["status:"]; ok {
		status = v
	}
	var alter string
	if v, ok := kv["alternative names:"]; ok {
		alter = v
	}
	var desc []string
	doc.Find(`div.desc-text`).Find(`p`).Each(func(i int, selection *goquery.Selection) {
		desc = append(desc, strings.TrimSpace(selection.Text()))
	})
	items := chapters(novelId)
	return genre, status, alter, coverLink, strings.Join(desc, "|"), items
}

const chapterLink = "https://readnovelfull.com/ajax/chapter-archive?novelId=%s"

func chapters(nid string) []model.Chapter {
	output, err := net.FetchResponse(http.MethodGet, fmt.Sprintf(chapterLink, nid), nil, map[string]string{
		"authority":        "readnovelfull.com",
		"x-requested-with": "XMLHttpRequest",
		"user-agent":       "",
		"origin":           "https://readnovelfull.com/invincible.html",
		"referer":          "https://readnovelfull.com/invincible.html",
	}, 3)
	if err != nil {
		log.Fatal(err.Error())
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(output))
	if err != nil {
		log.Fatal(err)
	}

	var items []model.Chapter
	var index int
	doc.Find(`li`).Each(func(i int, selection *goquery.Selection) {
		a := selection.Find(`a`)
		chapterLink, ok := a.Attr(`href`)
		if !ok {
			log.Fatal(err)
		}
		chapterTitle := a.Text()
		index += 1
		items = append(items, model.Chapter{
			Index: index,
			Link:  chapterLink,
			Name:  chapterTitle,
		})
	})
	return items
}

func books() {
	bookFile := rootDir + "/books.json"
	if fs.FileExists(bookFile) {
		return
	}

	var bookList []model.Book
	bookIdMap := make(map[string]struct{})
	for i := 1; i <= 55; i++ {
		targetLink := fmt.Sprintf(bookLink, i)
		pageCotent, err := net.FetchResponse(http.MethodGet, targetLink, nil, map[string]string{
			"referer":    "https://readnovelfull.com/latest-release-novel",
			"user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36",
		}, 3)
		if err != nil {
			log.Fatal(err)
		}
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageCotent))
		if err != nil {
			log.Fatal(err)
		}
		doc.Find(`div.list-novel`).Find(`div.row`).Each(func(i int, selection *goquery.Selection) {
			imageLink, ok := selection.Find(`img`).Attr(`src`)
			if !ok {
				return
			}
			imageLink = strings.TrimSpace(imageLink)
			title := strings.TrimSpace(selection.Find(`h3.novel-title a`).Text())
			detailLink, _ := selection.Find(`h3.novel-title a`).Attr(`href`)
			detailLink = strings.TrimSpace(detailLink)
			author := strings.TrimSpace(selection.Find(`span.author`).Text())
			if len(imageLink) == 0 || len(title) == 0 || len(detailLink) == 0 {
				log.Fatal("imageLink.title.detailLink.empty.author", author)
			}
			var isHot bool
			if selection.Find(`span.label-hot`).Length() > 0 {
				isHot = true
			}
			var isNew bool
			if selection.Find(`span.label-new`).Length() > 0 {
				isNew = true
			}
			if len(author) == 0 {
				author = "<anonymous>"
			}
			id := encode.CrcEncode(detailLink)
			if _, ok := bookIdMap[id]; ok {
				return
			}
			bookIdMap[id] = struct{}{}
			book := model.Book{
				Id:     id,
				Cover:  imageLink,
				Title:  title,
				Author: author,
				IsHot:  isHot,
				IsNew:  isNew,
				Link:   detailLink}
			fmt.Println(book.Title, book.Link)
			bookList = append(bookList, book)
		})
		fmt.Println("DONE", i, "TOTAL", len(bookList))
		time.Sleep(time.Second * 2)
	}
	buf, _ := json.Marshal(bookList)
	_ = ioutil.WriteFile(bookFile, buf, 0644)
}
