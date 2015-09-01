package main

import (
	"github.com/russross/blackfriday"
)

type Wiki struct {
	Name     string
	Slug     string
	Sections []*Section
}

type Section struct {
	Name        string
	Slug        string
	Subsections []*Section
	Article     *Article
	Weight      int
}

type Article struct {
	Name    string
	Slug    string
	Content []byte
	Weight  int
}

func (self *Article) GetHtml() (string, error) {
	output := blackfriday.MarkdownBasic(self.Content)
	return string(output), nil
}
