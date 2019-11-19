package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Validate and get args
	if len(os.Args) < 3 {
		fmt.Println("Usage: get <url> <dir>")
		os.Exit(1)
		return
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
	// Launch
	Launch("java", j2se, applicationDesc, paths)
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
		err = DownloadFile(jarURL.String(), path)
		if err != nil {
			return nil, err
		}
		paths[i] = path
	}
	return paths, nil
}

// Launch launches the Java application based on the JNLP descriptor.
func Launch(java string, j2se *J2SE, applicationDesc *ApplicationDesc, jars []string) {
	// Build command line in the form:
	// java -Xms<initheap> -Xmx<maxheap> -cp <cp> <mainclass> <args>
	// TODO Extra args
	cp := strings.Join(jars, string(os.PathListSeparator))
	count := 3
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
	args[pos] = "-cp"
	pos++
	args[pos] = cp
	pos++
	args[pos] = applicationDesc.MainClass
	pos++
	for i := 0; i < len(applicationDesc.Arguments); i++ {
		args[pos] = applicationDesc.Arguments[i]
		pos++
	}
	fmt.Println("Exec: java ", strings.Join(args, " "))
	cmd := exec.Command(java, args...)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
