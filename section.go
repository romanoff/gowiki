package main

import (
	"github.com/russross/blackfriday"
	"sort"
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

func (self *Wiki) Sort() {
	if self.Sections != nil {
		sort.Sort(BySectionWeight(self.Sections))
		for _, section := range self.Sections {
			section.Sort()
		}
	}
	if self.Articles != nil {
		sort.Sort(ByArticleWeight(self.Articles))
	}
}

type Section struct {
	Name        string
	Slug        string
	Subsections []*Section
	Articles    []*Article
	Weight      int
}

func (self *Section) Sort() {
	if self.Subsections != nil {
		sort.Sort(BySectionWeight(self.Subsections))
		for _, subsection := range self.Subsections {
			subsection.Sort()
		}
	}
	if self.Articles != nil {
		sort.Sort(ByArticleWeight(self.Articles))
	}
}

func (self *Section) Find(path string) (*Section, *Article) {
	parts := strings.Split(path, "/")
	if len(parts) == 1 && self.Articles != nil {
		for _, article := range self.Articles {
			if article.Slug == parts[0] {
				return nil, article
			}
		}
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

type BySectionWeight []*Section

func (self BySectionWeight) Len() int {
	return len(self)
}
func (self BySectionWeight) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}
func (self BySectionWeight) Less(i, j int) bool {
	if self[i].Weight == self[j].Weight {
		return self[i].Slug < self[j].Slug
	}
	return self[i].Weight < self[j].Weight
}

type ByArticleWeight []*Article

func (self ByArticleWeight) Len() int {
	return len(self)
}
func (self ByArticleWeight) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}
func (self ByArticleWeight) Less(i, j int) bool {
	if self[i].Weight == self[j].Weight {
		return self[i].Slug < self[j].Slug
	}
	return self[i].Weight < self[j].Weight
}
