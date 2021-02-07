package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"inla/inla-crawler/libs/fs"
	"inla/inla-crawler/libs/net"
	"inla/inla-crawler/libs/shell"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

/*
破解方式：
1. 通过详情页，遍历ID
2. 视频格式是把图片改成mp4后缀
*/

const pageStart = `http://za.v-gals.com/detail?id=%d`
const outputRoot = `/Volumes/extend/crawler/mix/vgals`

func main() {
	for i := 1000; i < 2000; i++ {
		page := fmt.Sprintf(pageStart, i)
		body, err := net.FetchResponse(http.MethodGet, page, nil, map[string]string{}, 1)
		if err != nil {
			log.Fatal(err.Error())
		}
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			log.Fatal(err.Error())
		}
		title := doc.Find(`title`).Text()
		if strings.Contains(title, "ErrorException") {
			log.Println("____empty", page)
			continue
		}
		doc.Find(`#detail_main`).Find(`img`).Each(func(j int, selection *goquery.Selection) {
			src, ok := selection.Attr(`src`)
			if !ok {
				log.Fatal("missing.src")
			}

			uri, err := url.Parse(src)
			if err != nil {
				log.Fatal("parse.src")
			}

			var subDir string
			pathParts := strings.Split(uri.Path, "/")
			if strings.HasPrefix(pathParts[1], "video") {
				if len(pathParts) > 3 {
					subDir = pathParts[1] + "-" + pathParts[2]
				}
			}
			if len(subDir) == 0 {
				subDir = pathParts[1]
			}
			outputDir := filepath.Join(outputRoot, subDir)
			log.Println("____output", outputDir)
			if !fs.FileExists(outputDir) {
				_ = os.MkdirAll(outputDir, 0755)
			}
			downloadImage(i, src, subDir)
			log.Println("DONE", i, src)

			if strings.HasPrefix(pathParts[1], "video") {
				downloadVideo(i, src, subDir)
			}
		})
	}
}

func downloadImage(i int, src, subDir string) {
	uri, _ := url.Parse(src)
	filePath := filepath.Join(outputRoot, subDir, fmt.Sprintf("%d%s", i, path.Ext(uri.Path)))
	if fs.FileExists(filePath) {
		log.Println("EXISTS", filePath)
		return
	}
	body, err := net.FetchResponse(http.MethodGet, src, nil, nil, 3)
	if err != nil {
		log.Fatal("fetch.src", err.Error())
	}
	err = ioutil.WriteFile(filePath, body, 0644)
	if err != nil {
		log.Fatal("write.file")
	}
}

func downloadVideo(i int, src, subDir string) {
	filePath := filepath.Join(outputRoot, subDir, fmt.Sprintf("%d.%s", i, "mp4"))
	if fs.FileExists(filePath) {
		log.Println("EXISTS", filePath)
		return
	}
	src = strings.ReplaceAll(src, ".jpg", ".mp4")
	src = strings.ReplaceAll(src, ".png", ".mp4")
	body, err := net.FetchResponse(http.MethodGet, src, nil, nil, 3)
	if err != nil {
		log.Fatal("fetch.src")
	}
	err = ioutil.WriteFile(filePath, body, 0644)
	if err != nil {
		log.Fatal("write.file")
	}
	log.Println("DONE MP4", i, src)
}

func resizeVideo() {
	_ = filepath.Walk(outputRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !strings.HasPrefix(info.Name(), ".mp4") {
			return nil
		}
		err, output := shell.Pipe("sh", "video_wh.sh", path)
		if err != nil {
			log.Fatal("____video.wh", err.Error())
		}
		parts := strings.Split(output, "x")
		width, height := parts[0], parts[1]
		widthInt, _ := strconv.Atoi(width)
		heightInt, _ := strconv.Atoi(height)
		if widthInt > 640 {
			resizeHeight := int64(640 * heightInt / widthInt)
			if resizeHeight%2 != 0 {
				resizeHeight = resizeHeight - 1
			}
			thumb := path + ".sm.mp4"
			_, _ = shell.Pipe("sh", "resize_video.sh", path, fmt.Sprintf("640:%d", resizeHeight), thumb)
			if fs.FileExists(thumb) {
				_ = os.Remove(path)
				_ = os.Rename(thumb, path)
			}
		}
		return nil
	})
}
