package main

import "net/http"

func (app *application) routes() *http.ServeMux {

  var mux *http.ServeMux = http.NewServeMux()

  var fileServer http.Handler = http.FileServer(http.Dir("./ui/static/"))
  mux.Handle("/static/", http.StripPrefix("/static", fileServer))
  mux.HandleFunc("/", app.Home)
  mux.HandleFunc("/snip/view", app.SnipView)
  mux.HandleFunc("/snip/create", app.SnipCreate)

  return mux
}
