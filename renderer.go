package main

import (
	"fmt"
	"net/http"
	"text/template"
)

type PageData struct {
	QueryString   string
	Breadcrumbs   []*Link
	Sections      []*Link
	Articles      []*Link
	SearchResults []*Result
	Html          string
}

func (self *PageData) Render(w http.ResponseWriter) {
	content := GetAsset("templates/article.tmpl")
	tmpl, err := template.New("article").Parse(string(content))
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	tmpl.Execute(w, self)
}

type Link struct {
	Name    string
	Path    string
	Current bool
}

type Result struct {
	Name string
	Path string
	Text string
}

func GetAsset(path string) []byte {
	data, err := Asset(path)
	if err != nil {
		panic(err)
	}
	return data
}
