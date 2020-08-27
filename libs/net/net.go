package net

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const maxRetry = 3
const defaultTimeout = 0

var defaultClient = &http.Client{
	Timeout: defaultTimeout,
}

func FetchResponse(method, url string, body io.Reader, headers map[string]string, retry int) ([]byte, error) {
	//fmt.Println("HTTP.METHOD", method)
	fmt.Println("FETCH.HTTP.URL", url)
	//fmt.Println("HTTP.BODY", body)
	//fmt.Println("HTTP.HEADERS", headers)
	//fmt.Println("HTTP.RETRY", retry)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.New("http.new.request " + err.Error())
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := defaultClient.Do(req)
	if err != nil {
		return nil, errors.New("http.do.request." + err.Error())
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if retry >= maxRetry {
			return nil, errors.New("http.read.response." + err.Error())
		}
		fmt.Println("http.read.response.error.do.retry", err.Error(), retry)
		retry++
		time.Sleep(time.Second * 5)
		return FetchResponse(method, url, body, headers, retry)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return buf, nil
}
