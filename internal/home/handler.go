package home

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/Nurlan270/weather-go/internal/home/view"
	"github.com/Nurlan270/weather-go/internal/location"
	"github.com/Nurlan270/weather-go/internal/logger"
	"github.com/Nurlan270/weather-go/internal/renderer"
	"github.com/Nurlan270/weather-go/internal/user"
)

type Handler struct {
	locationSvc *location.Service
	renderer    *renderer.Renderer
	log         *logger.Logger
}

func NewHandler(
	locationSvc *location.Service,
	renderer *renderer.Renderer,
	log *logger.Logger,
) *Handler {
	return &Handler{
		locationSvc: locationSvc,
		renderer:    renderer,
		log:         log,
	}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	const VIEW = "home"

	u := user.FromRequest(r)

	locations, err := h.locationSvc.ListLocationsByUserID(u.ID)
	if err != nil {
		h.log.Error("list locations by user id failed", zap.Error(err))

		_ = h.renderer.RenderServerError(w, u.Login)

		return
	}

	data := view.NewHomeViewData("Home", u.Login, locations)

	if err = h.renderer.RenderView(w, http.StatusOK, VIEW, data); err != nil {
		h.log.ErrorRenderPage(err)
		return
	}
}
