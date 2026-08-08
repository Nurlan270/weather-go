package home

import (
	"github.com/Nurlan270/weather-go/internal/config"
	"github.com/Nurlan270/weather-go/internal/logger"
	"github.com/Nurlan270/weather-go/internal/renderer"
	"github.com/Nurlan270/weather-go/internal/rest/request"
	"github.com/Nurlan270/weather-go/internal/user"
	"github.com/Nurlan270/weather-go/internal/view"
	"net/http"
)

type Handler struct {
	userSvc  *user.Service
	renderer *renderer.Renderer
	sessCfg  config.Session
	log      *logger.Logger
}

func NewHandler(
	renderer *renderer.Renderer,
	log *logger.Logger,
) *Handler {
	return &Handler{
		renderer: renderer,
		log:      log,
	}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	const VIEW = "home"

	login := request.GetLoginFromCtx(r.Context())

	data := view.NewBaseViewData("Home", login)

	if err := h.renderer.RenderView(w, http.StatusOK, VIEW, data); err != nil {
		h.log.ErrorRenderPage(err)
		return
	}
}
