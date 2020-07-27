package gz

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func UnGzip(input []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = r.Close()
	}()
	return ioutil.ReadAll(r)
}

func Gzip(input []byte) ([]byte, error) {
	var tmp bytes.Buffer
	wr := gzip.NewWriter(&tmp)
	_, err := wr.Write(input)
	if err != nil {
		return nil, err
	}
	_ = wr.Flush()
	_ = wr.Close()
	return tmp.Bytes(), nil
}
