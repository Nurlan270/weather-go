package user

import (
	"github.com/Nurlan270/weather-go/internal/logger"
	"github.com/Nurlan270/weather-go/internal/renderer"
	"github.com/Nurlan270/weather-go/internal/renderer/message"
	"github.com/Nurlan270/weather-go/internal/renderer/response"
	"github.com/Nurlan270/weather-go/internal/rest/request"
	"github.com/Nurlan270/weather-go/internal/validator"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"strings"
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

func (h *Handler) RegisterRoutes(r *chi.Mux) {
	r.Route("/auth", func(r chi.Router) {
		r.Get("/register", h.ShowRegister)
		r.Post("/register", h.RegisterUser)
	})
}

func (h *Handler) ShowRegister(w http.ResponseWriter, r *http.Request) {
	const PAGE = "register"

	resp := response.PageResponse{
		PageTitle: "Sign Up",
	}

	if err := h.renderer.RenderPage(w, PAGE, resp); err != nil {
		h.log.ErrorRenderPage(err)
		return
	}
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	const PAGE = "register"

	resp := response.PageResponse{
		PageTitle: "Sign Up",
		OldData:   make(response.OldData),
	}

	if err := r.ParseForm(); err != nil {
		h.log.Error("could not parse form", zap.Error(err))

		resp.Error.Message = message.InvalidData

		if err = h.renderer.Back(w, http.StatusBadRequest, PAGE, resp); err != nil {
			h.log.ErrorRenderPage(err)
		}
		return
	}

	login := strings.TrimSpace(r.PostForm.Get("login"))
	password := strings.TrimSpace(r.PostForm.Get("password"))
	passwordConfirmation := strings.TrimSpace(r.PostForm.Get("password_confirmation"))

	user := request.User{
		Login:                login,
		Password:             password,
		PasswordConfirmation: passwordConfirmation,
	}

	//	Validate data
	if err := h.validate.Struct(user); err != nil {
		resp.OldData["login"] = login
		resp.Error.Message = message.ValidationFailed
		resp.Error.Items = h.validate.MapErrors(err)

		h.log.Debug("debug", zap.Any("error", err))

		if err = h.renderer.Back(w, http.StatusUnprocessableEntity, PAGE, resp); err != nil {
			h.log.ErrorRenderPage(err)
		}
		return
	}

	//	Register user
	//	...
}
