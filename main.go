package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"

	"github.com/alorents/lenslocked/controllers"
	"github.com/alorents/lenslocked/models"
	"github.com/alorents/lenslocked/templates"
	"github.com/alorents/lenslocked/views"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func main() {
	// Load the .env file
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}
	// Setup the postgres db
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	// Setup the services
	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB: db,
	}
	pwResetService := &models.PasswordResetService{
		DB: db,
	}
	emailService := models.NewEmailService(cfg.SMTP)

	// Setup the controllers
	usersC := controllers.UsersController{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "signup.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "signin.gohtml"))
	usersC.Templates.Profile = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "profile.gohtml"))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "forgot-pw.gohtml"))
	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "check-your-email.gohtml"))

	// Setup middleware
	csrfMW := csrf.Protect([]byte(cfg.CSRF.Key), csrf.Secure(cfg.CSRF.Secure))
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
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
	fmt.Printf("Starting the server on %s...\n", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, router)
	if err != nil {
		panic(err)
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// TODO PSQL - read from env
	cfg.PSQL = models.DefaultPostgresConfig()
	// SMTP
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}
	cfg.SMTP.Port = port
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")
	// TODO CSRF - read from env
	cfg.CSRF.Key = "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX" // TODO fix before deploying to production
	cfg.CSRF.Secure = false
	// TODO the server values from env
	cfg.Server.Address = ":3000"

	return cfg, nil
}
