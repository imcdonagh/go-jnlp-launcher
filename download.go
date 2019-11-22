package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// DownloadFile downloads the specified file
// TODO Multi threaded downloads for speed
// TODO Conditional If-Modified-Since downloads
func DownloadFile(url string, file string, lastmod string) error {
	fmt.Println("Downloading ", url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if lastmod != "" {
		req.Header.Add("If-Modified-Since", lastmod)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	fmt.Println("Response: ", resp.StatusCode)
	for h := range resp.Header {
		fmt.Println(h, ": ", resp.Header.Get(h))
	}
	if resp.StatusCode == 304 {
		err = nil
	} else {
		out, err := os.Create(file)
		if err == nil {
			_, err = io.Copy(out, resp.Body)
		}
		out.Close()
		if err == nil {
			lastmod = resp.Header.Get("Last-Modified")
			err = touchFile(file, lastmod)
		}
	}
	return err
}

func touchFile(path string, lastmod string) error {
	mtime, err := time.Parse(time.RFC1123, lastmod)
	if err != nil {
		return err
	}
	err = os.Chtimes(path, mtime, mtime)
	return err
}
