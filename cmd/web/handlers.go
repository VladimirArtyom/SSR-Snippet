package main

import (
  "database/sql"
  "errors"
  "fmt"
  "net/http"
  "strconv"
  "unicode/utf8"
  "github.com/VladimirArtyom/SSR-snippet/internal/models"
  "github.com/VladimirArtyom/SSR-snippet/internal/validator"
  "github.com/julienschmidt/httprouter"
)

type snippetCreateFrom struct {
  Title string `form:"title"`
  Content string `form:"content"`
  Expire int `form:"expire"`
  validator.Validator `form:"-"`// So we can have access into the validator package
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

  var form *snippetCreateFrom = &snippetCreateFrom{}
  err := app.decodePostForm(r, form)

  if err != nil {
    app.clientError(w, http.StatusBadRequest)
    return
  }

  form.CheckField(validator.IsNotBlank(form.Title), "title" , validator.BLANK_MESSAGE)
  form.CheckField(validator.IsNotMaxChars(form.Title, 100), "title", fmt.Sprintf(validator.MAX_CHAR_MESSAGE,
    utf8.RuneCountInString(form.Title)))
  form.CheckField(validator.IsNotBlank(form.Content), "content", validator.BLANK_MESSAGE)
  form.CheckField(validator.PermittedInt(form.Expire, 1, 7, 30, 360), "expire", validator.NOT_IN_OPTIONS)


  if !form.IsValid() {
    var data *templateData = app.newTemplateData(r)
    data.Form = form

    app.render(w, "create.tmpl.html", data, http.StatusUnprocessableEntity)
    return
  }

  id, err := app.snippets.Insert(form.Title,
    form.Content,
    form.Expire)
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
