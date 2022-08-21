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
	Path        string
	Title       string
	IsDirectory bool
	Entries     []Entry
	RenderOnly  bool
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

var defaultPath = ".."

func getFileName(path string) (string, string, string) {
	dir, fileName := filepath.Split(path)
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))
	folderPath := filepath.Join(defaultPath, "public", dir)
	fileName = filepath.Join(folderPath, fileName+".html")
	return dir, fileName, folderPath
}

func buildFiles(entries []Entry, template *pongo2.Template, menu string) {
	for i, s := range entries {
		if s.IsDirectory {
			buildFiles(s.Entries, template, menu)
		}
		fmt.Println(i, s.Path, " | ", s.Title)
		markdownFile, err := os.ReadFile(filepath.Join(defaultPath, "content", s.Path))
		check(err)
		html := markdown.ToHTML(markdownFile, nil, nil)
		_, fileName, folderPath := getFileName(s.Path)
		if exist, _ := exists(folderPath); !exist {
			_ = os.MkdirAll(folderPath, 0775)
		}
		fmt.Println(fileName)
		f, err := os.Create(fileName)
		check(err)
		out, err := template.Execute(pongo2.Context{"title": s.Title, "file": string(html), "menu": menu})
		check(err)
		_, err = f.WriteString(out)
		check(err)
	}
}

func createMenu(entries []Entry, count int) string {
	menuStr := ""
	for _, s := range entries {
		if s.RenderOnly {
			continue
		}
		_, nameFile, _ := getFileName(s.Path)
		nameFile = strings.TrimPrefix(nameFile, "../public/")
		nameFile = "/" + nameFile
		if s.IsDirectory {
			data := createMenu(s.Entries, count+1)
			menuStr += fmt.Sprintf("<input type='checkbox' id='%s' class='toggle' />", s.Path)
			menuStr += fmt.Sprintf("<label for='%s' class='flex justify-between'>", s.Path)
			nameFile = strings.TrimSuffix(nameFile, ".html")
			nameFile = nameFile + ".html"
			menuStr += fmt.Sprintf("<a href='%s'>%s</a>", nameFile, s.Title)
			menuStr += "</label>"
			menuStr += "<ul>"
			menuStr += data
			menuStr += "</ul>"
		} else {
			if count == 0 {
				menuStr += "<li class='book-section-flat'>"
				menuStr += fmt.Sprintf("<a href='%s'>%s</a>", nameFile, s.Title)
				menuStr += "</li>"
			} else {
				menuStr += "<li>"
				menuStr += fmt.Sprintf("<a href='%s'>%s</a>", nameFile, s.Title)
				menuStr += "</li>"
			}
		}
	}
	return menuStr
}

func main() {
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
	menu := createMenu(entries, 0)
	buildFiles(entries, template, menu)
	err = gorecurcopy.CopyDirectory(filepath.Join(defaultPath, "src"), filepath.Join(defaultPath, "public"))
	check(err)
}
