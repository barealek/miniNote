package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/barealek/mininote/internal/controllers"
	"github.com/barealek/mininote/internal/database"
	"github.com/barealek/mininote/pkg/shutdown"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Prepare the exit code and defer the exit
	var exitCode int

	defer func() {
		os.Exit(exitCode)
	}()

	// Run the application
	cleanup, err := run()
	defer cleanup()

	// Handle any startup error
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		exitCode = 1
		return
	}

	// Wait for the application to exit and exit gracefully
	shutdown.Shutdown()
}

func run() (func(), error) {
	router, cleanup, err := buildMiniNote()
	if err != nil {
		return nil, err
	}

	// Run the app in a goroutine
	go func() {
		http.ListenAndServe(":3000", router)
	}()

	return func() {
		cleanup()
	}, nil
}

func buildMiniNote() (*chi.Mux, func(), error) {

	// Connect to the MongoDB database
	db, err := database.BootstrapMongo(os.Getenv("MONGO_CONN_URI"), "miniLink", 10*time.Second)
	if err != nil {
		return nil, nil, err
	}

	app := chi.NewRouter()
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)

	controllers.RegisterControllers(app)

	return app, func() {
		database.CloseMongo(db)
	}, nil
}
