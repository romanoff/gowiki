package main

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/writeas/go-strip-markdown"
	"sort"
	"strings"
)

type Wiki struct {
	Name     string
	Slug     string
	Sections []*Section
	Articles []*Article
	Index    bleve.Index
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

func (self *Wiki) IndexArticle(path string, article *Article) error {
	text := stripmd.Strip(string(article.Content))
	return self.Index.Index(path, ArticleData{Name: article.Name, Content: text})
}

func (self *Wiki) IndexSection(path string, section *Section) error {
	self.Index.Index(path, ArticleData{Name: section.Name})
	return nil
}

func (self *Wiki) Search(queryString string) ([]*SearchResult, error) {
	query := bleve.NewMatchPhraseQuery(queryString)
	search := bleve.NewSearchRequest(query)
	search.Highlight = bleve.NewHighlight()
	searchResult, err := self.Index.Search(search)
	if err != nil {
		return nil, err
	}
	results := []*SearchResult{}
	for _, result := range searchResult.Hits {
		section, article := self.Find(result.ID)
		if article == nil && section == nil {
			return nil, fmt.Errorf("%v section or article not found", result.ID)
		}
		text := ""
		for _, values := range result.Fragments {
			for _, value := range values {
				text += value
			}
		}
		name := ""
		if section != nil {
			name = section.Name
		}
		if article != nil {
			name = article.Name
		}
		results = append(results, &SearchResult{Path: result.ID, Name: name, Text: text})
	}
	return results, nil
}

type SearchResult struct {
	Path string
	Name string
	Text string
}

type Section struct {
	Name        string
	Icon        string
	Slug        string
	Parent      *Section
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

func (self *Section) GetPath() string {
	slugs := []string{self.Slug}
	section := self
	for section.Parent != nil {
		slugs = append(slugs, section.Slug)
	}
	sort.Reverse(sort.StringSlice(slugs))
	return strings.Join(slugs, "/")
}

func (self *Section) GetIcon() string {
	if self.Icon == "" {
		return "subject"
	}
	return self.Icon
}

type ArticleData struct {
	Name    string
	Content string
}

type Article struct {
	Name    string
	Slug    string
	Content []byte
	Weight  int
	Parent  *Section
}

func (self *Article) GetHtml() (string, error) {
	unsafe := blackfriday.Run(self.Content)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return string(html), nil
}

func (self *Article) GetPath() string {
	if self.Parent == nil {
		return self.Slug
	}
	return self.Parent.GetPath() + "/" + self.Slug
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
