package tool

import (
	"io/ioutil"
	"testing"
)

func TestGet(t *testing.T) {
	body, err := Get("https://raw.githubusercontent.com/billcoding/vip-video-downloader/main/README.md", false)
	if err != nil {
		t.Error(err)
	}
	defer body.Close()
	_, err = ioutil.ReadAll(body)
	if err != nil {
		t.Error(err)
	}
}
