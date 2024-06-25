package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/VladimirArtyom/SSR-snippet/internal/models"
	_ "github.com/go-sql-driver/mysql"
)


type application struct {
  errorLog *log.Logger
  infoLog *log.Logger
  snippets *models.SnippetModel
}

func openDB(dsn string) (*sql.DB, error) {

  db, err := sql.Open("mysql", dsn)
  if err != nil {
    return nil, err
  }

  if err = db.Ping(); err != nil {
    return nil, err
  }

  return db, nil
}

func main() {

  // Preparing but not parsing the command line
  var addr *string = flag.String("addr", "192.168.1.12:8080", "HTTP Network address")
  var dsn *string = flag.String("dsn", "web:123456@tcp(localhost:3306)/snippetbox?parseTime=true", "MYSQL data source name")



  // Parse the command line,
  flag.Parse()

  var infoLog *log.Logger = log.New(os.Stdout, "INFO\t", log.Ldate | log.Ltime)
  var errorLog *log.Logger = log.New(os.Stderr, "ERROR\t", log.Ldate | log.Ltime | log.Lshortfile)
  fmt.Println(*dsn)
  db, err := openDB(*dsn)
  if  err != nil {
    errorLog.Fatal(err)
  }

  defer db.Close()

  var app *application = &application{
    infoLog: infoLog,
    errorLog: errorLog,
    snippets: &models.SnippetModel{DB: db},
  }

  var mux *http.ServeMux = app.routes()

  infoLog.Printf("Listening on %s", *addr)
  var server *http.Server = &http.Server{
    Addr: *addr,
    Handler: mux,
    ErrorLog: errorLog,
  }

  err = server.ListenAndServe()
  errorLog.Fatal(err)
}
