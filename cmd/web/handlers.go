package main

import (
  "database/sql"
  "errors"
  "fmt"
  "net/http"
  "strconv"
  "strings"
  "unicode/utf8"
  "github.com/VladimirArtyom/SSR-snippet/internal/models"
  "github.com/julienschmidt/httprouter"
)
type snippetCreateFrom struct {
  Title string
  Content string
  Expire int
  FieldError map[string]string
}

func (app *application) Home(w http.ResponseWriter, r *http.Request) {

  snippets, err := app.snippets.Latest()
  if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
      app.notFound(w)
    } else {
      app.serverError(w, err) }
  }
  var templateData *templateData = app.newTemplateData(r)
  templateData.Snippets = snippets
  // Render page

  app.render(w, "home.tmpl.html", templateData, http.StatusOK)
}

func (app *application) SnipView(w http.ResponseWriter, r *http.Request) {


  var params httprouter.Params = httprouter.ParamsFromContext(r.Context())


  id, err := strconv.Atoi(params.ByName("id"))
  if id < 1 || err != nil {
    app.notFound(w)
    return
  } 


  snip, err :=  app.snippets.Get(id)
  if err != nil {
    if errors.Is(err, models.ErrNoRecord) {
      app.notFound(w)
    } else {
      app.serverError(w, err)
    }
    return
  }

  var templateData *templateData = app.newTemplateData(r)
  templateData.Snippet = snip

  app.render(w, "view.tmpl.html", templateData, http.StatusOK)

}

func (app *application) SnipCreatePost(w http.ResponseWriter, r *http.Request) {


  err := r.ParseForm()
  if err != nil {
    app.clientError(w, http.StatusBadRequest)
    return 
  }

  var title string = r.PostForm.Get("title")
  var content string = r.PostForm.Get("content")
  expire, err := strconv.Atoi(r.PostForm.Get("expires"))
  if err != nil {
    app.serverError(w, err)
    return
  }

  var fieldError map[string]string = make(map[string]string)

  if strings.TrimSpace(title) == "" {
    fieldError["title"] = "This field cannot be blank"
  } else if utf8.RuneCountInString(title) > 100 {
    fieldError["title"] = "This field cannot be more than 100 characters long"
  }

  if strings.TrimSpace(content) == "" {
    fieldError["content"] = "This field cannot be blank" 
  }

  if expire != 7 && expire != 30 && expire != 360 && expire != 1 {
    fieldError["expire"] = "This field must equal to 1, 7, 30, ou 365"
  }


  if len(fieldError) > 0 {
    var snipCreateForm *snippetCreateFrom = &snippetCreateFrom{
      Title: title,
      Content: content,
      Expire: expire,
      FieldError: fieldError,
    }
    // Render the page
    var data *templateData = app.newTemplateData(r)
    data.Form = snipCreateForm
    app.render(w, "create.tmpl.html", data, http.StatusUnprocessableEntity)
    return
  }

  id, err := app.snippets.Insert(title, content,  expire)
  if err != nil {
    app.serverError(w, err)
    return
  }

  http.Redirect(w, r, fmt.Sprintf("/snip/view/%d", id), http.StatusSeeOther)

}

func (app *application) SnipCreate(w http.ResponseWriter, r *http.Request) {
  // Render html create.tmpl
  data := app.newTemplateData(r)
  data.Form = &snippetCreateFrom{
    Expire: 365,
  }
  app.render(w, "create.tmpl.html",  data, http.StatusOK)
}
