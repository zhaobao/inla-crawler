package main

import (
	"bytes"
	"fmt"
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/libs/fs"
	"inla/inla-crawler/libs/shell"
	"inla/inla-crawler/tasks/comic/qgxymdmz/model"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// 漫画水印处理
const remoteDir = "/Volumes/extend/crawler/comic/imgs"
const rootDir = "factory/watermark/comic"
const dbLink = "tasks/comic/qgxymdmz/db.sqlite"

func setup() {
	database.Connect(dbLink)
}

func main() {
	setup()
	show()
	//removeUselessFile()
}

func removeUselessFile() {
	rows, err := database.GetInstance().Query(`select c.id, c.comic_id, c.chapter_index, c.chapter_part_index, 
		c.token from comic_chapter as c left join comic_book as b on c.comic_id = b.comic_id WHERE b.dirty = 0 
		                                                                                       AND c.clean = 1 
		                                                                                       AND c.chapter_index <= 20 order by c.id`)
	if err != nil {
		log.Fatal(err)
	}
	imgs := make(map[string]struct{})
	var items []model.ChapterItem
	for rows.Next() {
		var item model.ChapterItem
		err = rows.Scan(&item.Id, &item.ComicId, &item.ChapterIndex, &item.ChapterPartIndex, &item.Token)
		if err == nil {
			item.ImgLink = fmt.Sprintf("%s/input/%s/%.f_%d_%s.jpg", rootDir, item.ComicId, item.ChapterIndex, item.ChapterPartIndex, item.Token)
			items = append(items, item)
			imgs[item.ImgLink] = struct{}{}
		}
	}
	_ = filepath.Walk(fmt.Sprintf("%s/input", rootDir), func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".jpg") {
			return nil
		}
		if _, ok := imgs[path]; !ok {
			_ = os.Remove(path)
		}
		return nil
	})
	fmt.Println(len(items))
}

