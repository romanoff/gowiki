package main

import (
	"testing"
)

func TestArticleGetHtml(t *testing.T) {
	content := `
Article 1
=========

First *markdown* article
`
	a := &Article{Name: "First", Content: []byte(content)}
	html, err := a.GetHtml()
	if err != nil {
		t.Errorf("Unexpected error while getting article html content: %v", err)
	}
	expected := `<h1>Article 1</h1>

<p>First <em>markdown</em> article</p>
`
	if html != expected {
		t.Errorf("Expected to get:\n%v, but got:\n%v\n", expected, html)
	}
}

func TestFind(t *testing.T) {
	wiki := &Wiki{}
	wiki.Sections = []*Section{
		{Slug: "some", Subsections: []*Section{{Slug: "section"}}},
		{Slug: "new_some", Articles: []*Article{{Slug: "article"}}},
	}
	s, a := wiki.Find("some")
	if a != nil {
		t.Errorf("Expected to find section, but found article")
	}
	if s == nil {
		t.Errorf("Expected to find section, but found nil")
	}
	s, a = wiki.Find("some/section")
	if a != nil {
		t.Errorf("Expected to find section, but found article")
	}
	if s == nil {
		t.Errorf("Expected to find section, but found nil")
	}
	s, a = wiki.Find("new_some/article")
	if a == nil {
		t.Errorf("Expected to find article, but got nil")
	}
	if s != nil {
		t.Errorf("Expected to find article, but found section")
	}
}

func TestSort(t *testing.T) {
	wiki := &Wiki{}
	wiki.Sections = []*Section{
		{Slug: "some", Weight: 10, Subsections: []*Section{{Slug: "section", Weight: 9}, {Slug: "abc", Weight: 1}}},
		{Slug: "new_some", Weight: 1, Articles: []*Article{{Slug: "article"}}},
	}
	wiki.Sort()
	if wiki.Sections[0].Slug != "new_some" {
		t.Error("Expected wiki sections to be sorted, but got unsorted")
	}
	if wiki.Sections[1].Subsections[0].Slug != "abc" {
		t.Error("Expected wiki sections to be sorted, but got unsorted")
	}
}
