package main

import (
	"encoding/json"
	"fmt"
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/libs/encode"
	"inla/inla-crawler/libs/fs"
	"inla/inla-crawler/libs/shell"
	"inla/inla-crawler/tasks/novel/gdsbook/dao"
	"inla/inla-crawler/tasks/novel/gdsbook/model"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	database.Connect(fmt.Sprintf("%s/db.sqlite", rootDir))
	serve()
	//categoryData()
	//renameCover()
}

const rootDir = "tasks/novel/gdsbook"

//type epubHandler struct {
//}
//
//func (e epubHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
//	fmt.Println(".......")
//	path := filepath.Join(rootDir, "output/data", request.URL.Path)
//	buf, _ := ioutil.ReadFile(path)
//	fmt.Println(path)
//	writer.Header().Set("Access-Control-Allow-Origin", "*")
//	writer.Header().Set("Access-Control-Allow-Methods", "POST,GET")
//	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
//	_, _ = writer.Write(buf)
//}

func serve() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(".......")
		path := filepath.Join(rootDir, "output/data", request.URL.Path)
		buf, _ := ioutil.ReadFile(path)
		fmt.Println(path)
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "POST,GET")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		_, _ = writer.Write(buf)
	})
	_ = http.ListenAndServe(":8809", nil)
}

func renameCover() {
	_ = filepath.Walk(fmt.Sprintf("%s/output/data", rootDir), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if info.Name() == "illustration.jpg" {
			_ = os.Rename(path, filepath.Join(filepath.Dir(path), "cover.jpg"))
		}
		return nil
	})
}

func categoryData() {

	genres := make(map[string]int)
	var books []model.BookItem

	_ = filepath.Walk(fmt.Sprintf("%s/output/data", rootDir), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".txt") {
			err := os.Remove(path)
			fmt.Println("remove", err)
		}
		if strings.HasSuffix(path, ".epub") && info.Name() != "book.epub" {
			_ = os.Rename(path, filepath.Join(filepath.Dir(path), "book.epub"))
			return nil
		}
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		return nil

		var cbooks []model.BookItem
		buf, _ := ioutil.ReadFile(path)
		_ = json.Unmarshal(buf, &cbooks)

		books = append(books, cbooks...)
		for _, book := range cbooks {
			if _, ok := genres[book.Category]; !ok {
				genres[book.Category] = 0
			}
			genres[book.Category] += 1
		}

		return nil
	})

	return

	for genre, count := range genres {
		gid, err := dao.NewGenre().Save(genre, encode.CrcEncode(genre), count)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(gid, err)
		oldCatePath := fmt.Sprintf("%s/output/data/%s", rootDir, genre)
		newCatePath := fmt.Sprintf("%s/output/data/%s", rootDir, gid)
		_ = os.Rename(oldCatePath, newCatePath)
	}
	for _, book := range books {
		gid, _ := dao.NewGenre().Find(book.Category)
		bookId := encode.CrcEncode(book.Title)
		bookNames := strings.Split(book.Illustration, "/")
		bookNames[0] = gid
		rawCoverImg := fmt.Sprintf("%s/output/data/%s", rootDir, strings.Join(bookNames, "/"))
		oldBookPath := filepath.Dir(rawCoverImg)
		bookNames[1] = bookId
		coverImg := fmt.Sprintf("%s/output/data/%s", rootDir, strings.Join(bookNames, "/"))
		newBookPath := filepath.Dir(coverImg)
		if oldBookPath != newBookPath {
			fmt.Println("do.rename")
			_ = os.Rename(oldBookPath, newBookPath)
		}
		if !fs.FileExists(coverImg) {
			log.Fatal("not.exists")
		}
		err, ret := shell.Pipe("sh", "primary_color.sh", coverImg)
		if err != nil {
			log.Fatal(err)
		}
		ret = strings.TrimSpace(ret)
		color := parseRgba(ret)
		bid, err := dao.NewBook().Save(book.Title, "cover.jpg", gid, bookId, "gdsbook", book.Url, color)
		fmt.Println(bid, err)

	}
	fmt.Println(genres)
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
