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
	// Setup the postgres db
	postgresConfig := models.DefaultPostgresConfig()
	// TODO not prod safe
	fmt.Println(postgresConfig)
	db, err := models.Open(postgresConfig)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Setup the services
	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}

	// Setup the controllers
	usersC := controllers.UsersController{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "signup.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "signin.gohtml"))
	usersC.Templates.Profile = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "profile.gohtml"))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "forgot-pw.gohtml"))

	// Setup middleware
	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX" // TODO fix before deploying to production≈
	csrfMW := csrf.Protect([]byte(csrfKey), csrf.Secure(false))
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	// Create the router and apply middleware
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(csrfMW)
	router.Use(umw.SetUser)

	// Define the routes
	router.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "layout.gohtml", "home.gohtml"))))
	router.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "layout.gohtml", "contact.gohtml"))))
	router.Get("/signup", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "layout.gohtml", "signup.gohtml"))))
	router.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "layout.gohtml", "faq.gohtml"))))

	// User routes
	router.Get("/signup", usersC.New)
	router.Post("/users", usersC.Create)
	router.Get("/signin", usersC.SignIn)
	router.Post("/signin", usersC.ProcessSignin)
	router.Post("/signout", usersC.ProcessSignOut)
	router.Get("/forgot-pw", usersC.ForgotPassword)
	router.Post("/forgot-pw", usersC.ProcessForgotPassword)
	router.Route("/users/me", func(router chi.Router) {
		router.Use(umw.RequireUser)
		router.Get("/", usersC.CurrentUser)
	})

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// Start the server
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", router)
}
