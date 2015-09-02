package main

import (
	"fmt"
	"github.com/elazarl/go-bindata-assetfs"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

//go:generate go-bindata templates static
var wiki *Wiki

func handler(w http.ResponseWriter, r *http.Request) {
	pageData := &PageData{}
	// Search page
	if r.URL.Path[1:] == "search" {
		query := r.URL.Query().Get("q")
		pageData.QueryString = query
		results, err := wiki.Search(query)
		if err != nil {
			log.Fatal(err)
			return
		}
		pageData.SearchResults = make([]*Result, 0, 0)
		for _, r := range results {
			_, article := wiki.Find(r.Path)
			if article != nil {
				pageData.SearchResults = append(pageData.SearchResults, &Result{Name: article.Name, Path: r.Path})
			}
		}
		pageData.Render(w)
		return
	}
	// Index page
	if r.URL.Path[1:] == "" {
		handleIndex(pageData, nil)
		pageData.Render(w)
		return
	}
	section, article := wiki.Find(r.URL.Path[1:])
	// Section page
	if section != nil {
		handleSection(pageData, section, nil)
	}
	// Article page
	if article != nil {
		// If root article
		if article.Parent == nil {
			handleIndex(pageData, article)
		} else {
			handleSection(pageData, article.Parent, article)
		}
	}
	pageData.Render(w)
}

func handleIndex(pageData *PageData, currentArticle *Article) {
	if wiki.Sections != nil {
		if pageData.Sections == nil {
			pageData.Sections = make([]*Link, 0, 0)
		}
		for _, section := range wiki.Sections {
			pageData.Sections = append(pageData.Sections, &Link{Name: section.Name, Path: section.Slug})
		}
	}
	if wiki.Articles != nil {
		if pageData.Articles == nil {
			pageData.Articles = make([]*Link, 0, 0)
		}
		for i, article := range wiki.Articles {
			pageData.Articles = append(pageData.Articles, &Link{Name: article.Name, Path: article.Slug, Current: (currentArticle == nil && i == 0) || (currentArticle == article)})
		}
		if currentArticle == nil && len(wiki.Articles) > 0 {
			currentArticle = wiki.Articles[0]
		}
		if currentArticle != nil {
			html, err := currentArticle.GetHtml()
			if err != nil {
				log.Fatal(err)
				return
			}
			pageData.Html = html
		}
	}
}

func handleSection(pageData *PageData, section *Section, currentArticle *Article) {
	addBreadcrumbs(pageData, section)
	if section.Subsections != nil {
		if pageData.Sections == nil {
			pageData.Sections = make([]*Link, 0, 0)
		}
		for _, section := range section.Subsections {
			pageData.Sections = append(pageData.Sections, &Link{Name: section.Name, Path: section.Slug})
		}
	}
	if section.Articles != nil {
		if pageData.Articles == nil {
			pageData.Articles = make([]*Link, 0, 0)
		}
		for i, article := range section.Articles {
			pageData.Articles = append(pageData.Articles, &Link{Name: article.Name, Path: article.Slug, Current: (currentArticle == nil && i == 0) || (currentArticle == article)})
		}
		if currentArticle == nil && len(section.Articles) > 0 {
			currentArticle = section.Articles[0]
		}
		if currentArticle != nil {
			html, err := currentArticle.GetHtml()
			if err != nil {
				log.Fatal(err)
				return
			}
			pageData.Html = html
		}
	}
}

func addBreadcrumbs(pageData *PageData, section *Section) {
	if section == nil {
		return
	}
	slugs := []string{section.Slug}
	loopSection := section
	for loopSection.Parent != nil {
		loopSection = loopSection.Parent
		slugs = append(slugs, loopSection.Parent.Slug)
	}
	sort.Reverse(sort.StringSlice(slugs))
	if pageData.Breadcrumbs == nil {
		pageData.Breadcrumbs = make([]*Link, 0, 0)
	}
	pageData.Breadcrumbs = append(pageData.Breadcrumbs, &Link{Name: section.Name, Path: strings.Join(slugs, "/")})
	loopSection = section
	for i := 0; loopSection.Parent != nil; i++ {
		pageData.Breadcrumbs = append([]*Link{{Name: loopSection.Name, Path: strings.Join(slugs[:len(slugs)-i], "/")}}, pageData.Breadcrumbs...)
	}
}

func main() {
	var err error
	wiki, err = ParseWiki(".")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fs := http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: "static"})
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
