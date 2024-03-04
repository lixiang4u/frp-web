package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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

func HttpGet(requestUrl string, params ...url.Values) (int, []byte, error) {
	if len(params) != 0 {
		requestUrl = fmt.Sprintf("%s?%s", requestUrl, params[0].Encode())
	}
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

func HttpDelete(requestUrl string, requestBody []byte) (int, []byte, error) {
	req, err := http.NewRequest("DELETE", requestUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("[deleteError]", err.Error())
		return http.StatusInternalServerError, nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[responseReadError]", err.Error())
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, buf, nil
}
