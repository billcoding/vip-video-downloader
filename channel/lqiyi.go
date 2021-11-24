package channel

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

type lQiYi struct {
	baseAPI string
}

func newLQiYi() *lQiYi {
	return &lQiYi{`http://touyongsima.lqiyi.co:5566/analysi.php?v=`}
}

func (l *lQiYi) Parse(URL string) string {
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
	return URL
}
