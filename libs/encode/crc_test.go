package encode

import (
	"fmt"
	"testing"
)

func TestCrcEncode(t *testing.T) {
	fmt.Println(CrcEncode("wiz"))
	fmt.Println(CrcEncode("inla-video-1"))
}
