package main

import (
	//"log"
)

func GetUrls() ([]string, error) {
	urls := []string{}
	urls = append(urls, `https://google.com`)
	return urls, nil
}
