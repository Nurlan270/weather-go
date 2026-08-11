package home

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/Nurlan270/weather-go/internal/home/view"
	"github.com/Nurlan270/weather-go/internal/location"
	"github.com/Nurlan270/weather-go/internal/logger"
	"github.com/Nurlan270/weather-go/internal/renderer"
	"github.com/Nurlan270/weather-go/internal/rest/request"
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

	user := request.GetUserFromCtx(r.Context())

	locations, err := h.locationSvc.ListLocationsByUserID(user.ID)
	if err != nil {
		h.log.Error("list locations by user id failed", zap.Error(err))

		_ = h.renderer.RenderServerError(w, r)

		return
	}

	data := view.NewHomeViewData("Home", user.Login, locations)

	if err = h.renderer.RenderView(w, http.StatusOK, VIEW, data); err != nil {
		h.log.ErrorRenderPage(err)
		return
	}
}
