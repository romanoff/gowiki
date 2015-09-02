package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var wiki *Wiki

func handler(w http.ResponseWriter, r *http.Request) {
	pageData := &PageData{}
	if r.URL.Path[1:] == "" {
		return
	}
	_, article := wiki.Find(r.URL.Path[1:])
	if article != nil {
		html, err := article.GetHtml()
		if err != nil {
			log.Fatal(err)
		}
		pageData.Html = html
	}
	fmt.Fprintf(w, pageData.Render())
}

func main() {
	var err error
	wiki, err = ParseWiki(".")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
