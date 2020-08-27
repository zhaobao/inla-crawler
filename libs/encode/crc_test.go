package encode

import (
	"fmt"
	"testing"
)

func TestCrcEncode(t *testing.T) {
	fmt.Println(CrcEncode("qgxymdmz"))
	fmt.Println(CrcEncode("cozy"))
}
