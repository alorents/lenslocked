package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/gorilla/csrf"

	"github.com/alorents/lenslocked/context"
	"github.com/alorents/lenslocked/models"
)

type UsersController struct {
	Templates struct {
		New            Template
		SignIn         Template
		Profile        Template
		ForgotPassword Template
		CheckEmail     Template
	}
	UserService          *models.UserService
	SessionService       *models.SessionService
	PasswordResetService *models.PasswordResetService
	EmailService         *models.EmailService
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
	setCookie(w, CookeSession, session.Token)
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
	setCookie(w, CookeSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
	return
}

func (c UsersController) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := readCooke(r, CookeSession)
	if err != nil || tokenCookie.Value == "" {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	err = c.SessionService.DeleteByToken(tokenCookie.Value)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unexpepcted error", http.StatusInternalServerError)
		return
	}

	deleteCookie(w, CookeSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
	return
}

func (c UsersController) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())

	var data struct {
		User *models.User
	}
	data.User = user
	c.Templates.Profile.Execute(w, r, data)
}

func (c UsersController) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	c.Templates.ForgotPassword.Execute(w, r, data)
}

func (c UsersController) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	pwResetService, err := c.PasswordResetService.Create(data.Email)
	if err != nil {
		// TODO handle other cases in the future. For instance if a user doesn't exist with that email
		fmt.Println(err)
		http.Error(w, "Unexpepcted error", http.StatusInternalServerError)
		return
	}

	vals := url.Values{
		"token": {pwResetService.Token},
	}
	resetURL := "http://localhost:3000/reset-password?" + vals.Encode()
	err = c.EmailService.ForgotPassword(data.Email, resetURL)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unexpepcted error", http.StatusInternalServerError)
		return
	}

	// Don't render the token to the user. We need them to verify their email address
	c.Templates.CheckEmail.Execute(w, r, data)
}
