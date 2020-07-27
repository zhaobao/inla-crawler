package net

import (
	"log"
	"net/http/cookiejar"
)

func NewCookieJar() *cookiejar.Jar {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: Publicsuffix.List})
	if err != nil {
		log.Fatal("new.cookie.jar", err.Error())
	}
	return jar
}
