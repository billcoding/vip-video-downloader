package parse

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/billcoding/vip-video-downloader/m3u8/tool"
)

type Result struct {
	URL  *url.URL
	M3u8 *M3u8
	Keys map[int]string
}

func FromURL(link string, isFile bool) (*Result, error) {
	var (
		m3u8 *M3u8
		err  error
		u    *url.URL
	)
	if !isFile {
		u, err = url.Parse(link)
		if err != nil {
			return nil, err
		}
		link = u.String()
		body, err2 := tool.Get(link, false)
		if err2 != nil {
			return nil, fmt.Errorf("request m3u8 URL failed: %s", err2.Error())
		}
		defer body.Close()
		m3u8, err = parse(body)
		if len(m3u8.MasterPlaylist) != 0 {
			sf := m3u8.MasterPlaylist[0]
			return FromURL(tool.ResolveURL(u, sf.URI), isFile)
		}
	} else {
		file, err2 := os.Open(link)
		if err2 != nil {
			return nil, fmt.Errorf("open m3u8 File failed: %s", err2.Error())
		}
		defer file.Close()
		m3u8, err = parse(file)
	}
	if len(m3u8.Segments) == 0 {
		return nil, errors.New("can not found any TS file description")
	}
	result := &Result{
		URL:  u,
		M3u8: m3u8,
		Keys: make(map[int]string),
	}
	for idx, key := range m3u8.Keys {
		switch {
		case key.Method == "" || key.Method == CryptMethodNONE:
			continue
		case key.Method == CryptMethodAES:
			// Request URL to extract decryption key
			keyURL := key.URI
			keyURL = tool.ResolveURL(u, keyURL)
			resp, err2 := tool.Get(keyURL, false)
			if err2 != nil {
				return nil, fmt.Errorf("extract key failed: %s", err2.Error())
			}
			keyByte, err2 := ioutil.ReadAll(resp)
			_ = resp.Close()
			if err2 != nil {
				return nil, err2
			}
			fmt.Println("decryption key: ", string(keyByte))
			result.Keys[idx] = string(keyByte)
		default:
			return nil, fmt.Errorf("unknown or unsupported cryption method: %s", key.Method)
		}
	}
	return result, nil
}
