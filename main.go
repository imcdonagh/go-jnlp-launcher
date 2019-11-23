package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {

	// Validate and get args
	if len(os.Args) < 3 {
		fmt.Println("Usage: go-jnlp-launcher <url> <dir> [<args>...]")
		os.Exit(1)
	}
	jnlpURL := os.Args[1]
	dir := os.Args[2]

	// Download JNLP file
	jnlpFile, err := DownloadJnlp(jnlpURL)
	if err != nil {
		panic(err)
	}
	// Get JAR resources
	jars, err := jnlpFile.GetJarResources()
	if err != nil {
		panic(err)
	}
	// Download JAR files
	// TODO Progress indicator
	paths, err := DownloadJars(jnlpFile, jars, dir)
	if err != nil {
		panic(err)
	}
	// Get J2SE descriptor
	j2se := jnlpFile.GetJ2SE()
	// Get application desc
	applicationDesc := jnlpFile.GetApplicationDesc()
	// Launch appliation
	Launch("java", j2se, applicationDesc, paths, 3)
}

// DownloadJnlp downloads and parses the specified JNLP descriptor
func DownloadJnlp(url string) (*JnlpFile, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	jnlp, err := ParseJnlp(resp.Body)
	return jnlp, err
}

// DownloadJars downloads JAR resources to target directory.
func DownloadJars(jnlpFile *JnlpFile, jars []*JarResource, dir string) ([]string, error) {
	paths := make([]string, len(jars))
	cache := InitCache(dir)
	for i := 0; i < len(jars); i++ {
		jar := jars[i]
		if jar.Download == Lazy {
			continue
		}
		jarURL, err := jnlpFile.GetJarURL(jar)
		if err != nil {
			return nil, err
		}
		path, err := cache.EnsurePath(jarURL)
		if err != nil {
			return nil, err
		}
		lastmod, err := cache.GetLastModified(path)
		if err != nil {
			return nil, err
		}
		etag, err := cache.GetEtag(path)
		if err != nil {
			return nil, err
		}
		fmt.Println(path, ": ", lastmod)
		lastmod, etag, err = DownloadFile(jarURL.String(), path, lastmod, etag)
		if err != nil {
			return nil, err
		}
		err = cache.TouchFile(path, lastmod, etag)
		paths[i] = path
	}
	return paths, nil
}

// Launch launches the Java application based on the JNLP descriptor.
func Launch(java string, j2se *J2SE, applicationDesc *ApplicationDesc, jars []string, extraArgsOffset int) {
	// Build command line in the form:
	// java -Xms<initheap> -Xmx<maxheap> -cp <cp> <mainclass> <args>
	// TODO Extra args
	cp := strings.Join(jars, string(os.PathListSeparator))
	count := 3 + len(os.Args) - extraArgsOffset
	if j2se.InitialHeapSize != "" {
		count++
	}
	if j2se.MaxHeapSize != "" {
		count++
	}
	count += len(applicationDesc.Arguments)
	args := make([]string, count)
	pos := 0
	if j2se.InitialHeapSize != "" {
		args[pos] = "-Xms" + j2se.InitialHeapSize
		pos++
	}
	if j2se.MaxHeapSize != "" {
		args[pos] = "-Xmx" + j2se.MaxHeapSize
		pos++
	}
	args[pos], pos = "-cp", pos+1
	args[pos], pos = cp, pos+1
	args[pos], pos = applicationDesc.MainClass, pos+1
	for i := 0; i < len(applicationDesc.Arguments); i++ {
		args[pos], pos = applicationDesc.Arguments[i], pos+1
	}
	for i := extraArgsOffset; i < len(os.Args); i++ {
		args[pos], pos = os.Args[i], pos+1
	}
	fmt.Println("Exec: java ", strings.Join(args, " "))
	cmd := exec.Command(java, args...)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
