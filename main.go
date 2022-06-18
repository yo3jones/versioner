package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

func main() {
	var (
		cwd             string
		pkg             = os.Getenv("GOPACKAGE")
		file            = os.Getenv("GOFILE")
		v               = os.Getenv("VERSION")
		syntax          []byte
		r               *regexp.Regexp
		generateComment []byte
		tmpl            *template.Template
		f               *os.File
		err             error
	)

	if cwd, err = os.Getwd(); err != nil {
		panic(err)
	}

	fmt.Println(cwd)
	fmt.Println(pkg)
	fmt.Println(file)
	fmt.Println()

	sourcePath := filepath.Join(cwd, file)
	fmt.Println(sourcePath)

	if syntax, err = ioutil.ReadFile(sourcePath); err != nil {
		panic(err)
	}

	if r, err = regexp.Compile(`(?m)^\s*//\s*go:generate\s*\S.*$`); err != nil {
		panic(err)
	}

	generateComment = r.Find(syntax)

	fmt.Println()
	fmt.Println(string(generateComment))

	versionContext := &VersionContext{
		Package:         pkg,
		GenerateComment: strings.TrimSpace(string(generateComment)),
		Version:         v,
	}

	tmpl, err = template.New("versioner").Parse(versionTemplate)
	if err != nil {
		panic(err)
	}

	if f, err = os.OpenFile(sourcePath, os.O_WRONLY, 0644); err != nil {
		panic(err)
	}
	defer f.Close()

	if err = tmpl.Execute(f, versionContext); err != nil {
		panic(err)
	}
}

type VersionContext struct {
	Package         string
	GenerateComment string
	Version         string
}

var versionTemplate = `// Code generated DO NOT EDIT.
package {{ .Package }}

{{ .GenerateComment }}

func Get() string {
	return "{{- .Version -}}"
}
`
