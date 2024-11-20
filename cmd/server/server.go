package main

import (
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
	"github.com/alorents/lenslocked/migrations"
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
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}
	err = run(cfg)
	if err != nil {
		panic(err)
	}
}

func run(cfg config) error {
	// Load the .env file
	cfg, err := loadEnvConfig()
	if err != nil {
		return err
	}

	// Setup the postgres db
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		return err
	}
	defer db.Close()

	// run migrations
	err = models.MigrateFS(db, migrations.FS, "")
	if err != nil {
		panic(err)
	}

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
	galleriesService := &models.GalleryService{
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
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "forgot-password.gohtml"))
	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "check-your-email.gohtml"))
	usersC.Templates.ResetPassword = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "reset-password.gohtml"))

	galleriesC := controllers.GalleriesController{
		GalleryService: galleriesService,
	}
	galleriesC.Templates.New = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "galleries/new.gohtml"))
	galleriesC.Templates.Edit = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "galleries/edit.gohtml"))
	galleriesC.Templates.Index = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "galleries/index.gohtml"))
	galleriesC.Templates.Show = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "galleries/show.gohtml"))

	// Setup middleware
	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
		csrf.Path("/"),
	)
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	// Create the router and apply middleware
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(csrfMw)
	router.Use(umw.SetUser)

	// Define the routes
	router.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "layout.gohtml", "home.gohtml"))))
	router.Get("/home", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "layout.gohtml", "home.gohtml"))))
	router.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "layout.gohtml", "contact.gohtml"))))
	router.Get("/signup", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "layout.gohtml", "signup.gohtml"))))
	router.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "layout.gohtml", "faq.gohtml"))))

	// User routes
	router.Get("/signup", usersC.New)
	router.Post("/users", usersC.Create)
	router.Get("/signin", usersC.SignIn)
	router.Post("/signin", usersC.ProcessSignin)
	router.Post("/signout", usersC.ProcessSignOut)
	router.Get("/forgot-password", usersC.ForgotPassword)
	router.Post("/forgot-password", usersC.ProcessForgotPassword)
	router.Get("/reset-password", usersC.ResetPassword)
	router.Post("/reset-password", usersC.ProcessResetPassword)
	router.Route("/users/me", func(router chi.Router) {
		router.Use(umw.RequireUser)
		router.Get("/", usersC.CurrentUser)
	})

	// Gallery routes
	router.Route("/galleries", func(router chi.Router) {
		router.Get("/{id}", galleriesC.Show)
		router.Get("/{id}/images/{filename}", galleriesC.Image)
		router.Group(func(router chi.Router) {
			router.Use(umw.RequireUser)
			router.Get("/", galleriesC.Index)
			router.Post("/", galleriesC.Create)
			router.Get("/new", galleriesC.New)
			router.Get("/{id}/edit", galleriesC.Edit)
			router.Post("/{id}", galleriesC.Update)
			router.Post("/{id}/delete", galleriesC.Delete)
			router.Post("/{id}/images", galleriesC.UploadImage)
			router.Post("/{id}/images/{filename}/delete", galleriesC.DeleteImage)
		})
	})

	assetsHandler := http.FileServer(http.Dir("assets"))
	router.Get("/assets/*", http.StripPrefix("/assets", assetsHandler).ServeHTTP)

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// Start the server
	fmt.Printf("Starting the server on %s...\n", cfg.Server.Address)
	return http.ListenAndServe(cfg.Server.Address, router)
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// PSQL
	cfg.PSQL = models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		DBName:   os.Getenv("PSQL_DBNAME"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	}
	if cfg.PSQL.Host == "" || cfg.PSQL.Port == "" {
		return cfg, fmt.Errorf("no PSQL config provided")
	}
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

	// CSRF
	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	cfg.CSRF.Secure = os.Getenv("CSRF_SECURE") == "true"

	// Server
	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")

	return cfg, nil
}
