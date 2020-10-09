package main

import (
	"encoding/json"
	"fmt"
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/libs/downloader"
	"inla/inla-crawler/libs/encode"
	"inla/inla-crawler/libs/fs"
	"inla/inla-crawler/libs/net"
	"inla/inla-crawler/libs/number"
	"inla/inla-crawler/libs/shell"
	"inla/inla-crawler/libs/str"
	"inla/inla-crawler/tasks/video/wizlimited/constant"
	"inla/inla-crawler/tasks/video/wizlimited/dao"
	"inla/inla-crawler/tasks/video/wizlimited/model"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const rootDir = `tasks/video/wizlimited`
const outputRoot = `/Volumes/extend/crawler/video/wiz`
const downloaderPoolSize = 16

func init() {
	database.Connect(fmt.Sprintf("%s/db.sqlite", rootDir))
}

func main() {
	//download()
	serve()
	//thumb()
	//optimize()
}

func thumb() {
	_ = filepath.Walk(outputRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".mp4") {
			return nil
		}
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		thumbPath := path + ".jpg"
		if !fs.FileExists(thumbPath) {
			args := []string{"-i", path, "-ss", "00:00:01", "-vframes", "1", thumbPath}
			_, _ = shell.Pipe("ffmpeg", args...)
			log.Println("MISSING", path)
		}
		return nil
	})
}

func serve() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		path := filepath.Join(outputRoot, request.URL.Path)
		buf, _ := ioutil.ReadFile(path)
		fmt.Println(path)
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "POST,GET")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		_, _ = writer.Write(buf)
	})
	_ = http.ListenAndServe(":8809", nil)
}

func optimize() {
	start := 0
	limit := 10000
	querySql := `select res_id, type_id, group_id, raw_link from assets_video_4 order by id asc limit %d offset %d`
	rows, err := database.GetInstance().Query(fmt.Sprintf(querySql, limit, start))
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var resId, typeId, groupId, rawLink string
		err := rows.Scan(&resId, &typeId, &groupId, &rawLink)
		if err != nil {
			log.Fatal(err)
		}
		output1 := fmt.Sprintf("%s/%s/%s.mp4", outputRoot, groupId, resId)
		if fs.FileExists(output1) {
			output2 := fmt.Sprintf("%s/res/%s/%s.mp4", outputRoot, typeId, resId)
			outputDir := filepath.Dir(output2)
			if !fs.FileExists(outputDir) {
				_ = os.MkdirAll(outputDir, 0777)
			}
			_, _ = shell.Pipe("mv", output1, output2)
		} else {
			log.Fatal(output1, "not.exists")
		}
	}
	fmt.Println("DONE")
}

func download() {

	start := 0
	limit := 10000
	querySql := `select res_id, type_id, group_id, raw_link from assets_video_4 order by id asc limit %d offset %d`
	rows, err := database.GetInstance().Query(fmt.Sprintf(querySql, limit, start))
	if err != nil {
		log.Fatal(err)
	}

	var total, finished, diff uint32

	pool := downloader.New(downloaderPoolSize)
	pool.Start()
	done := make(chan struct{})

	for rows.Next() {
		var resId, typeId, groupId, rawLink string
		err := rows.Scan(&resId, &typeId, &groupId, &rawLink)
		if err != nil {
			log.Fatal(err)
		}
		output := fmt.Sprintf("%s/%s/%s.mp4", outputRoot, groupId, resId)
		if !fs.FileExists(filepath.Dir(output)) {
			err = os.MkdirAll(filepath.Dir(output), 0755)
			if err != nil {
				log.Fatal(err)
			}
		}
		if fs.FileExists(output) {
			atomic.AddUint32(&diff, 1)
			log.Println("continue", output, diff)
			continue
		}
		pool.PutTask(func() {
			atomic.AddUint32(&total, 1)
			defer atomic.AddUint32(&finished, 1)
			log.Println("PUT >>> ", rawLink)
			data, err := net.FetchResponse(http.MethodGet, rawLink, nil, nil, 3)
			if err != nil {
				log.Println(err)
				return
			}
			err = ioutil.WriteFile(output, data, 0644)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("🍺 TOTAL:", total, "FINISHED:", finished, "REMAINING:", 5240-total-diff)
			if (finished + diff) == total {
				done <- struct{}{}
			}
		})
	}
	<-done
	log.Println("DONE")
}

func prepare() {
	database.Connect(fmt.Sprintf("%s/db.sqlite", rootDir))
	buf, err := ioutil.ReadFile(fmt.Sprintf("%s/output/subscribe_assets_video.json", rootDir))
	if err != nil {
		log.Fatal(err)
	}

	var videos []model.VideoItem
	err = json.Unmarshal(buf, &videos)
	if err != nil {
		log.Fatal(err)
	}

	numberReg := regexp.MustCompile(`\d+`)

	var count int
	for _, video := range videos {
		if video.Watermark == constant.NoWaterMark {
			name, ok := constant.VideoTypeMap[video.Type]
			if !ok {
				log.Fatal(video.Type, "not.found")
			}

			if video.GIndex == 0 && len(video.GTitle) > 0 {
				if strings.HasPrefix(video.GTitle, "Day") {
					guessIndexStr := strings.TrimSpace(numberReg.FindString(video.GTitle))
					guessIndex, _ := strconv.ParseInt(guessIndexStr, 10, 64)
					video.GIndex = guessIndex
					fmt.Println("guess", guessIndex, video.GTitle)
				}
			}
			if video.GIndex == 0 && video.GId > 0 {
				video.GIndex = video.GId
			}

			item := model.VideoRow{
				ResId:         encode.CrcEncode(fmt.Sprintf("%d", video.SId)),
				ResIndex:      video.SIndex,
				ResTitle:      video.Title,
				ResLink:       encode.CrcEncode(str.RandStr(8)),
				GroupIndex:    video.GIndex,
				GroupTitle:    video.GTitle,
				VideoWidth:    video.Width,
				VideoHeight:   video.Height,
				VideoSize:     video.Size,
				VideoDuration: video.Duration,
				TypeId:        encode.CrcEncode(fmt.Sprintf("%d_%s", video.Type, name)),
				TypeName:      name,
				Source:        "wiz",
				CountPlay:     number.RandInt(100, 1000),
				CountLove:     number.RandInt(100, 1000),
				CountDown:     number.RandInt(100, 1000),
				CtTime:        number.RandInt64(time.Now().Unix(), time.Now().Unix()+10*86400),
				RawLink:       video.Link,
			}
			if video.GId > 0 {
				item.GroupId = encode.CrcEncode(fmt.Sprintf("%d", video.GId))
			}
			id, err := dao.NewVideo().Save(&item)
			fmt.Println(id, err)
		}
	}
	fmt.Println("count", count)
}
