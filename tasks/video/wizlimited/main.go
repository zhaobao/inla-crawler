package main

import (
	"encoding/json"
	"fmt"
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/libs/encode"
	"inla/inla-crawler/libs/number"
	"inla/inla-crawler/libs/str"
	"inla/inla-crawler/tasks/video/wizlimited/constant"
	"inla/inla-crawler/tasks/video/wizlimited/dao"
	"inla/inla-crawler/tasks/video/wizlimited/model"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const rootDir = "tasks/video/wizlimited"

func main() {
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