func show() {

	http.HandleFunc("/img", func(writer http.ResponseWriter, request *http.Request) {
		name := request.URL.Query().Get("name")
		fpath, err := url.PathUnescape(name)
		if err != nil {
			log.Fatal(err.Error())
		}
		if !fs.FileExists(fpath) {
			log.Fatal("not", name)
		}
		buf, _ := ioutil.ReadFile(fpath)
		_, _ = writer.Write(buf)
	})
	http.HandleFunc("/books", func(writer http.ResponseWriter, request *http.Request) {
		rows, err := database.GetInstance().Query(`select comic_id from comic_book where dirty = 0`)
		if err != nil {
			log.Fatal(err)
		}
		var body bytes.Buffer
		for rows.Next() {
			var id string
			err := rows.Scan(&id)
			if err == nil {
				body.WriteString(fmt.Sprintf("<a href='/chapters?book_id=%s'>%s</a><br>", id, id))
			}
		}
		writer.Header().Set("content-type", "text/html;charset=utf-8")
		_, _ = writer.Write(body.Bytes())

	})
	http.HandleFunc("/action/record", func(writer http.ResponseWriter, request *http.Request) {
		appendFile(fmt.Sprintf("%s/output/action.txt", rootDir), fmt.Sprintf(request.URL.Query().Encode()+"\n"))

		query := request.URL.Query()
		action := query.Get("action")
		if action == "remove" {
			id := query.Get("id")
			if len(id) > 0 {
				removeOne(id)
			}
		} else if action == "add" {
			size := query.Get("s")
			booId := query.Get("book_id")
			var watermark string
			if size == "b" {
				watermark = "factory/watermark/comic/output/water.2.png"
			} else {
				watermark = "factory/watermark/comic/output/water.5.png"
			}
			positions := `northwest,northeast,southwest,southeast`
			offset := `+0+0`
			rows, _ := database.GetInstance().Query(`select id, comic_id, chapter_index, chapter_part_index,token from comic_chapter where comic_id = ?`, booId)
			var items []model.ChapterItem
			for rows.Next() {
				var item model.ChapterItem
				err := rows.Scan(&item.Id, &item.ComicId, &item.ChapterIndex, &item.ChapterPartIndex, &item.Token)
				if err == nil {
					item.ImgLink = fmt.Sprintf("%s/input/%s/%.f_%d_%s.jpg", rootDir, item.ComicId, item.ChapterIndex, item.ChapterPartIndex, item.Token)
					items = append(items, item)
				}
			}
			for _, item := range items {
				if fs.FileExists(item.ImgLink) {
					watermarkWorker(item.ImgLink, watermark, positions, offset)
					updateClean(item.Id)
				}
			}
		}

		_, _ = writer.Write([]byte("OK"))
	})
	http.HandleFunc("/chapters", func(writer http.ResponseWriter, request *http.Request) {
		rows, err := database.GetInstance().Query(`select id, comic_id, chapter_index, chapter_part_index, 
       		token from comic_chapter where comic_id = ? and chapter_index <= 20`, request.URL.Query().Get("book_id"))
		if err != nil {
			log.Fatal(err)
		}
		var body bytes.Buffer
		for rows.Next() {
			var item model.ChapterItem
			err := rows.Scan(&item.Id, &item.ComicId, &item.ChapterIndex, &item.ChapterPartIndex, &item.Token)
			if err == nil {
				item.ImgLink = fmt.Sprintf("%s/input/%s/%.f_%d_%s.jpg", rootDir, item.ComicId, item.ChapterIndex, item.ChapterPartIndex, item.Token)
				body.WriteString(fmt.Sprintf(`
<div style='box-sizing:border-box;padding:0.5em;border:1px solid red;position:relative;float:left;width:16.6666667%%'>
	<img style='width:100%%' src='/img?name=%s'/>
	<p style='background-color:#9C27B0;margin:0;padding:8px;text-align:right;'>
	<button onclick='addwaterbookbig("%s")' id='book-b-%s'>大水印</button>
	<button onclick='addwaterbooksmall("%s")' id='book-s-%s'>小水印</button>
	<button onclick='remove(%d)' id='remove-%d'>删除这个</button>
	</p>
</div>
`, item.ImgLink, item.ComicId, item.ComicId, item.ComicId, item.ComicId, item.Id, item.Id))
			}
		}
		body.WriteString(`
<script>
function addwaterbooksmall(bookid) {
	let e = document.getElementById('book-s-' + bookid);
	e.style.backgroundColor = 'red';
	ajax('action=add&size=s&book_id=' + bookid, function() {
		e.style.removeProperty('background-color');
	});
}
function addwaterbookbig(bookid) {
	let e = document.getElementById('book-b-' + bookid);
	e.style.backgroundColor = 'red';
	ajax('action=add&size=b&book_id=' + bookid, function() {
		e.style.removeProperty('background-color');
	});
}
function remove(id) {
	let e = document.getElementById('remove-' + id);
	e.style.backgroundColor = 'red';
	ajax('action=remove&id=' + id, function() {
		e.style.removeProperty('background-color');
		e.parentNode.parentNode.style.opacity = '0.1';
	});
}
function ajax(q, cb) {
	let request = new XMLHttpRequest();
	request.open('GET', '/action/record?' + q, true);
	request.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded; charset=UTF-8');
	request.onload = function () {
		if (this.status >= 200 && this.status < 400) {
			cb(this.response);
		} else {
			alert('?');
		}
	}
	request.send();
}
</script>
`)
		writer.Header().Set("content-type", "text/html;charset=utf-8")
		_, _ = writer.Write(body.Bytes())
	})
	_ = http.ListenAndServe(":22222", nil)

}

func removeOne(id string) {
	_, err := database.GetInstance().Exec(`delete from comic_chapter where id = ?`, id)
	fmt.Println(id, err)
}

