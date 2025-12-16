package crawling

import (
	"os"
	"testing"
)

func TestCheckGoogle(t *testing.T) {
	urls := []string{"https://google.com"}
	res := CheckLinks(urls)

	if res[0].Status != "200 OK" {
		t.Errorf("Unsuccessfully visited google: wanted 200, got %s", res[0].Status)
	}
}

func TestProxySet(t *testing.T) {
	usn, pwd := "testusn", "testpassword"
	expected := `http://` + usn + `:` + pwd + `@exampleproxy.com:8080`

	res := SetProxy(usn, pwd)
	defer os.Unsetenv("HTTP_PROXY")
	defer os.Unsetenv("HTTPS_PROXY")

	if res != true {
		t.Errorf("no proxy details set")
	}

	httpProxy := os.Getenv("HTTP_PROXY")
	httpsProxy := os.Getenv("HTTPS_PROXY")
	if httpProxy != expected || httpsProxy != expected {
		t.Errorf("Proxy not set correctly:\nexpected %s\n", expected)
		t.Errorf("got http_proxy: %s\n", httpProxy)
		t.Errorf("got https_proxy: %s\n", httpsProxy)
	}
}

func TestNoProxy(t *testing.T) {
	usn, pwd := "", ""
	expected := ""

	res := SetProxy(usn, pwd)

	if res != false {
		t.Errorf("Proxy set despite empty username and password")
	}

	httpProxy := os.Getenv("HTTP_PROXY")
	httpsProxy := os.Getenv("HTTPS_PROXY")
	if httpProxy != expected || httpsProxy != expected {
		t.Errorf("Proxy not set correctly:\nexpected %s\n", expected)
		t.Errorf("got http_proxy: %s\n", httpProxy)
		t.Errorf("got https_proxy: %s\n", httpsProxy)
	}
}
