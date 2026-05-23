package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	userAgent   = "opendart-doc-crawler (github.com/kenshin579/opendart)"
	politeDelay = 300 * time.Millisecond
)

// httpGet 은 User-Agent 를 붙여 URL 본문을 문자열로 가져온다.
func httpGet(client *http.Client, url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GET %s: status %d", url, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
