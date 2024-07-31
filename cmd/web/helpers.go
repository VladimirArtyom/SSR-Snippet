package main

import (
  "bytes"
  "errors"
  "fmt"
  "net/http"
  "runtime/debug"

  "github.com/go-playground/form/v4"
)


func (app *application) decodePostForm(r *http.Request, destination any) error {

  err := r.ParseForm()

  if err != nil {
    return err
  }

  err = app.formDecoder.Decode(destination, r.PostForm)
  if err != nil {

    var invalidEncode *form.InvalidDecoderError = &form.InvalidDecoderError{}
    if errors.As(err, &invalidEncode) {
      panic(err)
    }

    return err
  }

  return nil
}


func (app *application) serverError(w http.ResponseWriter,err error ) {
  var trace string = fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
  app.errorLog.Output(2, trace)

  http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
  return
}

func (app *application) render(w http.ResponseWriter, pageName string,
  data *templateData, status int) {

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

func (app *application) isAuthenticated(r *http.Request) bool{
  isAuth,ok := r.Context().Value(isAuthenticatedContextKey).(bool) 
  if !ok {
    return false
  }
  return isAuth
}

func (app *application) notFound(w http.ResponseWriter) {
  app.clientError(w, http.StatusNotFound)
  return
}
