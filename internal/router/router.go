package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/raulcarlin/go-backend/internal/app"
	"github.com/raulcarlin/go-backend/internal/util/logger"
)

func New(application *app.Application) *chi.Mux {
	l := application.Logger

	r := chi.NewRouter()

	// Welcome Page
	r.Method("GET", "/", logger.NewHandler(application.HandleIndex, l))

	// Routes for health check
	r.Get("/health", app.HandleLive)
	r.Method("GET", "/health/ready", logger.NewHandler(application.HandleReady, l))

	// Routes for APIs
	r.Route("/api/v1", func(r chi.Router) {
		// All JSON
		r.Use(ContentTypeJson)

		r.Method("GET", "/books", logger.NewHandler(application.HandleListBooks, l))
		r.Method("POST", "/books", logger.NewHandler(application.HandleCreateBook, l))
		r.Method("GET", "/books/{id}", logger.NewHandler(application.HandleReadBook, l))
		r.Method("PUT", "/books/{id}", logger.NewHandler(application.HandleUpdateBook, l))
		r.Method("DELETE", "/books/{id}", logger.NewHandler(application.HandleDeleteBook, l))
	})

	return r
}
