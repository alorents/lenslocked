package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"

	"github.com/alorents/lenslocked/controllers"
	"github.com/alorents/lenslocked/models"
	"github.com/alorents/lenslocked/templates"
	"github.com/alorents/lenslocked/views"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	// TODO fix before deploying to production
	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	csrfMW := csrf.Protect([]byte(csrfKey), csrf.Secure(false))
	router.Use(csrfMW)

	router.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "layout.gohtml", "home.gohtml"))))
	router.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "layout.gohtml", "contact.gohtml"))))
	router.Get("/signup", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "layout.gohtml", "signup.gohtml"))))
	router.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "layout.gohtml", "faq.gohtml"))))

	postgresConfig := models.DefaultPostgresConfig()
	db, err := models.Open(postgresConfig)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create the services
	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}

	// Create the controllers
	usersC := controllers.UsersController{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "signup.gohtml"))
	router.Get("/signup", usersC.New)
	router.Post("/users", usersC.Create)
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "signin.gohtml"))
	router.Get("/signin", usersC.SignIn)
	router.Post("/signin", usersC.ProcessSignin)
	router.Get("/users/me", usersC.CurrentUser)

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", router)
}
