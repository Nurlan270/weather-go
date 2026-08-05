package controller

import (
	"github.com/Nurlan270/weather-go/internal/logger"
	"github.com/Nurlan270/weather-go/internal/renderer"
	"github.com/Nurlan270/weather-go/internal/renderer/response"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const PAGE = "home"

type HomeController struct {
	renderer *renderer.Renderer
	log      *logger.Logger
}

func NewHomeController(renderer *renderer.Renderer, log *logger.Logger) *HomeController {
	return &HomeController{
		renderer: renderer,
		log:      log,
	}
}

func (c *HomeController) RegisterRoutes(r *chi.Mux) {
	r.Get("/", c.Index)
}

func (c *HomeController) Index(w http.ResponseWriter, r *http.Request) {
	resp := response.PageResponse{
		PageTitle: "Home",
	}

	if err := c.renderer.RenderPage(w, r, PAGE, resp); err != nil {
		c.log.ErrorRenderPage(err)
		return
	}
}
