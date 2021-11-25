package channel

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type lQiYi struct {
	baseAPI string
}

func newLQiYi() *lQiYi {
	return &lQiYi{`http://touyongsima.lqiyi.co:5566/analysi.php?v=`}
}

func (l *lQiYi) Parse(URL string) (string, string, bool) {
	resp, err := http.Get(l.baseAPI + url.PathEscape(URL))
	if err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile(`[\s\S]+var\surls\s=\s"(.+)";[\s\S]+`)
	rawHTML := string(bytes)
	if matches := re.FindStringSubmatch(rawHTML); len(matches) >= 2 {
		URL = matches[1]
	}
	parsedURL, err := url.Parse(URL)
	if err != nil {
		panic(err)
	}
	extIndex := strings.LastIndexByte(parsedURL.Path, '.')
	if extIndex != -1 {
		URLExt := strings.ToLower(parsedURL.Path[extIndex+1:])
		switch URLExt {
		case "mp4", "mkv", "avi":
			// No need to download m3u8
			return URL, URLExt, false
		}
	}
	return URL, "", true
}
