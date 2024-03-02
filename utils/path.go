package utils

import (
	"os"
	"path/filepath"
)

func AppPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return dir
}
