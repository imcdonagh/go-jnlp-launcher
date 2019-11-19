package main

import (
	"io"
	"net/http"
	"os"
)

// DownloadFile downloads the specified file
// TODO Multi threaded downloads for speed
// TODO Conditional If-Modified-Since downloads
func DownloadFile(url string, file string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(file)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}
