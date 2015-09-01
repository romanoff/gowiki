package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ParseWiki(dirpath string) (*Wiki, error) {
	wiki := &Wiki{}
	filepath.Walk(dirpath, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			parts := strings.Split(path, "/")
			// var section *Section
			for i, part := range parts {
				if i == 0 && wiki.Slug == "" {
					name, slug, _ := ParseSlug(part)
					wiki.Name = name
					wiki.Slug = slug
					continue
				}
				if i == len(parts)-1 {
					continue
				}
			}
		}
		return nil
	})
	return wiki, nil
}

func ParseSlug(pathname string) (name string, slug string, weight int) {
	parts := strings.Split(pathname, "_")
	var err error
	weight, err = strconv.Atoi(parts[0])
	if err == nil {
		parts = parts[1:]
	}
	slug = strings.Join(parts, "_")
	name = strings.Join(parts, " ")
	name = strings.ToUpper(string(name[0])) + name[1:]
	return
}
