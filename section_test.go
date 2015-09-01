package main

import (
	"testing"
	"fmt"
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
	fmt.Println(html)
}
