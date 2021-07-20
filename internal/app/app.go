package app

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"github.com/raulcarlin/go-backend/internal/config"
	"github.com/raulcarlin/go-backend/internal/util/logger"
)

const (
	appErrDataAccessFailure      = "data access failure"
	appErrJsonCreationFailure    = "json creation failure"
	appErrDataCreationFailure    = "data creation failure"
	appErrFormDecodingFailure    = "form decoding failure"
	appErrDataUpdateFailure      = "data update failure"
	appErrFormErrResponseFailure = "form error response failure"
)

// Application dependencies
type Application struct {
	Logger    *logger.Logger
	DB        *gorm.DB
	Validator *validator.Validate
	Conf      *config.Config
}

// Dependency Injection
func New(logger *logger.Logger, db *gorm.DB, validator *validator.Validate, config *config.Config) *Application {
	return &Application{Logger: logger, DB: db, Conf: config, Validator: validator}
}

func (app *Application) HandleIndex(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Write([]byte("Sample App Backend"))
}
