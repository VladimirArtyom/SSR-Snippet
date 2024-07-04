package main

import (
  "bytes"
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


func (app *application) render(w http.ResponseWriter, pageName string, data *templateData, status int) {

  ts, exist := app.templateCache[pageName]
  if !exist {
    app.serverError(w, fmt.Errorf("Template does not exist", pageName))
  }

  var buf *bytes.Buffer = new(bytes.Buffer)

  err := ts.ExecuteTemplate(buf, "master", data)
  if err != nil {
    app.serverError(w, err)
    return
  }

  w.WriteHeader(status)
  _, err = buf.WriteTo(w)

  if err != nil {
    app.serverError(w, err)
    return
  }

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
