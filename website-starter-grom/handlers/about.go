package handlers

import (
	"net/http"
	"website-starter/core/logging"
	"website-starter/core/templating"
	"website-starter/services"
)

type AboutHandler struct {
	service   *services.AboutService
	templates *templating.Engine
	logger    *logging.Logger
	errors    *ErrorHandler
}

func NewAboutHandler(
	service *services.AboutService,
	templates *templating.Engine,
	logger *logging.Logger,
	errors *ErrorHandler,
) *AboutHandler {
	return &AboutHandler{
		service:   service,
		templates: templates,
		logger:    logger,
		errors:    errors,
	}
}

func (handler *AboutHandler) GetAbout(
	w http.ResponseWriter,
	r *http.Request,
) {
	message := handler.service.GetWelcomeMessage()

	data := map[string]any{
		"Title":   "About",
		"Message": message,
	}

	err := handler.templates.Render(
		w,
		"about",
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
