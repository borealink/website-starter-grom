package main

import (
	"log"
	"os"
	"time"

	core_database "website-starter/core/database"
	core_http "website-starter/core/http"
	"website-starter/core/logging"
	"website-starter/core/templating"
	"website-starter/handlers"
	"website-starter/models"
	"website-starter/services"
)

func main() {

	// =========================
	// LOGGER
	// =========================

	logger := logging.NewLogger(logging.Config{
		Out:    os.Stdout,
		Prefix: "[APP] ",
		Flags:  log.LstdFlags,
	})

	// =========================
	// DATABASE
	// =========================

	database, err := core_database.NewDatabase(
		"app.db",
	)

	if err != nil {
		logger.Get().Fatal(
			err,
		)
	}

	// =========================
	// AUTO MIGRATION
	// =========================

	err = database.Get().AutoMigrate(
		&models.Note{},
	)

	if err != nil {
		logger.Get().Fatal(
			err,
		)
	}

	// =========================
	// TEMPLATES
	// =========================

	templates := templating.NewTemplateEngine()

	err = templates.Load()

	if err != nil {
		logger.Get().Fatal(
			err,
		)
	}

	// =========================
	// SERVICES
	// =========================

	homeService := services.NewHomeService(
		database.Get(),
	)

	aboutService := services.NewAboutService()

	// =========================
	// HANDLERS
	// =========================

	errorHandler := handlers.NewErrorHandler(
		templates,
		logger,
	)

	homeHandler := handlers.NewHomeHandler(
		homeService,
		templates,
		logger,
		errorHandler,
	)

	aboutHandler := handlers.NewAboutHandler(
		aboutService,
		templates,
		logger,
		errorHandler,
	)

	// =========================
	// ROUTER
	// =========================

	router := core_http.NewRouter()

	router.RegisterStatic(
		"/static/",
		"./static",
	)

	router.RegisterRoutes([]core_http.Route{
		{
			Method:  "GET",
			Path:    "/",
			Handler: homeHandler.GetHome,
		},
		{
			Method:  "POST",
			Path:    "/notes/create",
			Handler: homeHandler.CreateNote,
		},
		{
			Method:  "POST",
			Path:    "/notes/update",
			Handler: homeHandler.UpdateNote,
		},
		{
			Method:  "POST",
			Path:    "/notes/delete",
			Handler: homeHandler.DeleteNote,
		},
		{
			Method:  "GET",
			Path:    "/about",
			Handler: aboutHandler.GetAbout,
		},
	})

	router.SetNotFound(
		errorHandler.NotFound,
	)

	handler := router.HttpHandler()

	// =========================
	// SERVER
	// =========================

	server := core_http.NewServer(
		":8081",
		handler,
		logger,
	)

	server.Start()

	server.WaitForShutdown(
		10 * time.Second,
	)
}
