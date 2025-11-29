package main

import (
	"log"
	"net/http"
)

type checkedLink struct {
	url       string
	status    string
	timeTaken uint
}

func check_link(url string) (checkedLink, error) {
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to get %s", url)
		return checkedLink{}, err
	}
	return checkedLink{url, response.Status, 0.0}, nil
}

func check_links(urls []string) []checkedLink {
	results := []checkedLink{}
	for _, url := range urls {
		result, err := check_link(url)
		if err == nil {
			results = append(results, result)
		}
	}
	return results
}

// func get_urls() []string {
// 	urls := []string{}
// 	urls = append(urls, "https://google.com")
// 	return urls
// }
