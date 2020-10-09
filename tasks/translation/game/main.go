package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/language"
	"inla/inla-crawler/libs/fs"
	googleTranslation "inla/inla-crawler/service/google.translation"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type entity struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Src         string `json:"src"`
	Icon        string `json:"icon"`
	Source      string `json:"source"`
	GameId      string `json:"game_id"`
	CtTime      int64  `json:"ct_time"`
	SubTitle    string `json:"sub_title"`
	ShortDesc   string `json:"short_desc"`
	Imgs        string `json:"imgs"`
	Cover       string `json:"cover"`
	CClick      int    `json:"c_click"`
	CLove       int    `json:"c_love"`
	CVisit      int    `json:"c_visit"`
	Language    string `json:"language"`
	Orientation int    `json:"orientation"`
	Quality     int    `json:"quality"`
}

const (
	workDir = "./tasks/translation/game"
	rootDir = "./"
)

func main() {
	if fs.FileExists("cache.txt") {
		return
	}

	client, err := googleTranslation.New(filepath.Join(rootDir, "assets", "inla-translate-9e4214bd7993.json"))
	if err != nil {
		log.Fatal(err)
	}
	buf, err := ioutil.ReadFile(filepath.Join(workDir, "games.en.json"))
	if err != nil {
		log.Fatal(err)
	}

	var games []entity
	err = json.Unmarshal(buf, &games)
	if err != nil {
		log.Fatal(err)
	}

	var names, categories []string
	for _, game := range games {
		names = append(names, strings.ReplaceAll(game.Name, "-", " "))
		categories = append(categories, game.Category)
	}

	if len(names) > 0 {
		outputs, err := client.TranslateText(language.Amharic, names)
		if err != nil {
			log.Fatal(err)
		}
		for i, v := range outputs {
			appendFile("cache.txt", fmt.Sprintf("%s::%s::%s", v.Source, names[i], v.Text))
			games[i].Name = v.Text
		}
	}
	if len(categories) > 0 {
		outputs, err := client.TranslateText(language.Amharic, categories)
		if err != nil {
			log.Fatal(err)
		}
		for i, v := range outputs {
			appendFile("cache.txt", fmt.Sprintf("%s::%s::%s", v.Source, categories[i], v.Text))
			games[i].Category = v.Text
		}
	}
	retBytes, _ := json.Marshal(games)
	_ = ioutil.WriteFile(filepath.Join(workDir, "game.am.json"), retBytes, 0644)
	fmt.Println(len(names))
}

func appendFile(path, message string) {
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	_, _ = fd.WriteString(message + "\n")
	_ = fd.Close()
}
