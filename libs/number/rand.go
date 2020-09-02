package number

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func RandInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func RandInt64(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}
