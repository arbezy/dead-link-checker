package main

import (
	"testing"
)

func TestCheckGoogle(t *testing.T) {
	urls := []string{"https://google.com"}
	res := check_links(urls)

	if res[0].status != "200 OK" {
		t.Errorf("Unsuccessfully visited google: wanted 200, got %s", res[0].status)
	}
}
