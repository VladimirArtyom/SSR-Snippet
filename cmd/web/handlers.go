package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
  
    app.notFound(w)
    return
  }

  // Render page
  
  var templates []string = []string{
    "./ui/html/master.tmpl.html",
    "./ui/html/partials/nav.tmpl.html",
    "./ui/html/pages/home.tmpl.html",
  }
  ts, err := template.ParseFiles(templates...)
  if err != nil {
    app.errorLog.Println(err.Error())
    app.serverError(w, err)
    return
  }

  err = ts.ExecuteTemplate(w, "master", nil)
  if err != nil {
    app.errorLog.Println(err.Error())
    app.serverError(w, err)
    return
  }

}

func (app *application) SnipView(w http.ResponseWriter, r *http.Request) {
  
  id, err := strconv.Atoi(r.URL.Query().Get("id"))
  if id < 1 || err != nil {
    app.notFound(w)
    return
  } 

  fmt.Fprintf(w, "Display snippets from user %d", id)
  
}


func (app *application) SnipCreate(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    w.Header().Set("Allow", http.MethodPost)
    w.WriteHeader(http.StatusMethodNotAllowed)
    app.clientError(w, http.StatusMethodNotAllowed)
    return
  }

  w.Write([]byte("Create a snippet"))

}
