package encode

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(message string) string {
	h := md5.New()
	_, _ = h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
