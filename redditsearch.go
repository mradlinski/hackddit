package main

import (
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

var baseRedditURL = "https://www.reddit.com/r/programming/hot"

type RedditLink struct {
	Title  string
	URL    string
	Points int
}

func getTopRedditLinks() ([]RedditLink, error) {
	resp, err := httpGet(baseRedditURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	page, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	links := make([]RedditLink, 0)
	page.Find("#siteTable .link").Each(func(idx int, s *goquery.Selection) {
		titleSel := s.Find("a.title")
		title := titleSel.Text()
		href, hrefExists := titleSel.Attr("href")
		points, err := strconv.Atoi(s.Find(".score.unvoted").Text())

		if hrefExists && err == nil {
			links = append(links, RedditLink{
				title,
				href,
				points,
			})
		}
	})

	return links, nil
}
