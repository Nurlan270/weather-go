package location

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/Nurlan270/weather-go/internal/location/view"
	"github.com/Nurlan270/weather-go/internal/logger"
	"github.com/Nurlan270/weather-go/internal/renderer"
	"github.com/Nurlan270/weather-go/internal/rest/openweather/response"
	"github.com/Nurlan270/weather-go/internal/rest/request"
	"github.com/Nurlan270/weather-go/internal/validator"
	baseview "github.com/Nurlan270/weather-go/internal/view"
)

type Handler struct {
	service  *Service
	validate *validator.Validate
	renderer *renderer.Renderer
	log      *logger.Logger
}

func NewHandler(
	service *Service,
	validate *validator.Validate,
	renderer *renderer.Renderer,
	log *logger.Logger,
) *Handler {
	return &Handler{
		service:  service,
		validate: validate,
		renderer: renderer,
		log:      log,
	}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	const VIEW = "search"

	user := request.GetUserFromCtx(r.Context())

	var location *response.Location

	data := view.NewSearchViewData("Search Results", user.Login, location)

	query := strings.TrimSpace(r.URL.Query().Get("q"))

	data.Old["query"] = query

	reqLocation := request.SearchLocation{
		Query: query,
	}

	//	Validate data
	if err := h.validate.Struct(reqLocation); err != nil {
		h.renderer.Redirect(w, r, "/")
		return
	}

	location, err := h.service.SearchLocation(reqLocation)
	if errors.Is(err, ErrNoResults) {
		data.Error.Message = fmt.Sprintf(baseview.MessageNoResults, query)

		if err = h.renderer.RenderView(w, http.StatusNotFound, VIEW, data); err != nil {
			h.log.ErrorRenderPage(err)
		}

		return
	}

	if err != nil {
		h.log.Error("could not get search results", zap.Error(err))

		if err = h.renderer.RenderServerError(w, r); err != nil {
			h.log.ErrorRenderPage(err)
		}

		return
	}

	data.Location = location
	if err = h.renderer.RenderView(w, http.StatusOK, VIEW, data); err != nil {
		h.log.ErrorRenderPage(err)
		return
	}
}

func (h *Handler) AddLocation(w http.ResponseWriter, r *http.Request) {
	const VIEW = "search"

	user := request.GetUserFromCtx(r.Context())

	var location *response.Location

	data := view.NewSearchViewData("Search Results", user.Login, location)

	if err := r.ParseForm(); err != nil {
		h.log.Error("could not parse form", zap.Error(err))

		data.Error.Message = baseview.MessageInvalidData

		if err = h.renderer.RenderView(w, http.StatusBadRequest, VIEW, data); err != nil {
			h.log.ErrorRenderPage(err)
		}

		return
	}

	name := strings.TrimSpace(r.PostForm.Get("name"))
	lat := strings.TrimSpace(r.PostForm.Get("lat"))
	lon := strings.TrimSpace(r.PostForm.Get("lon"))

	reqLocation := request.AddLocation{
		Name: name,
		Lat:  lat,
		Lon:  lon,
	}

	//	Validate data
	if err := h.validate.Struct(reqLocation); err != nil {
		data.Error.Message = baseview.MessageValidationFailed

		if err = h.renderer.RenderView(w, http.StatusUnprocessableEntity, VIEW, data); err != nil {
			h.log.ErrorRenderPage(err)
		}

		return
	}

	//	Add new location
	err := h.service.AddLocation(user.ID, reqLocation)

	//	Location already exists
	if errors.Is(err, ErrLocationAlreadyExists) {
		data.Error.Message = baseview.MessageLocationAlreadyAdded

		if err = h.renderer.RenderView(w, http.StatusUnprocessableEntity, VIEW, data); err != nil {
			h.log.ErrorRenderPage(err)
		}

		return
	}

	if err != nil {
		h.log.Error("could not add location", zap.Error(err))

		if err = h.renderer.RenderServerError(w, r); err != nil {
			h.log.ErrorRenderPage(err)
		}

		return
	}

	//	Location added -> redirect user to home page
	h.renderer.Redirect(w, r, "/")
}

func (h *Handler) DeleteLocation(w http.ResponseWriter, r *http.Request) {
	user := request.GetUserFromCtx(r.Context())
	name := chi.URLParam(r, "name")

	if err := h.service.DeleteLocation(user.ID, name); err != nil {
		h.log.Error("could not delete location", zap.Error(err))

		if err = h.renderer.RenderServerError(w, r); err != nil {
			h.log.ErrorRenderPage(err)
		}

		return
	}

	h.renderer.Redirect(w, r, "/")
}
