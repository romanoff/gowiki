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
