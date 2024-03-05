package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
)

func FileExists(files ...string) bool {
	if len(files) == 0 {
		return false
	}
	for _, file := range files {
		_, err := os.Stat(file)
		if err != nil {
			return false
		}
	}
	return true
}

func HashFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	sum := hash.Sum(nil)
	return hex.EncodeToString(sum), nil
}

func FileContents(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		return make([]byte, 0)
	}
	defer func() { _ = file.Close() }()

	buf, err := io.ReadAll(file)
	if err != nil {
		return make([]byte, 0)
	}
	return buf
}

func FileWriteContent(filename string, data []byte) error {
	return os.WriteFile(filename, data, fs.ModePerm)
}
