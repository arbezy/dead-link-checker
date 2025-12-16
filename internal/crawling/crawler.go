package crawling

import (
	"log"
	"net/http"
	"os"
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

func SetProxy(username string, password string) bool {
	if len(username) == 0 || len(password) == 0 {
		return false
	}

	proxyUrl := `http://` + username + `:` + password + `@PROXY:PORTNUMBER`

	// NOTE: should I defer reverting env vars after crawl has finished?
	// actually don't think I need to since setenv is only active for the current (NOT the parent) process
	// this means after exiting the program, the environment variables should be reset automatically

	os.Setenv("HTTP_PROXY", proxyUrl)
	os.Setenv("HTTPS_PROXY", proxyUrl)

	return false
}
