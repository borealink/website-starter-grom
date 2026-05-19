package handlers

import (
	"net/http"

	"website-starter/core/logging"
	"website-starter/core/templating"
)

type ErrorHandler struct {
	templates *templating.Engine
	logger    *logging.Logger
}

func NewErrorHandler(
	templates *templating.Engine,
	logger *logging.Logger,
) *ErrorHandler {
	return &ErrorHandler{
		templates: templates,
		logger:    logger,
	}
}

func (handler *ErrorHandler) NotFound(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.WriteHeader(
		http.StatusNotFound,
	)

	data := map[string]any{
		"Title": "404",
	}

	err := handler.templates.Render(
		w,
		"404",
		data,
	)

	if err != nil {
		http.Error(
			w,
			"404 page not found",
			http.StatusNotFound,
		)

		return
	}
}

func (handler *ErrorHandler) InternalServerError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	handler.logger.Get().Println(
		err,
	)

	w.WriteHeader(
		http.StatusInternalServerError,
	)

	data := map[string]any{
		"Title": "500",
	}

	renderErr := handler.templates.Render(
		w,
		"500",
		data,
	)

	if renderErr != nil {
		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError,
		)

		return
	}
}
