package str

import (
	"fmt"
	"testing"
)

func TestRandStr(t *testing.T) {
	for i := 0; i < 1000; i++ {
		fmt.Println(RandStr(40))
	}
}
