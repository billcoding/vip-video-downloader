package channel

import (
	"fmt"
	"net/url"
	"testing"
)

func TestURLParse(t *testing.T) {
	rawURL := "https://www.google.com"
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	fmt.Println(parsedURL.Path)
	fmt.Println(parsedURL.Query())
}
