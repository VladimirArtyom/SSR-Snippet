package main

import (
  "database/sql"
  "errors"
  "fmt"
  "net/http"
  "strconv"
  "unicode/utf8"
  "github.com/VladimirArtyom/SSR-snippet/internal/models"
  "github.com/VladimirArtyom/SSR-snippet/internal/validator"
  "github.com/julienschmidt/httprouter"
)

type snippetCreateFrom struct {
  Title string `form:"title"`
  Content string `form:"content"`
  Expire int `form:"expire"`
  validator.Validator `form:"-"`// So we can have access into the validator package
}

type userLoginForm struct {
  Email string `form:"email"`
  Password string `form:"password"`
  validator.Validator `form:"-"`
}

type userSingupForm struct {
  Name string `form:"name"`
  Email string `form:"email"`
  Password string `form:"password"` 
  validator.Validator `form:"-"`
}

func (app *application) Home(w http.ResponseWriter, r *http.Request) {

  snippets, err := app.snippets.Latest()
  if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
      app.notFound(w)
    } else {
      app.serverError(w, err) }
  }
  var templateData *templateData = app.newTemplateData(r)
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
  var templateData *templateData = app.newTemplateData(r)
  templateData.Snippet = snip

  app.render(w, "view.tmpl.html", templateData, http.StatusOK)

}

func (app *application) SnipCreatePost(w http.ResponseWriter, r *http.Request) {

  var form *snippetCreateFrom = &snippetCreateFrom{}
  err := app.decodePostForm(r, form)

  if err != nil {
    app.clientError(w, http.StatusBadRequest)
    return
  }

  form.CheckField(validator.IsNotBlank(form.Title), "title" , validator.BLANK_MESSAGE)
  form.CheckField(validator.IsNotMaxChars(form.Title, 100), "title", fmt.Sprintf(validator.MAX_CHAR_MESSAGE,
    utf8.RuneCountInString(form.Title)))
  form.CheckField(validator.IsNotBlank(form.Content), "content", validator.BLANK_MESSAGE)
  form.CheckField(validator.PermittedInt(form.Expire, 1, 7, 30, 360), "expire", fmt.Sprintf(validator.NOT_IN_OPTIONS, form.Expire))


  if !form.IsValid() {
    var data *templateData = app.newTemplateData(r)
    data.Form = form

    app.render(w, "create.tmpl.html", data, http.StatusUnprocessableEntity)
    return
  }

  id, err := app.snippets.Insert(form.Title,
    form.Content,
    form.Expire)
  if err != nil {
    app.serverError(w, err)
    return
  }

  app.sessionManager.Put(r.Context(), "flash", "Snippet cree avec succes!")

  http.Redirect(w, r, fmt.Sprintf("/snip/view/%d", id), http.StatusSeeOther)

}

func (app *application) SnipCreate(w http.ResponseWriter, r *http.Request) {
  // Render html create.tmpl
  data := app.newTemplateData(r)
  data.Form = &snippetCreateFrom{
    Expire: 365,
  }
  app.render(w, "create.tmpl.html",  data, http.StatusOK)
  return
}

// Authentications

func (app *application) UserSignup(w http.ResponseWriter, r* http.Request) {

  var data *templateData = app.newTemplateData(r)
  data.Form = &userSingupForm{}
  app.render(w, "signup.tmpl.html", data, http.StatusOK)
  return
}

func (app *application) UserSignupPost(w http.ResponseWriter, r* http.Request){

  var form *userSingupForm = &userSingupForm{}

  err := app.decodePostForm(r, &form)
  if err != nil {
    app.clientError(w, http.StatusBadRequest)
    return
  }

  form.CheckField(validator.IsNotBlank(form.Name), "name", validator.BLANK_MESSAGE)
  form.CheckField(validator.IsNotBlank(form.Email), "email", validator.BLANK_MESSAGE)
  form.CheckField(validator.IsNotBlank(form.Password), "password", validator.BLANK_MESSAGE)
  form.CheckField(validator.Matches(form.Email, *validator.EmailRX),"email", validator.INVALID_EMAIL)
  form.CheckField(validator.MinChars(form.Password, 8), "password", fmt.Sprintf(validator.MIN_CHAR_MESSAGE, 8, utf8.RuneCountInString(form.Password)))

  if !form.IsValid() {
    var data *templateData = app.newTemplateData(r)
    data.Form = form
    app.render(w, "signup.tmpl.html", data, http.StatusUnprocessableEntity)
    return
  }

  err = app.users.Insert(form.Name, form.Email, form.Password)
  if err != nil {
    if errors.Is(err, models.ErrDuplicateEmail) {
      form.AddFieldError("email", validator.DUPLICATE_EMAIL)
      var data *templateData = app.newTemplateData(r)
      data.Form = form
      app.render(w, "signup.tmpl.html", data, http.StatusUnprocessableEntity)
      return
    }

    app.serverError(w, err)
    return
  }

  app.sessionManager.Put(r.Context(), "flash", "Compte cree avec succes!")
  http.Redirect( w, r, "/user/login",http.StatusSeeOther)

}

func (app *application) UserLogin(w http.ResponseWriter, r* http.Request) {

  var data *templateData = app.newTemplateData(r)

  data.Form = &userLoginForm{}
  app.render(w, "signin.tmpl.html", data, http.StatusOK)
}

func (app *application) UserLoginPost(w http.ResponseWriter, r* http.Request) {

  var form *userLoginForm = &userLoginForm{}

  err := app.decodePostForm(r, &form)

  if err != nil {
    app.clientError(w, http.StatusBadRequest)
    return
  }

  form.CheckField(validator.IsNotBlank(form.Email), "email", validator.BLANK_MESSAGE)
  form.CheckField(validator.IsNotBlank(form.Password), "password", validator.BLANK_MESSAGE)
  form.CheckField(validator.Matches(form.Email, *validator.EmailRX),"email", validator.INVALID_EMAIL)

  if !form.IsValid() {
    var templateData *templateData = app.newTemplateData(r)
    templateData.Form = form
    app.render(w, "signin.tmpl.html", templateData, http.StatusUnprocessableEntity)
    return
  }


  id, err := app.users.Authenticate(form.Email, form.Password)
  if err != nil {
    if errors.Is(err, models.ErrInvalidCredentials) {
      form.AddNonFieldError("Mot de passe ou courrier invalide")
      var templateData* templateData = app.newTemplateData(r)
      templateData.Form = form
      app.render(w, "signin.tmpl.html", templateData, http.StatusUnauthorized)
      return
    }

    app.serverError(w, err)
    return
  }

  // Renew the token on current session to change the session ID
  err = app.sessionManager.RenewToken(r.Context())
  if err != nil {
    app.serverError(w, err)
    return
  }
  app.sessionManager.Put(r.Context(), "authenticateID", id)

  http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *application) UserLogoutPost(w http.ResponseWriter, r* http.Request) {

  app.sessionManager.Remove(r.Context(), "authenticateID")

  err := app.sessionManager.RenewToken(r.Context())
  if err != nil {
    app.serverError(w, err)
    return
  }

  app.sessionManager.Put(r.Context(), "flash", "You have been logged out successfully")

  http.Redirect(w, r, "/", http.StatusSeeOther)


}


