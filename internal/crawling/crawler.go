package crawling

import (
	"log"
	"net/http"
)

type CheckedLink struct {
	url       string
	status    string
	timeTaken uint
}

func checkLink(url string) (CheckedLink, error) {
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to get %s", url)
		return CheckedLink{}, err
	}
	return CheckedLink{url, response.Status, 0.0}, nil
}

func CheckLinks(urls []string) []CheckedLink {
	results := []CheckedLink{}
	for _, url := range urls {
		result, err := checkLink(url)
		if err == nil {
			results = append(results, result)
		}
	}
	return results
}
