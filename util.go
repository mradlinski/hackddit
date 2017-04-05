package main

import (
	"net/http"
	"time"
)

func intSliceContains(haystack []int, needle int) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}

	return false
}

var netClient = &http.Client{
	Timeout: time.Second * 10,
}

func httpGet(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")
	resp, err := netClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
