package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/flosch/pongo2/v5"
	"github.com/gomarkdown/markdown"
	"github.com/plus3it/gorecurcopy"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Entry struct {
	Path  string
	Title string
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func main() {
	defaultPath := ".."
	if len(os.Args) > 1 {
		defaultPath = os.Args[1]
	}
	os.RemoveAll(filepath.Join(defaultPath, "/public"))
	templatePath := filepath.Join(defaultPath, "build", "template.html")
	template := pongo2.Must(pongo2.FromFile(templatePath))
	data, err := os.ReadFile(filepath.Join(defaultPath, "content/index.json"))
	check(err)
	var entries []Entry
	json.Unmarshal([]byte(data), &entries)
	for i, s := range entries {
		fmt.Println(i, s.Path, " | ", s.Title)
		markdownFile, err := os.ReadFile(filepath.Join(defaultPath, "content", s.Path))
		check(err)
		html := markdown.ToHTML(markdownFile, nil, nil)
		dir, fileName := filepath.Split(s.Path)
		fileName = strings.TrimPrefix(fileName, "_")
		fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))
		folderPath := filepath.Join(defaultPath, "public", dir)
		fileName = filepath.Join(folderPath, fileName+".html")
		if exist, _ := exists(folderPath); !exist {
			_ = os.MkdirAll(folderPath, 0775)
		}
		fmt.Println(fileName)
		f, err := os.Create(fileName)
		check(err)
		out, err := template.Execute(pongo2.Context{"title": s.Title, "file": string(html)})
		check(err)
		_, err = f.WriteString(out)
		check(err)
	}
	err = gorecurcopy.CopyDirectory(filepath.Join(defaultPath, "src"), filepath.Join(defaultPath, "public"))
	check(err)
}
