package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ParseWiki(dirpath string) (*Wiki, error) {
	wiki := &Wiki{}
	err := filepath.Walk(dirpath, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			parts := strings.Split(path, "/")
			var section *Section
			var foundSection *Section
			var article *Article
			for i, part := range parts {
				if i == 0 {
					if wiki.Slug == "" {
						name, slug, _ := ParseSlug(part)
						wiki.Name = name
						wiki.Slug = slug
					}
					continue
				}
				if i == len(parts)-1 {
					article, err = ParseArticle(path, part)
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
						section.Articles = append(section.Articles, article)
					}
					continue
				}
				name, slug, weight := ParseSlug(part)
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
					foundSection = &Section{Name: name, Slug: slug, Weight: weight}
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

func ParseSlug(pathname string) (name string, slug string, weight int) {
	parts := strings.Split(pathname, "_")
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
	name, slug, weight := ParseSlug(part)
	article := &Article{
		Name:    name,
		Slug:    slug,
		Content: content,
		Weight:  weight,
	}
	return article, nil
}
