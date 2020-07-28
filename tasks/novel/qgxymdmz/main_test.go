package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	input := `A high mountain, a low-lying crag, a fresh spring, an ancient pine, a brightly burning oven, a pot of green tea, an old man, a young boy.

“What is the most fearsome weapon under heaven?” asked the young boy. “Is it Little Li's flying dagger, which never misses its target?”

“It used to be, but not anymore.”`

	content := strings.ReplaceAll(input, `“`, `"`)
	content = strings.ReplaceAll(content, `”`, `"`)
	content = strings.ReplaceAll(content, `‘`, `'`)
	content = strings.ReplaceAll(content, `’`, `'`)
	fmt.Println(content)

}
