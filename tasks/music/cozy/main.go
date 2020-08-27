package main

import (
	"encoding/json"
	"fmt"
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/libs/encode"
	"inla/inla-crawler/libs/fs"
	"inla/inla-crawler/libs/net"
	"inla/inla-crawler/libs/shell"
	"inla/inla-crawler/tasks/music/cozy/dao"
	"inla/inla-crawler/tasks/music/cozy/model"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//curl 'http://content.mcdmobi.com/api/cozy/list?&pageNo=1&pageSize=6' \
//  -H 'Connection: keep-alive' \
//  -H 'Accept: application/json, text/javascript, */*; q=0.01' \
//  -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36' \
//  -H 'Origin: http://cozyease.com' \
//  -H 'Referer: http://cozyease.com/index.html' \
//  -H 'Accept-Language: zh-CN,zh;q=0.9,en-SG;q=0.8,en;q=0.7' \
//  --compressed \
//  --insecure

const rootDir = "tasks/music/cozy"

func main() {
	database.Connect(fmt.Sprintf("%s/db.sqlite", rootDir))
	//serve()
	//createPrimaryColor()
	createDuration()
}

func serve() {
	http.Handle("/", http.FileServer(http.Dir(fmt.Sprintf("%s/output", rootDir))))
	_ = http.ListenAndServe(":8809", nil)
}

func createDuration() {
	_ = filepath.Walk(fmt.Sprintf("%s/output/res", rootDir), func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".mp3") {
			return nil
		}
		err, output := shell.Pipe("sh", "duration.sh", path)
		if err == nil {
			ret := strings.TrimSpace(output)
			duration, _ := strconv.ParseFloat(ret, 10)
			res := fmt.Sprintf("/res/" + info.Name())
			dao.Service.UpdateDurationByRes(res, duration)
		}
		return nil
	})
	// convert input.jpg -scale 1x1\! -format '%[pixel:u]' info:-
	// srgb(30.7705%,32.2335%,24.503%)%
}

func createPrimaryColor() {
	_ = filepath.Walk(fmt.Sprintf("%s/output/cover", rootDir), func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".png") {
			return nil
		}
		err, output := shell.Pipe("sh", "primary_color.sh", path)
		if err == nil {
			if strings.Index(output, "srgb") >= 0 {
				ret := strings.TrimSpace(output)
				color := parseRgba(ret)
				cover := fmt.Sprintf("/cover/" + info.Name())
				fmt.Println(cover, color)
				dao.Service.UpdateColorByCover(cover, color)
			}
		}
		return nil
	})
	// convert input.jpg -scale 1x1\! -format '%[pixel:u]' info:-
	// srgb(30.7705%,32.2335%,24.503%)%
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

func prepare() {
	var items []model.Cozy
	no := 1
	size := 6
	for {
		output, err := download(no, size)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Println("------", "pageNo", no, "pageSize", size)
		var item resp
		err = json.Unmarshal(output, &item)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		for _, v := range item.Data.ModelList {
			v.CoverHash = encode.CrcEncode(v.IconUrl)
			coverName := fmt.Sprintf("cover/%s.png", v.CoverHash)
			downloadRes(v.IconUrl, fmt.Sprintf("%s/output/%s", rootDir, coverName), false)
			v.CoverLink = "/" + coverName

			v.ResHash = encode.CrcEncode(v.AudioUrl)
			resName := fmt.Sprintf("res/%s.mp3", v.ResHash)
			downloadRes(v.AudioUrl, fmt.Sprintf("%s/output/%s", rootDir, resName), false)
			v.ResLink = "/" + resName

			items = append(items, v)
			id, err := dao.Service.SaveOrUpdate(v)
			fmt.Println("save", id, err)
		}
		if len(item.Data.ModelList) == 0 {
			break
		}
		no += 1
		time.Sleep(time.Second)
	}

	target := fmt.Sprintf("%s/output/cozy.json", rootDir)
	if !fs.FileExists(target) {
		buf, _ := json.Marshal(items)
		_ = ioutil.WriteFile(target, buf, 0644)
	}
	fmt.Println("DONE")
}

type resp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		PageNo     int          `json:"pageNo"`
		PageSize   int          `json:"pageSize"`
		PageCount  int          `json:"pageCount"`
		TotalCount int          `json:"totalCount"`
		ModelList  []model.Cozy `json:"modelList"`
	} `json:"data"`
}

func download(no, size int) ([]byte, error) {
	query := url.Values{}
	query.Set("pageNo", strconv.Itoa(no))
	query.Set("pageSize", strconv.Itoa(size))
	link := fmt.Sprintf(`http://content.mcdmobi.com/api/cozy/list?%s`, query.Encode())
	return net.FetchResponse(http.MethodGet, link, nil, map[string]string{
		"User-Agent": "Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Mobile Safari/537.36",
		"Origin":     "http://cozyease.com",
		"Host":       "content.mcdmobi.com",
		"Referer":    "http://cozyease.com/index.html",
	}, 3)
}

func downloadRes(from, to string, override bool) {
	if fs.FileExists(to) && !override {
		return
	}
	output, err := net.FetchResponse(http.MethodGet, from, nil, nil, 3)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = ioutil.WriteFile(to, output, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("DOWNLOAD", from, to, "SUCCESS")
}
