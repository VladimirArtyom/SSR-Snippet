package main

import ("net/http"
  "github.com/justinas/alice"
  "github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {

  var router *httprouter.Router = httprouter.New()
  
  
  router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    app.notFound(w) 
  })

  var fileServer http.Handler = http.FileServer(http.Dir("./ui/static/"))
  router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

  router.HandlerFunc(http.MethodGet, "/", app.Home)
  router.HandlerFunc(http.MethodGet, "/snip/view/:id", app.SnipView)
  router.HandlerFunc(http.MethodGet, "/snip/create", app.SnipCreate)
  router.HandlerFunc(http.MethodPost, "/snip/create", app.SnipCreatePost)

  var standardMiddleWare alice.Chain = alice.New(app.recoverPanic, app.logRequest, secureHeaders)
  return standardMiddleWare.Then(router)
}
