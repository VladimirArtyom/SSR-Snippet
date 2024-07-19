package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VladimirArtyom/SSR-snippet/internal/models"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

// The application struct contains all the things that web app needs access to.
// Mainly it is pour communicating with another modules.
// and a global access pour the app.
// This is not a good practice imo.
type application struct {
  errorLog *log.Logger
  infoLog *log.Logger
  templateCache map[string]*template.Template
  formDecoder *form.Decoder
  sessionManager *scs.SessionManager

  snippets *models.SnippetModel
  users *models.UserModel
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

  templateCache, err := newTemplateCache()
  if err != nil {
    errorLog.Fatal(err)
  } 
  infoLog.Printf("Cache succesfully saved by %s", "Primitif")

  // init form decoder
  var formDecoder *form.Decoder = form.NewDecoder()

  // init sessionManager
  var sessionManager *scs.SessionManager = scs.New()
  sessionManager.Store = mysqlstore.New(db)
  sessionManager.Lifetime = 12 * time.Hour // Life time is 12 hours, start when the session is created

  var app *application = &application{
    infoLog: infoLog,
    errorLog: errorLog,
    templateCache: templateCache,
    formDecoder: formDecoder,
    sessionManager: sessionManager,

    snippets: &models.SnippetModel{DB: db},
    users: &models.UserModel{DB: db},
  }

  var tlsConfig *tls.Config = &tls.Config{
    CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
    MinVersion: tls.VersionTLS12,
    MaxVersion: tls.VersionTLS13,
    CipherSuites: []uint16{
      tls.TLS_AES_128_GCM_SHA256,
      tls.TLS_CHACHA20_POLY1305_SHA256,
      tls.TLS_AES_256_GCM_SHA384,
      tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
      tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
      tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
    },
  }

  var mux http.Handler = app.routes()

  infoLog.Printf("Listening on %s", *addr)
  var server *http.Server = &http.Server{
    Addr: *addr,
    Handler: mux,
    ErrorLog: errorLog,
    TLSConfig: tlsConfig,
    IdleTimeout: 2 * time.Minute,
    ReadTimeout: 5 * time.Second,
    WriteTimeout: 10 * time.Second,
  }

  err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
  errorLog.Fatal(err)
}
