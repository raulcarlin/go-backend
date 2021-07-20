package main

import (
	"net/http"

	"github.com/raulcarlin/go-backend/internal/adapter"
	"github.com/raulcarlin/go-backend/internal/app"
	"github.com/raulcarlin/go-backend/internal/config"
	"github.com/raulcarlin/go-backend/internal/router"
	"github.com/raulcarlin/go-backend/internal/util/logger"
	"github.com/raulcarlin/go-backend/internal/util/validator"
)

func main() {
	app := configApplication()

	startHttpServer(app)
}

func configApplication() *app.Application {
	config := config.Get()
	logger := logger.New(true)

	db, dbErr := adapter.NewGORM(config)
	if dbErr != nil {
		logger.Fatal().Err(dbErr).Msg("Error connecting DB")
	}

	db.LogMode(true)

	validator := validator.New()

	return app.New(logger, db, validator, config)
}

func startHttpServer(app *app.Application) {
	app.Logger.Info().Msgf("Starting server on %s", app.Conf.ServerConfig.Port)

	appRouter := router.New(app)

	s := &http.Server{
		Addr:         app.Conf.ServerConfig.Port,
		Handler:      appRouter,
		ReadTimeout:  app.Conf.ServerConfig.ReadTimeout,
		WriteTimeout: app.Conf.ServerConfig.WriteTimeout,
		IdleTimeout:  app.Conf.ServerConfig.IdleTimeout,
	}

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		app.Logger.Fatal().Err(err).Msg("Server startup failed")
	}
}
