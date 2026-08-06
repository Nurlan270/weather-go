package app

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"net/mail"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Nurlan270/weather-go/internal/config"
	"github.com/Nurlan270/weather-go/internal/db"
	"github.com/Nurlan270/weather-go/internal/home"
	"github.com/Nurlan270/weather-go/internal/logger"
	"github.com/Nurlan270/weather-go/internal/renderer"
	"github.com/Nurlan270/weather-go/internal/rest"
	"github.com/Nurlan270/weather-go/internal/user"
	"github.com/Nurlan270/weather-go/internal/validator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	gpv "github.com/go-playground/validator/v10"

	"go.uber.org/zap"
)

type App struct {
	router    *chi.Mux
	config    *config.Config
	logger    *logger.Logger
	db        *sql.DB
	renderer  *renderer.Renderer
	validator *validator.Validate
}

func Run() error {
	app := &App{}

	//	Config
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	app.config = cfg

	//	Logger
	app.logger = logger.New(app.config.Env)

	//	Database
	database, err := db.Connect(app.config.DB)
	if err != nil {
		return err
	}
	app.db = database

	//	Renderer
	rend, err := app.setupRenderer()
	if err != nil {
		return err
	}
	app.renderer = rend

	//	Router
	app.router = app.setupRouter()

	//	Validator
	app.validator = app.setupValidator()

	//	Handlers
	handlers := app.registerHandlers()

	app.registerRoutes(handlers...)

	//	Start server
	if err = app.startServer(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (a *App) startServer() error {
	server := &http.Server{
		Handler:      a.router,
		Addr:         a.config.HTTPServer.Address,
		ReadTimeout:  a.config.HTTPServer.Timeout,
		WriteTimeout: a.config.HTTPServer.Timeout,
		IdleTimeout:  a.config.HTTPServer.IdleTimeout,
	}

	a.logger.Info("Starting server", zap.String("address", server.Addr))
	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (a *App) setupRouter() *chi.Mux {
	router := chi.NewRouter()

	//	Middlewares
	router.Use(
		middleware.URLFormat,
		middleware.Recoverer,
		middleware.CleanPath,
	)

	//	404-Page custom setup
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		err := a.renderer.RenderNotFound(w)
		if err != nil {
			a.logger.ErrorRenderPage(err)
		}
	})

	return router
}

func (a *App) registerHandlers() []rest.RouteRegistrar {
	//	Repositories
	userRepo := user.NewDBRepository(a.db)

	//	Services
	userSvc := user.NewService(userRepo)

	//	Handlers
	return []rest.RouteRegistrar{
		home.NewHandler(a.renderer, a.logger),
		user.NewHandler(userSvc, a.validator, a.renderer, a.logger),
	}
}

func (a *App) registerRoutes(handlers ...rest.RouteRegistrar) {
	//	Static files
	a.router.Handle(
		"/static/*",
		http.StripPrefix("/static/", static("templates/static")),
	)

	//	Routes
	for _, h := range handlers {
		h.RegisterRoutes(a.router)
	}
}

func (a *App) setupRenderer() (*renderer.Renderer, error) {
	templates := make(map[string]*template.Template)

	tmplFiles, _ := filepath.Glob("templates/pages/*.html")
	if tmplFiles == nil {
		return nil, fmt.Errorf("no template files found")
	}

	//	Parse templates
	for _, tmplPath := range tmplFiles {
		tmpl, err := template.ParseFiles(
			"templates/index.html", // Root template file
			"templates/components/header.html",
			tmplPath,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template: %w", err)
		}

		name := strings.TrimSuffix(strings.TrimPrefix(
			tmplPath, "templates/pages/",
		), ".html")

		templates[name] = tmpl
	}

	return renderer.New(templates), nil
}

func (a *App) setupValidator() *validator.Validate {
	validate := gpv.New(gpv.WithRequiredStructEnabled())

	//	Register custom validation rules
	validate.RegisterValidation("login", func(fl gpv.FieldLevel) bool {
		var re = regexp.MustCompile("^[a-zA-Z0-9._-]+$")

		login := fl.Field().String()

		//	Check whether it's a valid email, if yes - return true
		if _, err := mail.ParseAddress(login); err == nil {
			return true
		}

		// Otherwise, it's a username: reject if it contains forbidden chars
		return re.MatchString(login)
	})

	return validator.New(validate)
}

func static(dir string) http.Handler {
	fs := http.FileServer(http.Dir(dir))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//	Return 404 if user try to access list of static files (e.g. static/css/)
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		fs.ServeHTTP(w, r)
	})
}
