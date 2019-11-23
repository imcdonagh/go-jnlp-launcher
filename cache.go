package main

import (
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

// InitCache initializes a cache backed by the specified directory.
func InitCache(dir string) *Cache {
	return &Cache{
		Dir: dir,
	}
}

// Cache represents a simple cache backed by a local directory.
type Cache struct {
	Dir string
}

// ToPath converts the specified URL to the local directory in the cache.
func (x *Cache) ToPath(url *url.URL) string {
	host := strings.Replace(url.Host, ":", "/P", 1)
	abs := path.Join(x.Dir, url.Scheme, host, url.Path)
	return abs
}

// EnsurePath ensures that the parent directory for the specified resource exists.
func (x *Cache) EnsurePath(url *url.URL) (string, error) {
	abs := x.ToPath(url)
	dir := path.Dir(abs)
	err := os.MkdirAll(dir, 0755)
	return abs, err
}

// GetLastModified specifies the last modified timestamp for a cache entry if it exists
func (x *Cache) GetLastModified(path string) (string, error) {
	file, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	lastmod := file.ModTime().UTC().Format(time.RFC1123)
	return lastmod, nil
}

// GetEtag specifies the etag for a cache entry if it exists.
func (x *Cache) GetEtag(path string) (string, error) {
	etagFile := getEtagFile(path)
	_, err := os.Stat(etagFile)
	if os.IsNotExist(err) {
		return "", nil
	}
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// TouchFile touches the specified cache file with the given last modified and/or etag.
func (x *Cache) TouchFile(file string, lastmod string, etag string) error {
	if lastmod != "" {
		mtime, err := time.Parse(time.RFC1123, lastmod)
		if err != nil {
			return err
		}
		err = os.Chtimes(file, mtime, mtime)
		if err != nil {
			return err
		}
	}
	if etag != "" {
		etagFile := getEtagFile(file)
		f, err := os.Create(etagFile)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.WriteString(etag)
		if err != nil {
			return err
		}
		f.Sync()
	}
	return nil
}

func getEtagFile(file string) string {
	return path.Join(path.Dir(file), ".etag."+path.Base(file))
}
