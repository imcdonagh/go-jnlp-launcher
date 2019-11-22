package main

import (
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
	return file.ModTime().UTC().Format(time.RFC1123), nil
}

// GetEtag specifies the etag for a cache entry if it exists.
func (x *Cache) GetEtag(path string) string {
	// TODO implement
	return ""
}
