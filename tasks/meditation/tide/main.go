package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"inla/inla-crawler/libs/database"
	"inla/inla-crawler/libs/fs"
	"inla/inla-crawler/libs/net"
	"inla/inla-crawler/libs/shell"
	"inla/inla-crawler/tasks/meditation/tide/dao"
	"inla/inla-crawler/tasks/meditation/tide/model"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const rootDir = "tasks/meditation/tide"

func main() {
	database.Connect(fmt.Sprintf("%s/db.sqlite", rootDir))
	//parseMusic()
	//parseMusicTag()
	//parseMeditation()
	//parseSleepMed()
	//populateGroup()
	//populateMusicGroup()
	serve()
}

func serve() {
	http.Handle("/", http.FileServer(http.Dir(fmt.Sprintf("%s/output", rootDir))))
	_ = http.ListenAndServe(":8809", nil)
}

func populateMusicGroup() {
	buf, err := ioutil.ReadFile(fmt.Sprintf("%s/music.json", rootDir))
	if err != nil {
		log.Fatal(err)
	}
	// 解析ID
	var music model.MusicItem
	err = json.Unmarshal(buf, &music)
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range music.Scenes {
		var groups []string
		for _, t := range m.TagsV2 {
			groups = append(groups, strings.ToLower(t.Key))
		}
		dao.MusicService.UpdateGroup(m.Id, strings.Join(groups, ","))
	}
}

func populateGroup() {
	buf, err := ioutil.ReadFile(fmt.Sprintf("%s/meditation.json", rootDir))
	if err != nil {
		log.Fatal(err)
	}
	var obj model.MeditationObj
	err = json.Unmarshal(buf, &obj)
	if err != nil {
		log.Fatal(err)
	}
	for _, album := range obj.Albums {
		var keys []string
		for _, tag := range album.TagsV2 {
			keys = append(keys, tag.Key)
		}
		dao.MedService.UpdateMedGroup(album.Id, strings.Join(keys, ","))
	}
}

func parseSleepMed() {
	buf, err := ioutil.ReadFile(fmt.Sprintf("%s/meditation.sleep.json", rootDir))
	if err != nil {
		log.Fatal(err)
	}
	var obj model.MeditationObj
	err = json.Unmarshal(buf, &obj)
	if err != nil {
		log.Fatal(err)
	}
	for _, album := range obj.Albums {
		albumKey := album.TagsV2[0].Key
		albumId := album.Id
		dao.MedService.UpdateMedGroup(albumId, albumKey)
	}
}

func parseMeditation() {
	// 预遍历meditation_res目录，获得hash的map
	files := make(map[string]string)
	_ = filepath.Walk(fmt.Sprintf("%s/output/tide/resources.tide.moreless.io/meditation_res", rootDir), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.Index(info.Name(), `%3fe%3d`) < 0 {
			return nil
		}
		hash := strings.Split(info.Name(), `%3fe%3d`)[0]
		files[hash] = path
		return nil
	})
	// 解析文件拿到hash
	buf, err := ioutil.ReadFile(fmt.Sprintf("%s/meditation.json", rootDir))
	if err != nil {
		log.Fatal(err)
	}
	var obj model.MeditationObj
	err = json.Unmarshal(buf, &obj)
	if err != nil {
		log.Fatal(err)
	}
	for _, tag := range obj.AllTags {
		_, err := dao.TagService.Save(model.TagRow{
			TagId:   tag.Id,
			SortKey: tag.SortKey,
			Key:     tag.Key,
			Type:    tag.Type,
			Name:    tag.Name["en"]})
		if err != nil {
			log.Fatal(err)
		}
	}

	parentDir := fmt.Sprintf("%s/output/med", rootDir)
	// 三级：album => section => resource
	cache := make(map[string]struct{})
	for _, album := range obj.Albums {
		for _, section := range album.Sections {
			for _, res := range section.Resources {
				if res.Languages[0] == "en" {
					src, ok := files[res.Hash]
					if ok {
						cache[album.Id] = struct{}{}
						cache[section.Id] = struct{}{}
						cache[res.Hash] = struct{}{}
						target := fmt.Sprintf("%s/res/%s.mp3", parentDir, res.Hash)
						if !fs.FileExists(target) {
							_, _ = shell.Pipe("cp", src, target)
						}
					}
				}
			}
		}
	}

	// 5e480011e38b690007961e8d 奇怪的ID
	for _, album := range obj.Albums {
		if _, ok := cache[album.Id]; !ok {
			continue
		}
		var ids []string
		for _, tag := range obj.AllTags {
			ids = append(ids, tag.Id)
		}
		err = dao.MedService.SaveMed(model.MedRow{
			Type:         dao.TypeAlbum,
			MedId:        album.Id,
			TagIds:       strings.Join(ids, ","),
			Name:         album.Name["en"],
			Description:  album.Description["en"],
			PrimaryColor: album.PrimaryColor,
			CreatedAt:    album.CreatedAt,
			UpdatedAt:    album.UpdatedAt,
			SortKey:      album.SortKey,
			CoverLink:    album.Image,
		})
		if err != nil {
			log.Fatal(err.Error())
		}
		target := fmt.Sprintf("%s/cover/%s.jpg", parentDir, album.Id)
		if len(album.Image) > 0 {
			download(album.Image, target, false)
		}

		for _, section := range album.Sections {
			if _, ok := cache[section.Id]; !ok {
				continue
			}
			err = dao.MedService.SaveSection(model.MedRow{
				Type:        dao.TypeSection,
				MedId:       album.Id,
				SectionId:   section.Id,
				Name:        section.Name["en"],
				Description: section.Description["en"],
				DemoLink:    section.DemoSoundUrlMp3["en"],
			})
			if err != nil {
				log.Fatal(err.Error())
			}
			target := fmt.Sprintf("%s/demo/%s.mp3", parentDir, section.Id)
			if len(section.DemoSoundUrlMp3["en"]) > 0 {
				download(section.DemoSoundUrlMp3["en"], target, false)
			}
			for _, res := range section.Resources {
				if _, ok := cache[res.Hash]; !ok {
					continue
				}
				err = dao.MedService.SaveRes(model.MedRow{
					Type:      dao.TypeResource,
					MedId:     album.Id,
					SectionId: section.Id,
					ResId:     res.Hash,
					Name:      res.Name,
					Duration:  res.Duration,
				})
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		}
	}
}

