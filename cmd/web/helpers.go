package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter,err error ) {
  var trace string = fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
  app.errorLog.Output(2, trace)


  http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
  return
}


func (app *application) clientError(w http.ResponseWriter, status int) {
  http.Error(w, http.StatusText(status), status)
  return
}

func (app *application) notFound(w http.ResponseWriter) {
  app.clientError(w, http.StatusNotFound)
  return
}
