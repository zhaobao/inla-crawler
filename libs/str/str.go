package str

import (
	"math/rand"
	"strings"
	"time"
)

const sourceChar = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const sourceSize = 62

func init() {
	rand.Seed(time.Now().Unix())
}

func RandStr(size int) string {
	output := make([]byte, size, size)
	for i := 0; i < size; i++ {
		output[i] = sourceChar[rand.Intn(sourceSize)]
	}
	return string(output)
}

func CleanText(content string) string {
	content = strings.ReplaceAll(content, `“`, `"`)
	content = strings.ReplaceAll(content, `”`, `"`)
	content = strings.ReplaceAll(content, `‘`, `'`)
	content = strings.ReplaceAll(content, `’`, `'`)
	content = strings.ReplaceAll(content, `—`, `-`)
	content = strings.ReplaceAll(content, `…`, `...`)
	return content
}
