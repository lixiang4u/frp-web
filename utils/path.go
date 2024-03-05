package utils

import (
	"io/fs"
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

func AppTempFile(elem ...string) string {
	elem = append([]string{os.TempDir(), "frp-web"}, elem...)
	var tmpFile = filepath.Join(elem...)
	_ = os.MkdirAll(filepath.Dir(tmpFile), fs.ModePerm)
	return tmpFile
}
