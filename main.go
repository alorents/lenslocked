package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Contact Page</h1><p>To get in touch, email me at <a href=\"mailto:andreasphoenix@gmail.com\">andreasphoenix@gmail.com</a>.</p>")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>FAQ Page</h1>"+
		"<p>"+
		"Q: Why does this exist?</br>"+
		"A: To learn WebDev stuff!</br></br>"+
		"Q: How do we contact you?</br>"+
		"A: See the <a href=\"/contact\">contact us</a> page</br></br>"+
		"</p>")
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	w.Write([]byte(fmt.Sprintf("<h1>User Page</h1><p>User: %s</p>", userID)))
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", homeHandler)
	router.Get("/contact", contactHandler)
	router.Get("/faq", faqHandler)
	router.Route("/users/{userID}", func(router chi.Router) {
		router.Use(middleware.RequestID)
		router.Use(middleware.RealIP)
		router.Use(middleware.Logger)
		router.Get("/", userHandler)
	})
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", router)
}
