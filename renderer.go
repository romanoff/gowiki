package main

import (
	"fmt"
	"net/http"
	"text/template"
)

type PageData struct {
	Name          string
	QueryString   string
	Breadcrumbs   []*Link
	Sections      []*Link
	Articles      []*Link
	SearchResults []*Result
	Html          string
}

func (self *PageData) Render(w http.ResponseWriter) {
	templateName := "article"
	if self.SearchResults != nil {
		templateName = "search"
	}
	content := GetAsset("templates/" + templateName + ".tmpl")
	tmpl, err := template.New(templateName).Parse(string(content))
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	tmpl.Execute(w, self)
}

type Link struct {
	Name    string
	Icon    string
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
