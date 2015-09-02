package main

import (
	"fmt"
	"net/http"
	"os"
)

var wiki *Wiki

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there %s!", r.URL.Path[1:])
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
