package str

import (
	"fmt"
	"testing"
)

func TestNewBookId(t *testing.T) {
	fmt.Println(NewBookId(), len(NewBookId()))
}
