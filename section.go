package main

import (
	"github.com/russross/blackfriday"
	"strings"
)

type Wiki struct {
	Name     string
	Slug     string
	Sections []*Section
	Articles []*Article
}

func (self *Wiki) Find(path string) (*Section, *Article) {
	parts := strings.Split(path, "/")
	if self.Sections == nil && self.Articles == nil {
		return nil, nil
	}
	if self.Articles != nil && len(parts) == 1 {
		for _, article := range self.Articles {
			if article.Slug == parts[0] {
				return nil, article
			}
		}
	}
	if self.Sections == nil {
		return nil, nil
	}
	for _, section := range self.Sections {
		if section.Slug == parts[0] {
			if len(parts) == 1 {
				return section, nil
			}
			return section.Find(strings.Join(parts[1:], "/"))
		}
	}
	return nil, nil
}

type Section struct {
	Name        string
	Slug        string
	Subsections []*Section
	Article     *Article
	Weight      int
}

func (self *Section) Find(path string) (*Section, *Article) {
	parts := strings.Split(path, "/")
	if len(parts) == 1 && self.Article != nil && self.Article.Slug == parts[0] {
		return nil, self.Article
	}
	if self.Subsections == nil {
		return nil, nil
	}
	for _, section := range self.Subsections {
		if section.Slug == parts[0] {
			if len(parts) == 1 {
				return section, nil
			}
			return section.Find(strings.Join(parts[1:], "/"))
		}
	}
	return nil, nil
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