func parseMusicTag() {
	buf, err := ioutil.ReadFile(fmt.Sprintf("%s/music.json", rootDir))
	if err != nil {
		log.Fatal(err)
	}
	// 解析ID
	var music model.MusicItem
	err = json.Unmarshal(buf, &music)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range music.AllTags {
		id, err := dao.TagService.Save(model.TagRow{
			TagId:   item.Id,
			SortKey: item.SortKey,
			Key:     item.Key,
			Type:    item.Type,
			Name:    item.Name["en"],
		})
		fmt.Println(id, err)
	}
}

func parseMusic() {
	files := make(map[string]string)
	_ = filepath.Walk(fmt.Sprintf("%s/output/tide/resources.tide.moreless.io/sounds", rootDir), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.Index(info.Name(), `%3fe%3d`) < 0 {
			return nil
		}
		hash := strings.Split(info.Name(), `%3fe%3d`)[0]
		files[hash] = path
		return nil
	})

	buf, err := ioutil.ReadFile(fmt.Sprintf("%s/music.json", rootDir))
	if err != nil {
		log.Fatal(err)
	}
	// 解析ID
	var music model.MusicItem
	err = json.Unmarshal(buf, &music)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range music.Scenes {
		detail := fmt.Sprintf("%s/output/tide/tide-api.moreless.io/v1/market/scenes/%s", rootDir, item.Id)
		if !fs.FileExists(detail) {
			continue
		}
		var scene model.MusicScene
		buf, _ := ioutil.ReadFile(detail)
		err = json.Unmarshal(buf, &scene)
		if err != nil {
			log.Fatal(err)
		}
		sceneHash := scene.PlayList[0].Sounds[0].Hash
		if _, ok := files[sceneHash]; !ok {
			log.Fatal(sceneHash, "not.found")
		}
		dest := fmt.Sprintf("%s/output/music/res/%s.mp3", rootDir, sceneHash)
		if !fs.FileExists(dest) {
			src := files[sceneHash]
			_, _ = shell.Pipe("cp", src, dest)
		}
		var tags []string
		for _, v := range scene.TagsV2 {
			tags = append(tags, v.Id)
		}
		row := model.MusicRow{
			Id:             item.Id,
			Name:           item.Name["en"],
			DemoLink:       item.DemoSoundUrlMp3,
			CoverLink:      item.CoverUrl,
			Duration:       item.Duration,
			Description:    item.Description["en"],
			SubTitle:       item.SubTitle["en"],
			PrimaryColor:   item.PrimaryColor,
			SecondaryColor: item.SecondaryColor,
			UpdatedAt:      item.UpdatedAt,
			CreatedAt:      item.CreatedAt,
			TagIds:         strings.Join(tags, ","),
			Hash:           sceneHash,
		}
		// 下载demo
		if len(row.DemoLink) > 0 {
			parts := strings.Split(row.DemoLink, "/")
			download(row.DemoLink, fmt.Sprintf("%s/output/music/demo/%s.mp3", rootDir, parts[len(parts)-1]), false)
		}
		// 下载cover
		if len(row.CoverLink) > 0 {
			fmt.Println("cover ", row.CoverLink)
			parts := strings.Split(row.CoverLink, "/")
			download(row.CoverLink, fmt.Sprintf("%s/output/music/cover/%s.jpg", rootDir, parts[len(parts)-1]), false)
		}
		id, err := dao.MusicService.Save(row)
		fmt.Println(id, err)
	}
}

func downloadMp3() {
	fd, err := os.Open(fmt.Sprintf("%s/music.mp3.txt", rootDir))
	if err != nil {
		log.Fatal(err.Error())
	}
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "/")
		parts = strings.Split(parts[4], "?")
		download(line, fmt.Sprintf("%s/output/music/%s.mp3", rootDir, parts[0]), false)
	}
	fmt.Println("DONE")
}

func download(input, output string, override bool) {
	if !override && fs.FileExists(output) {
		return
	}
	if len(input) == 0 {
		log.Fatal("input.empty", input, output, override)
	}
	out, err := net.FetchResponse(http.MethodGet, input, bytes.NewBuffer([]byte(``)), map[string]string{
		"pragma":                    "no-cache",
		"cache-control":             "no-cache",
		"upgrade-insecure-requests": "1",
		"authority":                 "pics.tide.moreless.io",
		"user-agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36",
	}, 1)
	if err != nil {
		log.Fatal("net.FetchResponse", err.Error())
	}
	if len(out) == 0 {
		log.Fatal("out", "empty")
	}
	_ = ioutil.WriteFile(output, out, 0644)
	fmt.Println("download", output, "done")
}
