package net

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
)

type PTransport struct {
	http.RoundTripper
	filters map[string][]RoundTripFunc
	lock    sync.Mutex
}

type RoundTripFunc func(req *http.Request, code int, header http.Header, buf []byte) ([]byte, error)

func (p *PTransport) RegisterRoundFunc(host string, f ...RoundTripFunc) {
	if len(f) == 0 {
		return
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.filters == nil {
		p.filters = make(map[string][]RoundTripFunc)
	}
	p.filters[host] = append(p.filters[host], f...)
}

func (p *PTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := p.RoundTripper.RoundTrip(req)
	if err == nil {
		buf, _ := ioutil.ReadAll(resp.Body)
		defer func() {
			_ = resp.Body.Close()
		}()
		host := req.URL.Host
		if p.filters != nil {
			var ff []RoundTripFunc
			if items, ok := p.filters[host]; ok {
				ff = items
			}
			if len(ff) == 0 {
				if items, ok := p.filters["*"]; ok {
					ff = items
				}
			}
			for _, f := range ff {
				newBuf, err := f(req, resp.StatusCode, resp.Header, buf)
				if err != nil {
					if req.Body != nil {
						_ = req.Body.Close()
					}
					return nil, errors.New("http: nil Request.URL")
				}
				buf = newBuf
			}
		}
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	}
	return resp, err
}

func New() *PTransport {
	return &PTransport{RoundTripper: http.DefaultTransport}
}
