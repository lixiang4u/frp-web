package utils

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func HttpPost(requestUrl string, requestBody []byte) (int, []byte, error) {
	resp, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("[postError]", err.Error())
		return http.StatusInternalServerError, nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[responseReadError]", err.Error())
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, buf, nil
}

func HttpGet(requestUrl string) (int, []byte, error) {
	log.Println("[requestUrl]", requestUrl)
	resp, err := http.Get(requestUrl)
	if err != nil {
		log.Println("[postError]", err.Error())
		return http.StatusInternalServerError, nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[responseReadError]", err.Error())
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, buf, nil
}
