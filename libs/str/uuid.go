package str

import "github.com/satori/go.uuid"

func NewGenreId() string {
	return uuid.NewV4().String()
}

func NewBookId() string {
	return uuid.NewV4().String()
}
