package redoc

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"

	_ "embed"
)

// ErrSpecNotFound error for when spec file not found.
var ErrSpecNotFound = errors.New("spec not found")

// Redoc configuration.
type Redoc struct {
	// <title> HTML tag.
	Title string `json:"title" yaml:"title"`
	// Meta description.
	Description string `json:"description" yaml:"description"`
	// Represents the URL path in case you use the Redoc.Handler function.
	DocsPath string `json:"docsPath" yaml:"docsPath"`
	// Represents the URL where the openapi.yaml/json will be exposed(as static file
	// in case the relative path doesn't works).
	SpecPath    string `json:"specPath" yaml:"specPath"`
	SpecFile    string `json:"specFile" yaml:"specFile"`
	FaviconPath string `json:"faviconPath" yaml:"faviconPath"`
}

// HTML represents the redoc index.html page
//go:embed assets/index.html
var HTML string

// JavaScript represents the redoc standalone javascript
//go:embed assets/redoc.standalone.js
var JavaScript string

// Body returns the final html with the js in the body.
func (r *Redoc) Body() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	tpl, err := template.New("redoc").Parse(HTML)
	if err != nil {
		return nil, err
	}

	var faviconType string

	if strings.HasSuffix(r.FaviconPath, ".png") {
		faviconType = "image/png"
	} else {
		faviconType = "image/x-icon"
	}

	err = tpl.Execute(buf, map[string]string{
		"body":        JavaScript,
		"title":       r.Title,
		"url":         r.SpecPath,
		"description": r.Description,
		"faviconPath": r.FaviconPath,
		"faviconType": faviconType,
	})

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Handler sets some defaults and returns a HandlerFunc.
func (r *Redoc) Handler() http.HandlerFunc {
	data, err := r.Body()
	if err != nil {
		panic(err)
	}

	specFile := r.SpecFile
	if specFile == "" {
		panic(ErrSpecNotFound)
	}

	spec, err := ioutil.ReadFile(specFile)
	if err != nil {
		panic(err)
	}

	docsPath := r.DocsPath

	return func(w http.ResponseWriter, req *http.Request) {
		method := strings.ToUpper(req.Method)
		if method != http.MethodGet && method != http.MethodHead {
			return
		}

		if strings.HasSuffix(req.URL.Path, r.SpecPath) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write(spec)

			return
		}

		if docsPath == "" || docsPath == req.URL.Path {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "text/html")
			w.Write(data)
		}
	}
}
