package crawling

import (
// "log"
)

// TODO: read urls from a txt file

func GetUrls() ([]string, error) {
	urls := []string{}
	urls = append(urls, `https://google.com`)
	return urls, nil
}
