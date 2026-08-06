package home

import (
	"github.com/Nurlan270/weather-go/internal/logger"
	"github.com/Nurlan270/weather-go/internal/renderer"
	"github.com/Nurlan270/weather-go/internal/renderer/response"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Handler struct {
	renderer *renderer.Renderer
	log      *logger.Logger
}

func NewHandler(renderer *renderer.Renderer, log *logger.Logger) *Handler {
	return &Handler{
		renderer: renderer,
		log:      log,
	}
}

func (h *Handler) RegisterRoutes(r *chi.Mux) {
	r.Get("/", h.Index)
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	const PAGE = "home"

	resp := response.PageResponse{
		PageTitle: "Home",
	}

	if err := h.renderer.RenderPage(w, PAGE, resp); err != nil {
		h.log.ErrorRenderPage(err)
		return
	}
}
