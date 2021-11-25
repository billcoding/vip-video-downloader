package download

import (
	"github.com/billcoding/vip-video-downloader/m3u8/tool"
	"io"
	"os"
)

type Downloader struct {
	URL    string
	Output string
}

func (d *Downloader) Start(verbose bool) {
	if reader, err := tool.Get(d.URL, verbose); err != nil {
		panic(err)
	} else {
		defer func() { _ = reader.Close() }()
		if file, err := os.OpenFile(d.Output, os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0700); err != nil {
			panic(err)
		} else {
			defer func() { _ = file.Close() }()
			buf := make([]byte, 1024*1024*100) // 100MB
			if _, err := io.CopyBuffer(file, reader, buf); err != nil {
				panic(err)
			}
		}
	}
}
