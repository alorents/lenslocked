package controllers

import (
	"fmt"
	"net/http"

	"github.com/alorents/lenslocked/models"
)

type UsersController struct {
	Templates struct {
		New    Template
		SignIn Template
	}
	UserService *models.UserService
}

func (c UsersController) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	c.Templates.New.Execute(w, data)
}

func (c UsersController) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := c.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unexpepcted error", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Created user %v", *user)
}

func (c UsersController) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	c.Templates.SignIn.Execute(w, data)
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
	cookie := http.Cookie{
		Name:  "email",
		Value: user.Email,
		Path:  "/",
	}
	http.SetCookie(w, &cookie)
	fmt.Fprintf(w, "User authenticated: %+v", user)
}
