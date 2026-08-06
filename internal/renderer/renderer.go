package renderer

import (
	"fmt"
	"github.com/Nurlan270/weather-go/internal/renderer/message"
	"github.com/Nurlan270/weather-go/internal/renderer/response"
	"html/template"
	"net/http"
)

type Renderer struct {
	templates map[string]*template.Template
}

func New(templates map[string]*template.Template) *Renderer {
	return &Renderer{
		templates: templates,
	}
}

func (r *Renderer) RenderPage(
	w http.ResponseWriter,
	page string,
	resp response.PageResponse,
) error {
	return r.render(w, http.StatusOK, page, resp)
}

func (r *Renderer) RenderNotFound(w http.ResponseWriter) error {
	resp := response.ErrorPageResponse{
		PageTitle: message.NotFound,
		Message:   message.NotFound,
	}

	return r.render(w, http.StatusNotFound, "error", resp)
}

func (r *Renderer) Back(
	w http.ResponseWriter,
	code int,
	page string,
	resp response.PageResponse,
) error {
	return r.render(w, code, page, resp)
}

// render is low-level helper; powers other methods
func (r *Renderer) render(
	w http.ResponseWriter,
	statusCode int,
	page string,
	resp any,
) error {
	tmpl, ok := r.templates[page]
	if !ok {
		return fmt.Errorf("template not found: %q", page)
	}

	r.writeMetaData(w, statusCode)

	if err := tmpl.Execute(w, resp); err != nil {
		return fmt.Errorf("failed to execute template %q: %w", page, err)
	}
	return nil
}

func (r *Renderer) writeMetaData(w http.ResponseWriter, statusCode int) {
	//	Headers
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	//	Status code
	w.WriteHeader(statusCode)
}
