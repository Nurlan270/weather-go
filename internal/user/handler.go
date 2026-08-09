package user

import (
	"errors"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/Nurlan270/weather-go/internal/config"
	"github.com/Nurlan270/weather-go/internal/logger"
	"github.com/Nurlan270/weather-go/internal/renderer"
	"github.com/Nurlan270/weather-go/internal/rest"
	"github.com/Nurlan270/weather-go/internal/rest/request"
	"github.com/Nurlan270/weather-go/internal/validator"
	"github.com/Nurlan270/weather-go/internal/view"
)

type Handler struct {
	service  *Service
	validate *validator.Validate
	renderer *renderer.Renderer
	sessCfg  config.Session
	log      *logger.Logger
}

func NewHandler(
	service *Service,
	validate *validator.Validate,
	renderer *renderer.Renderer,
	sessCfg config.Session,
	log *logger.Logger,
) *Handler {
	return &Handler{
		service:  service,
		validate: validate,
		renderer: renderer,
		sessCfg:  sessCfg,
		log:      log,
	}
}

func (h *Handler) ShowRegister(w http.ResponseWriter, r *http.Request) {
	const VIEW = "auth/register"

	data := view.NewBaseViewData("Sign Up", "")

	if err := h.renderer.RenderView(w, http.StatusOK, VIEW, data); err != nil {
		h.log.ErrorRenderPage(err)
		return
	}
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	const VIEW = "auth/register"

	data := view.NewBaseViewData("Sign Up", "")

	if err := r.ParseForm(); err != nil {
		h.log.Error("could not parse form", zap.Error(err))

		data.Error.Message = view.MessageInvalidData

		if err = h.renderer.RenderView(w, http.StatusBadRequest, VIEW, data); err != nil {
			h.log.ErrorRenderPage(err)
		}

		return
	}

	login := strings.TrimSpace(r.PostForm.Get("login"))
	password := strings.TrimSpace(r.PostForm.Get("password"))
	passwordConfirmation := strings.TrimSpace(r.PostForm.Get("password_confirmation"))

	user := request.RegisterUser{
		Login:                login,
		Password:             password,
		PasswordConfirmation: passwordConfirmation,
	}

	//	Validate data
	if err := h.validate.Struct(user); err != nil {
		data.Old["login"] = login
		data.Error.Message = view.MessageValidationFailed
		data.Error.Items = h.validate.MapErrors(err)

		if err = h.renderer.RenderView(w, http.StatusUnprocessableEntity, VIEW, data); err != nil {
			h.log.ErrorRenderPage(err)
		}

		return
	}

	//	Create user
	session, err := h.service.RegisterUser(user)

	//	User already exists
	if errors.Is(err, ErrUserAlreadyExists) {
		data.Old["login"] = login
		data.Error.Message = view.MessageUserAlreadyExists

		if err = h.renderer.RenderView(w, http.StatusUnprocessableEntity, VIEW, data); err != nil {
			h.log.ErrorRenderPage(err)
		}

		return
	}

	if err != nil {
		h.log.Error("could not register user", zap.Error(err))

		data.Old["login"] = login
		data.Error.Message = view.MessageServerError

		if err = h.renderer.RenderView(w, http.StatusInternalServerError, VIEW, data); err != nil {
			h.log.ErrorRenderPage(err)
		}

		return
	}

	//	Set session cookie
	rest.SetCookie(w,
		rest.GetSessionCookieName(h.sessCfg.Name), session.ID, session.ExpiresAt)

	h.renderer.Redirect(w, r, "/")
}

func (h *Handler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	const VIEW = "auth/login"

	data := view.NewBaseViewData("Sign In", "")

	if err := h.renderer.RenderView(w, http.StatusOK, VIEW, data); err != nil {
		h.log.ErrorRenderPage(err)
		return
	}
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	const VIEW = "auth/login"

	data := view.NewBaseViewData("Sign In", "")

	if err := r.ParseForm(); err != nil {
		h.log.Error("could not parse form", zap.Error(err))

		data.Error.Message = view.MessageInvalidData

		if err = h.renderer.RenderView(w, http.StatusBadRequest, VIEW, data); err != nil {
			h.log.ErrorRenderPage(err)
		}

		return
	}

	login := strings.TrimSpace(r.PostForm.Get("login"))
	password := strings.TrimSpace(r.PostForm.Get("password"))

	user := request.LoginUser{
		Login:    login,
		Password: password,
	}

	//	Validate data
	if err := h.validate.Struct(user); err != nil {
		data.Old["login"] = login
		data.Error.Message = view.MessageValidationFailed
		data.Error.Items = h.validate.MapErrors(err)

		if err = h.renderer.RenderView(w, http.StatusUnprocessableEntity, VIEW, data); err != nil {
			h.log.ErrorRenderPage(err)
		}

		return
	}

	//	Login user
	session, err := h.service.LoginUser(user)

	//	User doesn't exists
	if errors.Is(err, ErrUserNotFound) {
		data.Old["login"] = login
		data.Error.Message = view.MessageUserNotFound

		if err = h.renderer.RenderView(w, http.StatusUnprocessableEntity, VIEW, data); err != nil {
			h.log.ErrorRenderPage(err)
		}

		return
	}

	//	Invalid password
	if errors.Is(err, ErrInvalidPassword) {
		data.Old["login"] = login
		data.Error.Message = view.MessageInvalidPassword

		if err = h.renderer.RenderView(w, http.StatusUnprocessableEntity, VIEW, data); err != nil {
			h.log.ErrorRenderPage(err)
		}

		return
	}

	if err != nil {
		h.log.Error("could not login user", zap.Error(err))

		data.Old["login"] = login
		data.Error.Message = view.MessageServerError

		if err = h.renderer.RenderView(w, http.StatusInternalServerError, VIEW, data); err != nil {
			h.log.ErrorRenderPage(err)
		}

		return
	}

	//	Set session cookie
	rest.SetCookie(w,
		rest.GetSessionCookieName(h.sessCfg.Name), session.ID, session.ExpiresAt)

	h.renderer.Redirect(w, r, "/")
}

func (h *Handler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	rest.DeleteCookie(w, rest.GetSessionCookieName(h.sessCfg.Name))

	h.renderer.Redirect(w, r, "/")
}
