package main

import (
	"io"
	"net/url"

	"github.com/antchfx/xmlquery"
)

// JnlpFile represents a JNLP descriptor
type JnlpFile struct {
	Codebase string
	Base     *url.URL
	Root     *xmlquery.Node
}

// DownloadType constants resprenting resource download attribute values
type DownloadType int8

const (
	// Eager download type
	Eager DownloadType = 0
	// Progress download type
	Progress DownloadType = 1
	// Lazy download type
	Lazy DownloadType = 2
)

// JarResource represents a jar resource in a JNLP descriptor
type JarResource struct {
	Href     string
	Main     bool
	Download DownloadType
}

// ApplicationDesc represents an application-desc in a JNLP descriptor
type ApplicationDesc struct {
	MainClass string
	Arguments []string
}

// Property represents a property resource in a JNLP descriptor
type Property struct {
	Name  string
	Value string
}

// J2SE represents a j2se resource in a JNLP descriptor
type J2SE struct {
	Href            string
	Version         string
	InitialHeapSize string
	MaxHeapSize     string
}

// ParseJnlp parses the JNLP descriptor from the given input
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

// GetJarResources specifies all the jar resources contained in the JNLP descriptor
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

// GetJarURL specifies the full URL for the given jar resource that may be used for downloading
func (x *JnlpFile) GetJarURL(jar *JarResource) (*url.URL, error) {
	u, err := x.Base.Parse(jar.Href)
	return u, err
}

// GetProperties specifies all the property resources contained in the JNLP descriptor
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

// GetJ2SE specifies the j2se resource contained in the JNLP descriptor
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

// GetApplicationDesc specifies the application-desc contained in the JNLP descriptor
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
