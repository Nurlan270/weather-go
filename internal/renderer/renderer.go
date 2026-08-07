package renderer

import (
	"fmt"
	"github.com/Nurlan270/weather-go/internal/view"
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

func (r *Renderer) RenderView(
	w http.ResponseWriter,
	code int,
	page string,
	data any,
) error {
	return r.render(w, code, page, data)
}

func (r *Renderer) RenderNotFound(w http.ResponseWriter) error {
	data := view.NewErrorViewData(view.MessageNotFoundError, view.MessageNotFoundError, "")

	return r.render(w, http.StatusNotFound, "error", data)
}

func (r *Renderer) RenderServerError(w http.ResponseWriter) error {
	data := view.NewErrorViewData(view.MessageServerError, view.MessageServerError, "")

	return r.render(w, http.StatusInternalServerError, "error", data)
}

func (r *Renderer) Back(
	w http.ResponseWriter,
	code int,
	page string,
	data any,
) error {
	return r.render(w, code, page, data)
}

func (r *Renderer) Redirect(w http.ResponseWriter, req *http.Request, url string) {
	r.writeHeaders(w)

	http.Redirect(w, req, url, http.StatusSeeOther)
}

// render is low-level helper; powers other methods
func (r *Renderer) render(
	w http.ResponseWriter,
	statusCode int,
	page string,
	data any,
) error {
	tmpl, ok := r.templates[page]
	if !ok {
		return fmt.Errorf("template not found: %q", page)
	}

	r.writeMetaData(w, statusCode)

	if err := tmpl.Execute(w, data); err != nil {
		return fmt.Errorf("failed to execute template %q: %w", page, err)
	}
	return nil
}

func (r *Renderer) writeHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

// writeMetaData is used to set headers & write status code
func (r *Renderer) writeMetaData(w http.ResponseWriter, statusCode int) {
	//	Headers
	r.writeHeaders(w)

	//	Status code
	w.WriteHeader(statusCode)
}