func process() {
	setup()
	// 处理book的前20chapter
	rows, err := database.GetInstance().Query(`select c.id, c.comic_id, c.chapter_index, c.chapter_part_index, 
		c.token from comic_chapter as c left join comic_book as b on c.comic_id = b.comic_id WHERE b.dirty = 0 
		                                                                                       AND c.clean = 1 
		                                                                                       AND c.chapter_index <= 20 order by c.id`)
	if err != nil {
		log.Fatal(err)
	}

	// 缓存数据
	var items []model.ChapterItem
	for rows.Next() {
		var item model.ChapterItem
		err = rows.Scan(&item.Id, &item.ComicId, &item.ChapterIndex, &item.ChapterPartIndex, &item.Token)
		if err == nil {
			item.ImgLink = fmt.Sprintf("%s/input/%s/%.f_%d_%s.jpg", rootDir, item.ComicId, item.ChapterIndex, item.ChapterPartIndex, item.Token)
			items = append(items, item)
		}
	}
	log.Println("total", len(items))

	// 准备数据
	for _, item := range items {
		remoteLink := fmt.Sprintf("%s/%s/%.f_%d_%s.jpg", remoteDir, item.ComicId, item.ChapterIndex, item.ChapterPartIndex, item.Token)
		if !fs.FileExists(remoteLink) {
			log.Fatal(remoteLink, "not.exists")
		}
		if !fs.FileExists(item.ImgLink) {
			_ = os.MkdirAll(filepath.Dir(item.ImgLink), 0755)
			fmt.Println("cp", remoteLink, item.ImgLink)
			_, _ = shell.Pipe("cp", remoteLink, item.ImgLink)
		}
	}

	done := make(map[int64]struct{})

	http.HandleFunc("/next/book", func(writer http.ResponseWriter, request *http.Request) {
		booId := request.URL.Query().Get("book_id")
		for _, item := range items {
			if item.ComicId == booId {
				done[item.Id] = struct{}{}
				fmt.Println("/next/book", booId, item.Id)
			}
		}
		_, _ = writer.Write([]byte("ok"))
	})
	http.HandleFunc("/resize/cut", func(writer http.ResponseWriter, request *http.Request) {
		booId := request.URL.Query().Get("book_id")
		part := request.URL.Query().Get("part")
		size := request.URL.Query().Get("size")
		for _, item := range items {
			if item.ComicId == booId {
				done[item.Id] = struct{}{}
				cutPartImage(item.ImgLink, part, size)
				updateClean(item.Id)
			}
		}
		_, _ = writer.Write([]byte("ok"))
	})
	http.HandleFunc("/dirty/book", func(writer http.ResponseWriter, request *http.Request) {
		booId := request.URL.Query().Get("book_id")
		for _, item := range items {
			if item.ComicId == booId {
				done[item.Id] = struct{}{}
			}
		}
		updateDirty(booId)
		_, _ = writer.Write([]byte("ok"))
	})
	http.HandleFunc("/drop/this", func(writer http.ResponseWriter, request *http.Request) {
		id := request.URL.Query().Get("id")
		cid, _ := strconv.ParseInt(id, 10, 64)
		dropOne(cid)
		done[cid] = struct{}{}
		_, _ = writer.Write([]byte("ok"))
	})
	http.HandleFunc("/ignore/book", func(writer http.ResponseWriter, request *http.Request) {
		booId := request.URL.Query().Get("book_id")
		for _, item := range items {
			if item.ComicId == booId {
				done[item.Id] = struct{}{}
				updateClean(item.Id)
				fmt.Println("/ignore/book", booId, item.Id)
			}
		}
		_, _ = writer.Write([]byte("ok"))
	})
	http.HandleFunc("/ignore/item", func(writer http.ResponseWriter, request *http.Request) {
		id := request.URL.Query().Get("id")
		cid, _ := strconv.ParseInt(id, 10, 64)
		updateClean(cid)
		done[cid] = struct{}{}
		fmt.Println("/ignore/item", cid)
		_, _ = writer.Write([]byte("ok"))
	})
	http.HandleFunc("/apply/book", func(writer http.ResponseWriter, request *http.Request) {
		log.Println("/apply/book", request.URL.Query())
		query := request.URL.Query()
		booId := query.Get("book_id")
		watermark := strings.Replace(query.Get("watermark"), "/img?name=", "", 1)
		positions := query.Get("position")
		offset := query.Get("offset")
		for _, item := range items {
			if item.ComicId == booId {
				if _, ok := done[item.Id]; !ok {
					watermarkWorker(item.ImgLink, watermark, positions, offset)
					updateClean(item.Id)
					done[item.Id] = struct{}{}
				}
			}
		}
		_, _ = writer.Write([]byte("ok"))
	})
	http.HandleFunc("/img", func(writer http.ResponseWriter, request *http.Request) {
		name := request.URL.Query().Get("name")
		fpath, err := url.PathUnescape(name)
		if err != nil {
			log.Fatal(err.Error())
		}
		if !fs.FileExists(fpath) {
			log.Fatal("not", name)
		}
		buf, _ := ioutil.ReadFile(fpath)
		_, _ = writer.Write(buf)
	})
	http.HandleFunc("/cut", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		id, _ := strconv.ParseInt(query.Get("id"), 10, 64)
		watermark := strings.Replace(query.Get("watermark"), "/img?name=", "", 1)
		input := strings.Replace(query.Get("input"), "/img?name=", "", 1)
		positions := query.Get("position")
		offset := query.Get("offset")
		fmt.Println(id, watermark, input, positions, offset)
		watermarkWorker(input, watermark, positions, offset)
		updateClean(id)
		done[id] = struct{}{}
		_, _ = writer.Write([]byte("ok"))
	})
	http.HandleFunc("/id", func(writer http.ResponseWriter, request *http.Request) {
		id, _ := strconv.ParseInt(request.URL.Query().Get("id"), 10, 64)
		log.Println("id=", id)
		var target model.ChapterItem
		for _, item := range items {
			if _, ok := done[item.Id]; ok {
				continue
			}
			if item.Id > id {
				target = item
				break
			}
		}
		writer.Header().Set("content-type", "text/html;charset=utf-8")
		if target.Id == 0 {
			_, _ = writer.Write([]byte("<h1>DONE</h1>"))
		} else {
			buf, err := ioutil.ReadFile(fmt.Sprintf("%s/output/page.html", rootDir))
			if err != nil {
				log.Fatal("read.page", err)
			}
			q := url.Values{}
			q.Set("name", target.ImgLink)
			w, h := imgSize(target.ImgLink)
			output := strings.ReplaceAll(string(buf), "{{_Id_}}", fmt.Sprintf("%d", target.Id))
			output = strings.ReplaceAll(output, "{{_Img_}}", fmt.Sprintf("/img?%s", q.Encode()))
			output = strings.ReplaceAll(output, "{{_BookId_}}", target.ComicId)
			output = strings.ReplaceAll(output, "{{_w_}}", fmt.Sprintf("%d", w))
			output = strings.ReplaceAll(output, "{{_h_}}", fmt.Sprintf("%d", h))
			_, _ = writer.Write([]byte(output))
		}
	})
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("content-type", "text/html;charset=utf-8")
		_, _ = writer.Write([]byte("<h1>404</h1>"))
	})
	log.Println("http://127.0.0.1:8809/id?id=0")
	_ = http.ListenAndServe(":8809", nil)
}

