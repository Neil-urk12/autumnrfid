package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"

	"admin/internal/config"
	"admin/internal/handlers"
	"admin/repositories"
)

var templates *template.Template

func main() {
	// Load db config
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		log.Fatalf("Failed to load database config: %v", err)
	}

	// init db toma
	db, err := repositories.NewDatabaseClient(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Test the database connection
	if err := testDBConnection(db.DB); err != nil {
		log.Fatalf("Database connection test failed: %v", err)
	}
	log.Println("Database connection test successful.")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("Shutting down...")
		db.Close()
		os.Exit(0)
	}()
	defer db.Close()
	mux := http.NewServeMux()
	initTemplates()

	appHandler := handlers.NewHandler(db)
	// Routes
	// mux.HandleFunc("/", handlers.IndexHandler(templates))
	mux.HandleFunc("/", appHandler.IndexHandler(templates))
	mux.HandleFunc("/login", appHandler.LoginPageHandler(templates))
	mux.HandleFunc("/dashboard", appHandler.DashboardHandler(templates))
	mux.HandleFunc("/students", appHandler.StudentsPageHandler(templates))
	// mux.HandleFunc("/login", handlers.LoginPageHandler(templates))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static/"))))
	mux.HandleFunc("/logout", handlers.LogoutHandler())
	// mux.HandleFunc("/students", handlers.StudentsPageHandler(templates))
	// mux.HandleFunc("/dashboard", handlers.DashboardHandler(templates))

	fmt.Println("Server running on port: 8080")
	http.ListenAndServe(":8080", mux)
}

func initTemplates() {
	// Parse all templates at startup
	var err error
	templates, err = template.ParseFiles(
		"ui/html/templates/index.html",
		"ui/html/templates/login.html",
		"ui/html/partials/dashboard.html",
		"ui/html/partials/student_management.html",
	)
	if err != nil {
		// Prevent the admin from running if there are template errors
		panic(fmt.Sprintf("Failed to parse templates: %v", err))
	}
}

func testDBConnection(db *sql.DB) error {
	// Ping the database
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping failed: %v", err)
	}

	// Try a simple query
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Students").Scan(&count)
	if err != nil {
		return fmt.Errorf("test query failed: %v", err)
	}
	fmt.Printf("Found %d students in database\n", count)

	return nil
}
