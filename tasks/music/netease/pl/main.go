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
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const dbDir = `tasks/music/netease/pl`
const rootDir = `tasks/music/netease/pl`
const srcDir = `/Volumes/extend/crawler/music/163/pl`

func main() {
	database.Connect(filepath.Join(dbDir, "db.sqlite"))
	parse()
}

func serve() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		path := filepath.Join(rootDir, "output", request.URL.Path)
		buf, _ := ioutil.ReadFile(path)
		fmt.Println(path)
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "POST,GET")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		_, _ = writer.Write(buf)
	})
	log.Println("http://127.0.0.1:8809")
	_ = http.ListenAndServe(":8809", nil)
}

type MusicBo struct {
	Id           int64    `json:"id"`
	Name         string   `json:"name"`
	Duration     int      `json:"duration"`
	Cover        string   `json:"cover"`
	Album        string   `json:"album"`
	Src          string   `json:"src"`
	LyricSrc     string   `json:"lyric_src"`
	Artists      []string `json:"artists"`
	Illustration string   `json:"illustration"`

	ResId        string `json:"res_id"`
	ResLink      string `json:"res_link"`
	CoverLink    string `json:"cover_link"`
	LyricLink    string `json:"lyric_link"`
	PrimaryColor string `json:"primary_color"`
}

func parse() {
	listAry := []string{
		"list_POLSKA1_info.json",
		"list_POLSKA2_info.json",
		"list_POLSKA3_info.json",
		"list_POLSKA4_info.json",
	}

	var tt int
	for _, listName := range listAry {
		var items []MusicBo
		dataBuf, err := ioutil.ReadFile(filepath.Join(srcDir, listName))
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(dataBuf, &items)
		if err != nil {
			log.Fatal(err)
		}

		dbLink := database.GetInstance()
		for _, music := range items {
			rawFile := filepath.Join(srcDir, music.Src)
			if !fs.FileExists(rawFile) {
				fmt.Println("raw.file.not.exists")
				log.Fatal("not.exists", rawFile)
			}

			tt += 1
			music.ResId = encode.CrcEncode(fmt.Sprintf("%d::%s", music.Id, music.Name))
			music.ResLink = filepath.Join(music.ResId, "m.mp3")
			localAbsPath := filepath.Join(rootDir, "output", music.ResLink)
			if !fs.FileExists(filepath.Dir(localAbsPath)) {
				_ = os.MkdirAll(filepath.Dir(localAbsPath), 0755)
			}
			cpFile(filepath.Join(srcDir, music.Src), localAbsPath)

			rawCover := filepath.Join(srcDir, music.Illustration)
			if len(music.Illustration) > 0 && fs.FileExists(rawCover) {
				music.CoverLink = filepath.Join(music.ResId, "m.jpg")
				cpFile(rawCover, filepath.Join(rootDir, "output", music.CoverLink))
				music.PrimaryColor = primaryColor(filepath.Join(rootDir, "output", music.CoverLink))
			} else {
				log.Println(music.Src, "no.cover", rawCover)
			}

			rawLyric := filepath.Join(srcDir, music.LyricSrc)
			if len(music.LyricSrc) > 9 && fs.FileExists(rawLyric) {
				music.LyricLink = filepath.Join(music.ResId, "m.txt")
				cpFile(rawLyric, filepath.Join(rootDir, "output", music.LyricLink))
			} else {
				log.Println(music.Src, "no.lyric")
			}

			id, err := dbLink.Exec(`insert into music_163_pl(res_id, title, sub_title, cover_link, res_link, 
			lyric_link, category, cate_id, primary_color, duration, album, artists) values(?,?,?,?,?,?,?,?,?,?,?,?)`,
				music.ResId, music.Name, "", music.CoverLink, music.ResLink,
				music.LyricLink, "", "", music.PrimaryColor, music.Duration,
				music.Album, strings.Join(music.Artists, "|"))
			log.Println(id, err)
		}
	}

	fmt.Println("------>", "tt", tt)
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

func primaryColor(path string) string {
	if !fs.FileExists(path) {
		log.Fatal("not.exists", path)
	}
	err, output := shell.Pipe("sh", "primary_color.sh", path)
	if err == nil {
		if strings.Index(output, "srgb") >= 0 {
			ret := strings.TrimSpace(output)
			return parseRgba(ret)
		}
	}
	return ""
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
