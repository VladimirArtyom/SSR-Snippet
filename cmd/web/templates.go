package main

import (
  "html/template"
  "net/http"
  "path/filepath"
  "time"

  "github.com/VladimirArtyom/SSR-snippet/internal/models"
)

/*
templateData struct represent the dynamic data
that will be needed for the template to render

Snippet *models.Snippet : pointer to Snippet

*/

type templateData struct {
  CurrentYear int
  Snippet *models.Snippet
  Snippets []*models.Snippet
}

func newTemplateData(r *http.Request) (templ *templateData) {
  templ = &templateData{
    CurrentYear: time.Now().Year(),
  }
  return templ
}

func newTemplateCache() (map[string]*template.Template, error){

  cache := map[string]*template.Template{}

  pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
  if err != nil {
    return nil, err
  }

  for _, page_relative := range pages {
    fileName := filepath.Base(page_relative)

    files := "./ui/html/master.tmpl.html"

    ts, err := template.New(fileName).
      Funcs( template.FuncMap{
	"humanDate": humanDate,},
      ).
      ParseFiles(files)
    if err != nil {
      return nil, err
    }

    ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
    if err != nil {
      return nil, err
    }

    ts, err = ts.ParseFiles(page_relative)
    if err != nil {
      return nil, err
    }

    cache[fileName] = ts
  }

  return cache, nil
}

func humanDate(t time.Time) string{
  return t.Format("02 Jan 2006 15:04")
}

