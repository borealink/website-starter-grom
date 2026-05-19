package handlers

import (
	"net/http"
	"website-starter/core/logging"
	"website-starter/core/templating"
	"website-starter/services"
)

type HomeHandler struct {
	service   *services.HomeService
	templates *templating.Engine
	logger    *logging.Logger
	errors    *ErrorHandler
}

func NewHomeHandler(
	service *services.HomeService,
	templates *templating.Engine,
	logger *logging.Logger,
	errors *ErrorHandler,
) *HomeHandler {
	return &HomeHandler{
		service:   service,
		templates: templates,
		logger:    logger,
		errors:    errors,
	}
}

func (handler *HomeHandler) GetHome(
	w http.ResponseWriter,
	r *http.Request,
) {
	message := handler.service.GetWelcomeMessage()

	notes, err := handler.service.GetNotes()

	if err != nil {
		handler.errors.InternalServerError(
			w,
			r,
			err,
		)

		return
	}

	data := map[string]any{
		"Title":   "Home",
		"Message": message,
		"Notes":   notes,
	}

	err = handler.templates.Render(
		w,
		"home",
		data,
	)

	if err != nil {
		handler.errors.InternalServerError(
			w,
			r,
			err,
		)

		return
	}
}

func (handler *HomeHandler) CreateNote(
	w http.ResponseWriter,
	r *http.Request,
) {
	err := r.ParseForm()

	if err != nil {
		handler.errors.InternalServerError(
			w,
			r,
			err,
		)

		return
	}

	title := r.FormValue(
		"title",
	)

	description := r.FormValue(
		"description",
	)

	if title == "" || description == "" {
		http.Redirect(
			w,
			r,
			"/",
			http.StatusSeeOther,
		)

		return
	}

	err = handler.service.CreateNote(
		title,
		description,
	)

	if err != nil {
		handler.errors.InternalServerError(
			w,
			r,
			err,
		)

		return
	}

	http.Redirect(
		w,
		r,
		"/",
		http.StatusSeeOther,
	)
}

func (handler *HomeHandler) DeleteNote(
	w http.ResponseWriter,
	r *http.Request,
) {
	err := r.ParseForm()

	if err != nil {
		handler.errors.InternalServerError(
			w,
			r,
			err,
		)

		return
	}

	id := r.FormValue(
		"id",
	)

	if id == "" {
		http.Redirect(
			w,
			r,
			"/",
			http.StatusSeeOther,
		)

		return
	}

	err = handler.service.DeleteNote(
		id,
	)

	if err != nil {
		handler.errors.InternalServerError(
			w,
			r,
			err,
		)

		return
	}

	http.Redirect(
		w,
		r,
		"/",
		http.StatusSeeOther,
	)
}

func (handler *HomeHandler) UpdateNote(
	w http.ResponseWriter,
	r *http.Request,
) {
	err := r.ParseForm()

	if err != nil {
		handler.errors.InternalServerError(
			w,
			r,
			err,
		)

		return
	}

	id := r.FormValue(
		"id",
	)

	title := r.FormValue(
		"title",
	)

	description := r.FormValue(
		"description",
	)

	if id == "" || title == "" || description == "" {
		http.Redirect(
			w,
			r,
			"/",
			http.StatusSeeOther,
		)

		return
	}

	err = handler.service.UpdateNote(
		id,
		title,
		description,
	)

	if err != nil {
		handler.errors.InternalServerError(
			w,
			r,
			err,
		)

		return
	}

	http.Redirect(
		w,
		r,
		"/",
		http.StatusSeeOther,
	)
}
