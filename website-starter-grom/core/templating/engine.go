package templating

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

// Engine manages template loading, caching,
// and rendering for the application.
type Engine struct {
	cache map[string]*template.Template
	mutex sync.RWMutex
}

// NewTemplateEngine creates and initializes
// a new template engine instance.
func NewTemplateEngine() *Engine {
	return &Engine{
		cache: map[string]*template.Template{},
	}
}

// Load scans the templates directory,
// parses all templates and partials,
// and stores them in memory cache.
func (engine *Engine) Load() error {

	// Base templates directory.
	const path = "templates"

	// Template file extension.
	const suffix = "html"

	// Main layout template filename.
	const layout = "layout." + suffix

	// Find all page templates.
	pages, err := filepath.Glob(
		fmt.Sprintf(
			"%s/*.%s",
			path,
			suffix,
		),
	)

	if err != nil {
		return err
	}

	// Find all partial templates.
	partials, err := filepath.Glob(
		fmt.Sprintf(
			"%s/partials/*.%s",
			path,
			suffix,
		),
	)

	if err != nil {
		return err
	}

	// Lock cache for writing.
	engine.mutex.Lock()

	defer engine.mutex.Unlock()

	// Parse every page template.
	for _, page := range pages {

		// Skip the main layout file.
		if strings.Contains(
			page,
			layout,
		) {
			continue
		}

		// Extract template name from filename.
		name := strings.TrimSuffix(
			filepath.Base(
				page,
			),
			fmt.Sprintf(
				".%s",
				suffix,
			),
		)

		// Build template file list.
		files := []string{
			fmt.Sprintf(
				"%s/%s",
				path,
				layout,
			),
			page,
		}

		// Include all partial templates.
		files = append(
			files,
			partials...,
		)

		// Parse all template files together.
		parsedFile, err := template.ParseFiles(
			files...,
		)

		if err != nil {
			return err
		}

		// Store parsed template in cache.
		engine.cache[name] = parsedFile
	}

	return nil
}

// Render executes a cached template
// and writes the output to the response writer.
func (engine *Engine) Render(
	w http.ResponseWriter,
	name string,
	data any,
) error {
	// Lock cache for reading.
	engine.mutex.RLock()

	parsedFile, found := engine.cache[name]

	engine.mutex.RUnlock()

	// Return error if template does not exist.
	if !found {
		return fmt.Errorf(
			"template '%s' not found",
			name,
		)
	}

	// Use buffer before writing to response.
	var buffer bytes.Buffer

	// Execute the layout template.
	err := parsedFile.ExecuteTemplate(
		&buffer,
		"layout",
		data,
	)

	if err != nil {
		return err
	}

	// Write rendered content to response.
	_, err = buffer.WriteTo(
		w,
	)

	if err != nil {
		return err
	}

	return nil
}
