package main

import (
	"encoding/json"
	"fmt"
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/libs/encode"
	"inla/inla-crawler/libs/fs"
	"inla/inla-crawler/libs/shell"
	"inla/inla-crawler/tasks/music/netease/ng/dao"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const dbDir = `tasks/music/netease/us`
const rootDir = `tasks/music/netease/ng`
const srcDir = `/Volumes/extend/crawler/music/163/ng/Nigeria`

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

func parse() {
	var items []dao.MusicBo
	dataBuf, _ := ioutil.ReadFile(filepath.Join(rootDir, "data.json"))
	_ = json.Unmarshal(dataBuf, &items)

	dbLink := database.GetInstance()
	for _, music := range items {

		rawFile := filepath.Join(srcDir, music.Name, music.Name+".mp3")
		if !fs.FileExists(rawFile) {
			fmt.Println("raw.file.not.exists")
			continue
		}

		music.ResId = encode.CrcEncode(fmt.Sprintf("%d::%s", music.Id, music.Name))
		music.ResLink = filepath.Join(music.ResId, "m.mp3")
		localAbsPath := filepath.Join(rootDir, "output", music.ResLink)
		if !fs.FileExists(filepath.Dir(localAbsPath)) {
			_ = os.MkdirAll(filepath.Dir(localAbsPath), 0755)
		}
		cpFile(filepath.Join(srcDir, music.Name, music.Name+".mp3"), localAbsPath)
		music.CoverLink = filepath.Join(music.ResId, "m.jpg")
		cpFile(filepath.Join(srcDir, music.Name, music.Name+".jpg"), filepath.Join(rootDir, "output", music.CoverLink))
		music.LyricLink = filepath.Join(music.ResId, "m.txt")
		cpFile(filepath.Join(srcDir, music.Name, music.Name+".txt"), filepath.Join(rootDir, "output", music.LyricLink))
		music.PrimaryColor = primaryColor(filepath.Join(rootDir, "output", music.CoverLink))

		id, err := dbLink.Exec(`insert into music_163_ng_v1(res_id, title, sub_title, cover_link, res_link, 
			lyric_link, category, cate_id, primary_color, duration, album, artists) values(?,?,?,?,?,?,?,?,?,?,?,?)`, music.ResId,
			music.Name, "", music.CoverLink, music.ResLink, music.LyricLink, "", "", music.PrimaryColor, music.Duration,
			music.Album, strings.Join(music.Artists, "|"))
		log.Println(id, err)
	}
}

func cpFile(from, to string) {
	if fs.FileExists(to) {
		log.Println("to.exists")
		return
	}
	if !fs.FileExists(from) {
		log.Println("from.not.exists", from)
		return
	}
	log.Println("cp", from, to)
	_, _ = shell.Pipe("cp", from, to)
}

func primaryColor(path string) string {
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
