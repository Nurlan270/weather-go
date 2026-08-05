package app

import (
	"database/sql"
	"fmt"
	"github.com/Nurlan270/weather-go/internal/rest"
	"github.com/Nurlan270/weather-go/internal/rest/controller"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Nurlan270/weather-go/internal/config"
	"github.com/Nurlan270/weather-go/internal/db"
	"github.com/Nurlan270/weather-go/internal/logger"
	"github.com/Nurlan270/weather-go/internal/renderer"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type App struct {
	router   *chi.Mux
	config   *config.Config
	logger   *logger.Logger
	db       *sql.DB
	renderer *renderer.Renderer
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

	//	Controllers
	controllers := app.registerController()

	app.registerRoutes(controllers...)

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

func (a *App) registerController() []rest.RouteRegistrar {
	//	Storages

	//	Services

	//	Controllers
	return []rest.RouteRegistrar{
		controller.NewHomeController(a.renderer, a.logger),
	}
}

func (a *App) registerRoutes(controllers ...rest.RouteRegistrar) {
	//	Static files
	a.router.Handle(
		"/static/*",
		http.StripPrefix("/static/", static("templates/static")),
	)

	//	Routes
	for _, c := range controllers {
		c.RegisterRoutes(a.router)
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
