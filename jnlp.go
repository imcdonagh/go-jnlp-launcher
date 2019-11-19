package main

import (
	"io"
	"net/http"
	"net/url"

	"github.com/antchfx/xmlquery"
)

func DownloadJnlp(url string) (*JnlpFile, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	jnlp, err := ParseJnlp(resp.Body)
	return jnlp, err
}

type JnlpFile struct {
	Codebase string
	Base     *url.URL
	Root     *xmlquery.Node
}

type DownloadType int8

const (
	Eager    DownloadType = 0
	Progress DownloadType = 1
	Lazy     DownloadType = 2
)

type JarResource struct {
	Href     string
	Main     bool
	Download DownloadType
}

type ApplicationDesc struct {
	MainClass string
	Arguments []string
}

type Property struct {
	Name  string
	Value string
}

type J2SE struct {
	Href            string
	Version         string
	InitialHeapSize string
	MaxHeapSize     string
}

func ParseJnlp(reader io.Reader) (*JnlpFile, error) {
	doc, err := xmlquery.Parse(reader)
	if err != nil {
		return nil, err
	}
	node, err := xmlquery.Query(doc, "jnlp")
	if err != nil {
		return nil, err
	}
	codebase := node.SelectAttr("codebase")
	base, err := url.Parse(codebase)
	if err != nil {
		return nil, err
	}
	return &JnlpFile{
		Codebase: codebase,
		Base:     base,
		Root:     node,
	}, nil
}

func (x *JnlpFile) GetJarResources() ([]*JarResource, error) {
	// Get jars
	jars := xmlquery.Find(x.Root, "resources/jar")
	jarResources := make([]*JarResource, len(jars))
	for i := 0; i < len(jars); i++ {
		jar := jars[i]
		href := jar.SelectAttr("href")
		main := jar.SelectAttr("main")
		download := jar.SelectAttr("download")
		bMain := false
		if main == "true" {
			bMain = true
		}
		iDownload := Eager
		switch download {
		case "progress":
			iDownload = Progress
			break
		case "lazy":
			iDownload = Lazy
			break
		}
		jarResources[i] = &JarResource{
			Href:     href,
			Main:     bMain,
			Download: iDownload,
		}
	}
	return jarResources, nil
}

func (x *JnlpFile) GetJarURL(jar *JarResource) (*url.URL, error) {
	u, err := x.Base.Parse(jar.Href)
	return u, err
}

func (x *JnlpFile) GetProperties() []*Property {
	props := xmlquery.Find(x.Root, "resources/property")
	arr := make([]*Property, len(props))
	for i := 0; i < len(props); i++ {
		arr[i] = &Property{
			Name:  props[i].SelectAttr("name"),
			Value: props[i].SelectAttr("value"),
		}
	}
	return arr
}

func (x *JnlpFile) GetJ2SE() *J2SE {
	elem := xmlquery.FindOne(x.Root, "resources/j2se")
	href := elem.SelectAttr("href")
	version := elem.SelectAttr("version")
	initialHeapSize := elem.SelectAttr("initial-heap-size")
	maxHeapSize := elem.SelectAttr("max-heap-size")
	return &J2SE{
		Href:            href,
		Version:         version,
		InitialHeapSize: initialHeapSize,
		MaxHeapSize:     maxHeapSize,
	}
}

func (x *JnlpFile) GetApplicationDesc() *ApplicationDesc {
	desc := xmlquery.FindOne(x.Root, "application-desc")
	mainClass := desc.SelectAttr("main-class")
	args := xmlquery.Find(desc, "argument")
	argvals := make([]string, len(args))
	for i := 0; i < len(args); i++ {
		argvals[i] = args[i].InnerText()
	}
	return &ApplicationDesc{
		MainClass: mainClass,
		Arguments: argvals,
	}
}
