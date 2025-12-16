package crawling

import (
	"log"
	"net/http"
)

type CheckedLink struct {
	Url       string
	Status    string
	TimeTaken uint
}

var LinksCrawled int = 0

func checkLink(url string) (CheckedLink, error) {
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to get %s", url)
		return CheckedLink{}, err
	}
	return CheckedLink{url, response.Status, 0.0}, nil
}

// TODO: check multiple links at once
func CheckLinks(urls []string) []CheckedLink {
	results := []CheckedLink{}
	for _, url := range urls {
		result, err := checkLink(url)
		if err == nil {
			results = append(results, result)
		}
		LinksCrawled++
	}
	return results
}
