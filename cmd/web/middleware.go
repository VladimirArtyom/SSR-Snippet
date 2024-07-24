package main

import (
	"fmt"
	"net/http"
)

const (
  cspHeader = "Content-Security-Policy"
  referrerHeader = "Referrer-Policy"
  xContentTypeHeader = "X-Content-Type-Options"
  xFrameOptionsHeader = "X-Frame-Options"
  xXssProtectionHeader = "X-XSS-Protection"
)

func secureHeaders(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Code will be executed on the way down
    w.Header().Set(cspHeader, "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
    w.Header().Set(referrerHeader, "origin-when-cross-origin")
    w.Header().Set(xContentTypeHeader, "no-sniff")
    w.Header().Set(xFrameOptionsHeader, "deny")
    w.Header().Set(xXssProtectionHeader, "0")

    next.ServeHTTP(w, r)

    // Code below will be executed on the way up including deffered funcs
  });
}


// Use InformationLogger to record IP Request, Protocol, URL, et method 
func (app *application) logRequest(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r* http.Request){
    app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto,
      r.URL.RequestURI(), r.Method)
      next.ServeHTTP(w, r)
  })
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r* http.Request) {
    // Execute when stack is up i.e when respnse will be given
    defer func() {
      err := recover();
      if err != nil {
        w.Header().Set("Connection", "close")
        app.serverError(w, fmt.Errorf("%s", err))
      }
    }()
      // Execute when the request is achieve
      next.ServeHTTP(w, r)
  })
}


func (app *application) requireAuth(next http.Handler) http.Handler {

  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if !app.isAuthenticated(r) {
      http.Redirect(w, r, "/user/login", http.StatusSeeOther);
      return
    }

    w.Header().Add("Cache-Control", "no-store")
    next.ServeHTTP(w, r)
  }) 
  
}
