package main

import (
	"testing"
)

func TestParseWiki(t *testing.T) {
	wiki, err := ParseWiki("wiki_test")
	if err != nil {
		t.Errorf("Unexpected error while parsing wiki content: %v", err)
	}
	if len(wiki.Sections) != 2 {
		t.Errorf("Expected wiki to have 2 sections, but got %v", len(wiki.Sections))
	}
	if wiki.Sections[0].Name != "Section 1" {
		t.Errorf("Expected section name to be 'Section 1', but got '%v'", wiki.Sections[0].Name)
	}
	if wiki.Sections[1].Name != "Section 2" {
		t.Errorf("Expected section name to be 'Section 2', but got '%v'", wiki.Sections[1].Name)
	}
	_, article := wiki.Find("section_1/about")
	if article == nil {
		t.Error("Expected to find about article, but got nil")
	}
	if len(article.Content) == 0 {
		t.Error("Expected to get about article content, but was empty")
	}
}

func TestWikiSearch(t *testing.T) {
	wiki, err := ParseWiki("wiki_test")
	if err != nil {
		t.Errorf("Unexpected error while parsing wiki content: %v", err)
	}
	results, err := wiki.Search("text")
	if err != nil {
		t.Errorf("Not expected to get error while searching thorough wiki, but got %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected to get 1 search result while searching on wiki, but got %v", len(results))
	}
}

func TestParseSlug(t *testing.T) {
	name, slug, weight := ParseSlug("12_some_article")
	if name != "Some article" {
		t.Errorf("Parsed slug name is expected to be 'Some article', but got '%v'", name)
	}
	if slug != "some_article" {
		t.Errorf("Parsed slug slug is expected to be 'some_article', but got '%v'", slug)
	}
	if weight != 12 {
		t.Errorf("Parsed slug weight is expected to be '12', but got '%v'", weight)
	}
}
