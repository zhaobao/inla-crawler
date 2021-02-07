package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"inla/inla-crawler/libs/databasex"
	"inla/inla-crawler/libs/fs"
	"inla/inla-crawler/libs/net"
	"inla/inla-crawler/libs/rand"
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
	log.Println("DONE")
}

func downloadAll() {
	for i := 694; i < 695; i++ {
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
		if !strings.HasSuffix(info.Name(), ".mp4") {
			return nil
		}
		err, output := shell.Pipe("sh", "video_wh.sh", path)
		if err != nil {
			log.Fatal("____video.wh", err.Error())
		}
		parts := strings.Split(output, "x")
		log.Println(parts, path)
		width, height := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		widthInt, _ := strconv.Atoi(width)
		heightInt, _ := strconv.Atoi(height)

		if widthInt > 640 {
			resizeHeight := int64(640 * heightInt / widthInt)
			if resizeHeight%2 != 0 {
				resizeHeight = resizeHeight - 1
			}
			log.Println(width, height, 640, resizeHeight)
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

func uploadWallpaper() {

	_ = filepath.Walk(filepath.Join(outputRoot, "wallpaper"), func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(info.Name(), ".jpg") && !strings.HasSuffix(info.Name(), ".png") {
			return nil
		}

		suffix := filepath.Ext(info.Name())
		//ID := rand.NewUUID(fmt.Sprintf("%s", path))
		//_ = os.Rename(path, filepath.Join(filepath.Dir(path), ID+suffix))

		var imageWidth, imageHeight int
		err, out := shell.Pipe("sh", "identify_size.sh", path)
		if err == nil {
			whParts := strings.Split(out, "x")
			imageWidth, _ = strconv.Atoi(whParts[0])
			imageHeight, _ = strconv.Atoi(whParts[1])
		}

		fileSize := info.Size()
		log.Println(path, imageWidth, imageHeight, fileSize)

		ID := strings.Split(info.Name(), ".")[0]
		item := ImageRow{
			SkuID:    ID,
			Link:     "https://a.inlamob.com/wallpaper/vgals/" + ID + suffix,
			Src:      "vgals",
			W:        imageWidth,
			H:        imageHeight,
			S:        fileSize,
			Category: "",
		}

		ret, err := databasex.GetInstance().Exec(`
		insert into in_assets_wallpaper(sku_id, link, src, w, h, s, category)
		values(?,?,?,?,?,?,?)`, item.SkuID, item.Link, item.Src, item.W, item.H, item.S, item.Category)
		if err != nil {
			log.Fatal(err)
		}
		id, err := ret.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("____", id)
		return nil
	})
}

type ImageRow struct {
	SkuID    string `db:"sku_id"`
	Link     string `db:"link"`
	Src      string `db:"src"`
	W        int    `db:"w"`
	H        int    `db:"h"`
	S        int64  `db:"s"`
	Category string `db:"category"`
}

func renameVideo() {
	_ = filepath.Walk(filepath.Join(outputRoot, "videos-beyond-cooking"), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !isImage(info.Name()) {
			return nil
		}
		videoPath := checkMatchVideo(path)
		err, output := shell.Pipe("sh", "video_wh.sh", videoPath)
		if err != nil {
			log.Fatal("____video.xy", err.Error())
		}

		size := fileSize(videoPath)
		suffix := filepath.Ext(path)
		parts := strings.Split(output, "x")
		width, height := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		widthInt, _ := strconv.Atoi(width)
		heightInt, _ := strconv.Atoi(height)
		imageID := rand.NewUUID(fmt.Sprintf("%s", path))
		videoID := rand.NewUUID(fmt.Sprintf("%s", reverseString(videoPath)))
		_ = os.Rename(path, filepath.Join(filepath.Dir(path), imageID+suffix))
		_ = os.Rename(videoPath, filepath.Join(filepath.Dir(path), videoID+".mp4"))
		item := VideoRow{
			SkuID:    videoID,
			Link:     "https://a.inlamob.com/video/vgals/" + videoID + ".mp4",
			Cover:    "https://a.inlamob.com/video/vgals/" + imageID + suffix,
			Src:      "vgals",
			W:        widthInt,
			H:        heightInt,
			S:        size,
			Category: "cooking",
		}

		ret, err := databasex.GetInstance().Exec(`
		insert into in_assets_video(sku_id, link, cover, src, w, h, s, category)
		values(?,?,?,?,?,?,?,?)`, item.SkuID, item.Link, item.Cover, item.Src, item.W, item.H, item.S, item.Category)
		if err != nil {
			log.Fatal(err)
		}
		id, err := ret.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("____", id)

		return nil
	})
}

type VideoRow struct {
	SkuID    string `db:"sku_id"`
	Link     string `db:"link"`
	Cover    string `db:"cover"`
	Src      string `db:"src"`
	W        int    `db:"w"`
	H        int    `db:"h"`
	S        int64  `db:"s"`
	Category string `db:"category"`
}

func isImage(name string) bool {
	return strings.HasSuffix(name, ".png") || strings.HasSuffix(name, ".jpg")
}

func checkMatchVideo(path string) string {
	dir, name := filepath.Split(path)
	imageName := strings.Split(name, ".")
	videoFile := filepath.Join(dir, imageName[0]+".mp4")
	if !fs.FileExists(videoFile) {
		log.Fatal("____no.match.video", path, videoFile)
	}
	return videoFile
}

func reverseString(input string) string {
	ary := strings.Split(input, "")
	var output []string
	for i := len(ary) - 1; i >= 0; i-- {
		output = append(output, ary[i])
	}
	return strings.Join(output, "")
}

func fileSize(input string) int64 {
	info, err := os.Stat(input)
	if err != nil {
		return 0
	}
	return info.Size()
}

/*
desc in_assets_wallpaper;
+----------+--------------+------+-----+-------------------+-----------------------------+
| Field    | Type         | Null | Key | Default           | Extra                       |
+----------+--------------+------+-----+-------------------+-----------------------------+
| id       | int(11)      | NO   | PRI | <null>            | auto_increment              |
| sku_id   | char(36)     | NO   |     | <null>            |                             |
| link     | varchar(128) | NO   |     | <null>            |                             |
| src      | varchar(12)  | NO   |     | <null>            |                             |
| w        | int(11)      | YES  |     | 0                 |                             |
| h        | int(11)      | YES  |     | 0                 |                             |
| s        | int(11)      | YES  |     | 0                 |                             |
| category | varchar(12)  | NO   |     | <null>            |                             |
| up_time  | timestamp    | YES  |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP |
+----------+--------------+------+-----+-------------------+-----------------------------+
*/
