package main

import (
	"fmt"
	"testing"
)

func TestParseWiki(t *testing.T) {
	wiki, err := ParseWiki("wiki_test")
	if err != nil {
		t.Errorf("Unexpected error while parsing wiki content: %v", err)
	}
	if wiki != nil {
		fmt.Println(len(wiki.Sections))
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
