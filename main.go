package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/alorents/lenslocked/views"
)

func executeTemplate(w http.ResponseWriter, filepath string) {
	tpl, err := views.Parse(filepath)
	if err != nil {
		log.Printf("error parsing template: %v", err)
		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return
	}
	tpl.Execute(w, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "templates/home.gohtml")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "templates/contact.gohtml")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "templates/faq.gohtml")

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