type PageData struct {
	Item model.ChapterItem
	Id   int64
}

// 存储
func updateClean(id int64) {
	_, err := database.GetInstance().Exec(`update comic_chapter set clean = 1 where id = ?`, id)
	if err != nil {
		log.Fatal("update.clean", err.Error())
	}
	log.Println("++++++ update.clean", id)
}

// 删除
func dropOne(id int64) {
	appendFile(fmt.Sprintf("%s/output/drop.txt", rootDir), fmt.Sprintf(`delete from comic_chapter where id = %d`, id))
}

func updateDirty(bookId string) {
	_, err := database.GetInstance().Exec(`update comic_book set dirty = 1 where comic_id = ?`, bookId)
	if err != nil {
		log.Fatal("update.clean")
	}
	log.Println("update.dirty", bookId)
}

// 水印处理工厂
// 目前提供能力：
// 1. 覆盖水印，支持传入文件，制定位置覆盖
func watermarkWorker(input, watermark, positions, geometry string) string {
	// composite -gravity SouthEast -geometry +0+3 watermark.png input.jpg output.jpg
	geometry = strings.ReplaceAll(geometry, " ", "+")
	msg := fmt.Sprintf("%s::%s::%s::%s\n", input, watermark, positions, geometry)
	appendFile(fmt.Sprintf("%s/output/log.txt", rootDir), msg)
	for _, position := range strings.Split(positions, ",") {
		err, output := shell.Pipe("composite", "-gravity", position, "-geometry", geometry, watermark, input, input)
		if err != nil {
			log.Fatal("shell.pipe.composite", err)
		}
		log.Println(err, output)
	}
	return input
}

func cutPartImage(input, part, size string) {
	width, height := imgSize(input)
	crop, _ := strconv.Atoi(size)
	var cmd string
	if part == "T" {
		cmd = fmt.Sprintf("%dx%d+0+%s", width, height-crop, size)
	} else {
		cmd = fmt.Sprintf("%dx%d+0+0", width, height-crop)
	}
	_, _ = shell.Pipe("convert", "-crop", cmd, input, input)
}

// 拼接日志
func appendFile(input, msg string) {
	fd, err := os.OpenFile(input, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal("open.file", input, err)
	}
	_, err = fd.WriteString(msg)
	if err != nil {
		log.Fatal("write.file", input, err)
	}
	_ = fd.Close()
}

func imgSize(input string) (int, int) {
	_, out := shell.Pipe("identify", "-format", "%[fx:w]x%[fx:h]", input)
	out = strings.TrimSpace(out)
	wh := strings.Split(out, "x")
	w, _ := strconv.Atoi(wh[0])
	h, _ := strconv.Atoi(wh[1])
	return w, h
}
