package main

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ParseWiki(dirpath string) (*Wiki, error) {
	dirpath, err := filepath.Abs(dirpath)
	if err != nil {
		return nil, err
	}
	wiki := &Wiki{}
	mapping := bleve.NewIndexMapping()
	mapping.DefaultAnalyzer = "en"
	index, err := bleve.New("", mapping)
	if err != nil {
		return nil, err
	}
	wiki.Index = index
	currentPath, err := filepath.Abs(dirpath + "/../")
	if err != nil {
		return nil, err
	}
	err = filepath.Walk(dirpath, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			absPath = strings.Replace(absPath, currentPath, "", 1)
			if absPath[0] == '/' {
				absPath = strings.Replace(absPath, "/", "", 1)
			}
			parts := strings.Split(absPath, "/")
			var section *Section
			var foundSection *Section
			var article *Article
			pathParts := []string{}
			for i, part := range parts {
				name, slug, weight, icon := ParseSlug(part)
				if i == 0 {
					if wiki.Slug == "" {
						wiki.Name = name
						wiki.Slug = slug
					}
					continue
				}
				pathParts = append(pathParts, slug)
				if i == len(parts)-1 {
					article, err = ParseArticle(path, part)
					wiki.IndexArticle(strings.Join(pathParts, "/"), article)
					if err != nil {
						return err
					}
					if section == nil {
						if wiki.Articles == nil {
							wiki.Articles = make([]*Article, 0, 0)
						}
						wiki.Articles = append(wiki.Articles, article)
					} else {
						if section.Articles == nil {
							section.Articles = make([]*Article, 0, 0)
						}
						article.Parent = section
						section.Articles = append(section.Articles, article)
					}
					continue
				}
				if section == nil {
					foundSection, article = wiki.Find(slug)
					if article != nil {
						return fmt.Errorf("Article name matches section name: %v", slug)
					}
				} else {
					foundSection, article = section.Find(slug)
					if article != nil {
						return fmt.Errorf("Article name matches section name: %v", slug)
					}
				}
				if foundSection == nil {
					foundSection = &Section{Name: name, Slug: slug, Weight: weight, Icon: icon}
					if section == nil {
						if wiki.Sections == nil {
							wiki.Sections = make([]*Section, 0, 0)
						}
						wiki.Sections = append(wiki.Sections, foundSection)
					} else {
						if section.Subsections == nil {
							section.Subsections = make([]*Section, 0, 0)
						}
						section.Subsections = append(section.Subsections, foundSection)
					}
				}
				if section != nil {
					foundSection.Parent = section
				}
				section = foundSection
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return wiki, nil
}

func ParseSlug(pathname string) (name string, slug string, weight int, icon string) {
	parts := strings.Split(pathname, "__")
	if len(parts) == 2 {
		icon = parts[1]
		pathname = parts[0]
	}
	parts = strings.Split(pathname, "_")
	var err error
	weight, err = strconv.Atoi(parts[0])
	if err == nil {
		parts = parts[1:]
	} else {
		weight = 100
	}
	slug = strings.Join(parts, "_")
	name = strings.Join(parts, " ")
	name = strings.ToUpper(string(name[0])) + name[1:]
	return
}

func ParseArticle(pathname, part string) (*Article, error) {
	content, err := ioutil.ReadFile(pathname)
	if err != nil {
		return nil, err
	}
	name, slug, weight, _ := ParseSlug(part)
	article := &Article{
		Name:    name,
		Slug:    slug,
		Content: content,
		Weight:  weight,
	}
	return article, nil
}
