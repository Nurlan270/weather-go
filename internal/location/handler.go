package location

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/Nurlan270/weather-go/internal/location/view"
	"github.com/Nurlan270/weather-go/internal/logger"
	"github.com/Nurlan270/weather-go/internal/renderer"
	"github.com/Nurlan270/weather-go/internal/rest/request"
	"github.com/Nurlan270/weather-go/internal/user"
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

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	const VIEW = "search"

	u := user.FromRequest(r)

	data := view.NewSearchViewData("Search Results", u.Login, nil)

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

	locations, err := h.service.SearchLocations(reqLocation)
	if errors.Is(err, ErrNoResults) {
		data.Error.Message = fmt.Sprintf(baseview.MessageNoResults, query)

		if err = h.renderer.RenderView(w, http.StatusNotFound, VIEW, data); err != nil {
			h.log.ErrorRenderPage(err)
		}

		return
	}

	if err != nil {
		h.log.Error("could not get search results", zap.Error(err))

		if err = h.renderer.RenderServerError(w, u.Login); err != nil {
			h.log.ErrorRenderPage(err)
		}

		return
	}

	data.Locations = locations
	if err = h.renderer.RenderView(w, http.StatusOK, VIEW, data); err != nil {
		h.log.ErrorRenderPage(err)
		return
	}
}

func (h *Handler) AddLocation(w http.ResponseWriter, r *http.Request) {
	const VIEW = "search"

	u := user.FromRequest(r)

	data := view.NewSearchViewData("Search Results", u.Login, nil)

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
	err := h.service.AddLocation(u.ID, reqLocation)

	if err != nil && !errors.Is(err, ErrLocationAlreadyExists) {
		h.log.Error("could not add location", zap.Error(err))

		if err = h.renderer.RenderServerError(w, u.Login); err != nil {
			h.log.ErrorRenderPage(err)
		}

		return
	}

	//	Location added -> redirect user to home page
	h.renderer.Redirect(w, r, "/")
}

func (h *Handler) DeleteLocation(w http.ResponseWriter, r *http.Request) {
	u := user.FromRequest(r)

	int64Id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.log.Error("could not parse id", zap.Error(err))
		return
	}

	if err = h.service.DeleteLocation(int64Id, u); err != nil {
		h.log.Error("could not delete location", zap.Error(err))

		if err = h.renderer.RenderServerError(w, u.Login); err != nil {
			h.log.ErrorRenderPage(err)
		}

		return
	}

	h.renderer.Redirect(w, r, "/")
}
