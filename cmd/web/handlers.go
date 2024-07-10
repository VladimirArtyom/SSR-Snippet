package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
  "github.com/VladimirArtyom/SSR-snippet/internal/models"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {

  snippets, err := app.snippets.Latest()
  if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
      app.notFound(w)
    } else {
      app.serverError(w, err)
    }
  }
  var templateData *templateData = newTemplateData(r)
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

  var templateData *templateData = newTemplateData(r)
  templateData.Snippet = snip

  app.render(w, "view.tmpl.html", templateData, http.StatusOK)

}

func (app *application) SnipCreatePost(w http.ResponseWriter, r *http.Request) {

  var title string = "0 snail"
  var content string = "0 snail\n Climb Mount Fuji, \nBut slowly, slowly!\n\n- Kobayashi Issa"
  var expires int = 7

  id, err := app.snippets.Insert(title, content,  expires)
  if err != nil {
    app.serverError(w, err)
    return
  }

  http.Redirect(w, r, fmt.Sprintf("/snip/view/%d", id), http.StatusSeeOther)

}

func (app *application) SnipCreate(w http.ResponseWriter, r *http.Request) {

}
