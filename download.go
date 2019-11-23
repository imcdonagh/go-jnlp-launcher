package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// DownloadFile downloads the specified file
// TODO Multi threaded downloads for speed
// TODO Conditional If-Modified-Since downloads
func DownloadFile(url string, file string, lastmod string, etag string) (string, string, error) {
	fmt.Println("Downloading ", url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if lastmod != "" {
		req.Header.Add("If-Modified-Since", lastmod)
	} else if etag != "" {
		req.Header.Add("If-None-Matches", etag)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if err != nil {
		return "", "", err
	}
	fmt.Println("Response: ", resp.StatusCode)
	for h := range resp.Header {
		fmt.Println(h, ": ", resp.Header.Get(h))
	}
	lastmod = resp.Header.Get("Last-Modified")
	etag = resp.Header.Get("Etag")
	if resp.StatusCode == 304 {
		err = nil
	} else {
		out, err := os.Create(file)
		defer out.Close()
		if err == nil {
			_, err = io.Copy(out, resp.Body)
		}
	}
	return lastmod, etag, err
}
