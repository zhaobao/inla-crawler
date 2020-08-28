package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"inla/inla-crawler/libs/net"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const startLink = "http://www.popsanime.com/filter/class"

func main() {
	output, err := net.FetchResponse(http.MethodGet, startLink, nil, map[string]string{
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36",
		"Cookie":     "IPLOC=CN1100; SUV=200815124734LTU8; H5UID=1597466854412787; beegosessionID=b211294841d9b6fb52ce68009b232b70; localChannelInfo=%5B%7B%22appointUrl%22%3A%22http%3A%2F%2Fm.tv.sohu.com%2Fapp%2F%22%2C%22startapp%22%3A%221%22%2C%22channelSrc%22%3A0%2C%22cover%22%3A0%2C%22isClosed%22%3A0%2C%22timeLimit%22%3A300%2C%22channelNum%22%3A680%2C%22cid%22%3A%22%22%2C%22quality%22%3A%22nor%2Chig%2Csup%22%2C%22time%22%3A1598580491254%7D%5D",
	}, 3)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(output))
	if err != nil {
		log.Fatal(err)
	}

	matchCount := regexp.MustCompile(`^\d+`)
	doc.Find(`li.content_item`).Each(func(i int, selection *goquery.Selection) {
		itemLink, _ := selection.Find(`a`).Attr(`href`)
		uri, _ := url.Parse(startLink)
		uri.Path = itemLink
		imageLink, _ := selection.Find(`img`).Attr(`src`)
		countStr := selection.Find(`.bottom_title`).Text()
		count, err := strconv.Atoi(strings.TrimSpace(matchCount.FindString(countStr)))
		if err != nil {
			log.Fatal(err)
		}
		title := selection.Find(`.main_title`).Text()
		played := selection.Find(`.play_count`).Text()
		fmt.Println(i, uri.String(), imageLink, count, title, played)
	})
}
