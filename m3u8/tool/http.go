package tool

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func Get(url string, verbose bool) (io.ReadCloser, error) {
	c := http.Client{
		Timeout: time.Minute * 5,
	}
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	if verbose {
		fmt.Println(fmt.Sprintf("[total %.4fMB] %s", float64(resp.ContentLength)/1024/1024, url))
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http error: status code %d", resp.StatusCode)
	}
	return resp.Body, nil
}
