package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jinzhu/gorm"
	"github.com/raulcarlin/go-backend/internal/model"
	"github.com/raulcarlin/go-backend/internal/repository"
	"github.com/raulcarlin/go-backend/internal/util/validator"
)

func (app *Application) HandleListBooks(w http.ResponseWriter, r *http.Request) {
	books, err := repository.ListBooks(app.DB)
	if err != nil {
		app.Logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrDataAccessFailure)
		return
	}
	if books == nil {
		fmt.Fprint(w, "[]")
		return
	}
	dtos := books.ToDto()
	if err := json.NewEncoder(w).Encode(dtos); err != nil {
		app.Logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)
		return
	}
}

func (app *Application) HandleCreateBook(w http.ResponseWriter, r *http.Request) {
	form := &model.BookForm{}

	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		app.Logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrFormDecodingFailure)
		return
	}

	if err := app.Validator.Struct(form); err != nil {
		app.Logger.Warn().Err(err).Msg("")
		resp := validator.ToErrResponse(err)
		if resp == nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, appErrFormErrResponseFailure)
			return
		}
		respBody, err := json.Marshal(resp)
		if err != nil {
			app.Logger.Warn().Err(err).Msg("")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(respBody)
		return
	}

	bookModel, err := form.ToModel()
	if err != nil {
		app.Logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrFormDecodingFailure)
		return
	}

	book, err := repository.CreateBook(app.DB, bookModel)
	if err != nil {
		app.Logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrDataCreationFailure)
		return
	}
	app.Logger.Info().Msgf("New book created: %d", book.ID)
	w.WriteHeader(http.StatusCreated)
}

func (app *Application) HandleReadBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		app.Logger.Info().Msgf("can not parse ID: %v", id)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	book, err := repository.ReadBook(app.DB, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		app.Logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrDataAccessFailure)
		return
	}

	dto := book.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		app.Logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)
		return
	}
}

func (app *Application) HandleUpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		app.Logger.Info().Msgf("can not parse ID: %v", id)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	form := &model.BookForm{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		app.Logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrFormDecodingFailure)
		return
	}

	if err := app.Validator.Struct(form); err != nil {
		app.Logger.Warn().Err(err).Msg("")
		resp := validator.ToErrResponse(err)
		if resp == nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, appErrFormErrResponseFailure)
			return
		}
		respBody, err := json.Marshal(resp)
		if err != nil {
			app.Logger.Warn().Err(err).Msg("")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(respBody)
		return
	}

	bookModel, err := form.ToModel()
	if err != nil {
		app.Logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrFormDecodingFailure)
		return
	}

	bookModel.ID = uint(id)
	if err := repository.UpdateBook(app.DB, bookModel); err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		app.Logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrDataUpdateFailure)
		return
	}

	app.Logger.Info().Msgf("Book updated: %d", id)

	w.WriteHeader(http.StatusAccepted)
}

func (app *Application) HandleDeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		app.Logger.Info().Msgf("can not parse ID: %v", id)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := repository.DeleteBook(app.DB, uint(id)); err != nil {
		app.Logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrDataAccessFailure)
		return
	}

	app.Logger.Info().Msgf("Book deleted: %d", id)

	w.WriteHeader(http.StatusAccepted)
}
