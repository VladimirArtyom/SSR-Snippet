package main

import (
  "flag"
  "log"
  "net/http"
  "os"
)


type application struct {
  errorLog *log.Logger
  infoLog *log.Logger
}

func main() {

  // Preparing but not parsing the command line
  var addr *string = flag.String("addr", "192.168.1.12:8080", "HTTP Network address")


  // Parse the command line,
  flag.Parse()

  var infoLog *log.Logger = log.New(os.Stdout, "INFO\t", log.Ldate | log.Ltime)
  var errorLog *log.Logger = log.New(os.Stderr, "ERROR\t", log.Ldate | log.Ltime | log.Lshortfile)

  var app *application = &application{
    infoLog: infoLog,
    errorLog: errorLog,
  }

  var mux *http.ServeMux = app.routes()

  infoLog.Printf("Listening on %s", *addr)
  var server *http.Server = &http.Server{
    Addr: *addr,
    Handler: mux,
    ErrorLog: errorLog,
  }

  err := server.ListenAndServe()
  errorLog.Fatal(err)
}
