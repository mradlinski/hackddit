package main

import (
	"net/url"
	"strings"
)

var baseHNSubmitURL = "http://news.ycombinator.com/submitlink"

func createSubmitLink(scrapeURL string, title string) (string, error) {
	submitURL, _ := url.Parse(baseHNSubmitURL)

	submitQuery := url.Values{}
	submitQuery.Set("u", scrapeURL)
	submitQuery.Set("t", title)

	submitURL.RawQuery = strings.Replace(submitQuery.Encode(), "+", "%20", -1)

	return submitURL.String(), nil
}
