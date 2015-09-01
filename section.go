package main

import (
	"github.com/russross/blackfriday"
)

type Section struct {
	Name string
	Subsections []*Section
	Article *Article
}

type Article struct {
	Name string
	Content []byte
}

func (self *Article) GetHtml() (string, error) {
	output := blackfriday.MarkdownBasic(self.Content)
	return string(output), nil
}
