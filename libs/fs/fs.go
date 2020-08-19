package fs

import "os"

func FileExists(input string) bool {
	_, err := os.Stat(input)
	return err == nil
}
