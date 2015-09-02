package main

type PageData struct {
	QueryString   string
	Breadcrumbs   []*Link
	Sections      []*Link
	Articles      []*Link
	SearchResults []*Result
	Html          string
}

func (self *PageData) Render() string {
	return self.Html
}

type Link struct {
	Name    string
	Path    string
	Current bool
}

type Result struct {
	Name string
	Path string
	Text string
}
