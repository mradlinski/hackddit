package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"
)

var baseSearchURL = "http://hn.algolia.com/api/v1/search"

type HNAlgoliaResponse struct {
	Hits []HNStory `json:"hits"`
}

type HNStory struct {
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Points      int       `json:"points"`
	NumComments int       `json:"num_comments"`
	CreatedAt   time.Time `json:"created_at"`
}

func normalizePath(path string) string {
	if strings.HasSuffix(path, "/") {
		return path[:len(path)-1]
	}

	return path
}

func checkIfStoryOnHN(storyURL string) (bool, error) {
	searchURL, _ := url.Parse(baseSearchURL)

	searchQuery := url.Values{}
	searchQuery.Set("query", storyURL)
	searchQuery.Set("restrictSearchableAttributes", "url")
	searchQuery.Set("tags", "story")

	now := time.Now().Unix()
	twoWeeksAgo := now - 60*60*24*7*2
	searchQuery.Set("numericFilters", fmt.Sprintf("created_at_i>%d", twoWeeksAgo))

	searchURL.RawQuery = searchQuery.Encode()

	resp, err := httpGet(searchURL.String())
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var searchData HNAlgoliaResponse
	err = json.NewDecoder(resp.Body).Decode(&searchData)
	if err != nil {
		return false, err
	}

	tenHoursAgo := now - 60*60*10

	parsedStoryURL, _ := url.Parse(storyURL)
	for _, s := range searchData.Hits {
		if s.Points < 3 && s.CreatedAt.Unix() < tenHoursAgo {
			continue
		}

		linkURL, _ := url.Parse(s.URL)
		if err != nil {
			continue
		}

		if linkURL.Host == parsedStoryURL.Host &&
			normalizePath(linkURL.Path) == normalizePath(parsedStoryURL.Path) &&
			reflect.DeepEqual(linkURL.Query(), parsedStoryURL.Query()) {
			return true, nil
		}
	}

	return false, nil
}
