package google_translation

import (
	"fmt"
	"golang.org/x/text/language"
	"testing"
)

func TestNew(t *testing.T) {
	client, err := New("./inla-translate-9e4214bd7993.json")
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.TranslateText(language.Amharic, []string{"Hello, my friends! Good day"})
	if err != nil {
		t.Fatal(err)
	}

	for _, line := range resp {
		fmt.Println(line.Text)
	}
}
