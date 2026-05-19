package http

import (
	"fmt"
	"net/http"
)

// Route represents a single HTTP route definition.
type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

// Router is a lightweight wrapper around http.ServeMux.
type Router struct {
	mux      *http.ServeMux
	notFound http.HandlerFunc
}

// NewRouter creates and initializes a new router instance.
func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

// RegisterRoutes registers all application routes.
func (router *Router) RegisterRoutes(
	routes []Route,
) {
	for _, route := range routes {
		// Create local copies to avoid closure issues.
		path := route.Path
		method := route.Method
		handler := route.Handler

		// Register the route handler.
		router.mux.HandleFunc(
			path,
			func(
				w http.ResponseWriter,
				req *http.Request,
			) {
				// Ensure the requested path matches exactly.
				if req.URL.Path != path {
					// Use custom not found handler if available.
					if router.notFound != nil {
						router.notFound(
							w,
							req,
						)
					} else {
						// Fallback to default 404 response.
						http.NotFound(
							w,
							req,
						)
					}

					return
				}

				// Validate the HTTP method.
				if req.Method != method {
					w.WriteHeader(
						http.StatusMethodNotAllowed,
					)

					return
				}

				// Execute the route handler.
				handler(
					w,
					req,
				)
			},
		)
	}
}

// RegisterStatic serves static files from a directory.
func (router *Router) RegisterStatic(
	path,
	dir string,
) {
	// Create a file server for the static directory.
	fs := http.FileServer(
		http.Dir(
			dir,
		),
	)

	// Register the static route.
	router.mux.Handle(
		fmt.Sprintf(
			"GET %s",
			path,
		),
		http.StripPrefix(
			path,
			fs,
		),
	)
}

// SetNotFound sets a custom 404 handler.
func (router *Router) SetNotFound(
	handler http.HandlerFunc,
) {
	router.notFound = handler
}

// HttpHandler returns the internal HTTP handler.
func (router *Router) HttpHandler() http.Handler {
	return router.mux
}
