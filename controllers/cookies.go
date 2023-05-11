package controllers

import (
	"net/http"
)

const (
	CookeSession = "session"
)

func newCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}
}

func setCookie(w http.ResponseWriter, name, value string) {
	http.SetCookie(w, newCookie(name, value))
}

func readCooke(r *http.Request, name string) (*http.Cookie, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return nil, err
	}
	return cookie, nil
}
