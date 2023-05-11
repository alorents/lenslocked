package controllers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"

	"github.com/alorents/lenslocked/models"
)

type UsersController struct {
	Templates struct {
		New    Template
		SignIn Template
	}
	UserService    *models.UserService
	SessionService *models.SessionService
}

func (c UsersController) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email     string
		CSRFField template.HTML
	}
	data.Email = r.FormValue("email")
	data.CSRFField = csrf.TemplateField(r)
	c.Templates.New.Execute(w, r, data)
}

func (c UsersController) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := c.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
		return
	}

	session, err := c.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		// TODO - display error to user
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    session.Token,
		Path:     "/",
		HttpOnly: true,
	})
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (c UsersController) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	c.Templates.SignIn.Execute(w, r, data)
}

func (c UsersController) ProcessSignin(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := c.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unexpepcted error", http.StatusInternalServerError)
		return
	}
	session, err := c.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unexpepcted error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    session.Token,
		Path:     "/",
		HttpOnly: true,
	})
	fmt.Fprintf(w, "User authenticated: %+v", user)
}

func (c UsersController) CurrentUser(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("session")
	if err != nil || tokenCookie.Value == "" {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	user, err := c.SessionService.User(tokenCookie.Value)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	fmt.Fprintf(w, "Current user: %+v", user)
}

func (c UsersController) SignOut(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:   "email",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
	fmt.Fprintln(w, "Signed out")
}
