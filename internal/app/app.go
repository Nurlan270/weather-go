package app

import (
	"database/sql"
	"errors"
	"github.com/Nurlan270/weather-go/internal/config"
	"github.com/Nurlan270/weather-go/internal/db"
	"github.com/Nurlan270/weather-go/internal/home"
	"github.com/Nurlan270/weather-go/internal/logger"
	"github.com/Nurlan270/weather-go/internal/renderer"
	mw "github.com/Nurlan270/weather-go/internal/rest/middleware"
	"github.com/Nurlan270/weather-go/internal/user"
	"github.com/Nurlan270/weather-go/internal/validator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	gpv "github.com/go-playground/validator/v10"
	"html/template"
	"log"
	"net/http"
	"net/mail"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"go.uber.org/zap"
)

type App struct {
	router    chi.Router
	config    *config.Config
	logger    *logger.Logger
	db        *sql.DB
	renderer  *renderer.Renderer
	validator *validator.Validate
	services  Services
}

type Services struct {
	User *user.Service
}

func Run() {
	var app App

	//	Setup application
	app.setup()

	//	Start server
	app.startServer()
}

func (a *App) setup() {
	//	Initializing app components
	a.initConfig()
	a.initLogger()
	a.initDatabase()
	a.initRenderer()
	a.initServices()
	a.initValidator()
	a.initRouter()

	//	Registering app components
	a.registerRoutes()
}

func (a *App) initConfig() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	a.config = cfg
}

func (a *App) initLogger() {
	a.logger = logger.New(a.config.App.Env)
}

func (a *App) initDatabase() {
	database, err := db.Connect(a.config.DB)
	if err != nil {
		a.logger.Fatal("failed to connect to database", zap.Error(err))
	}

	a.db = database
}

func (a *App) initRenderer() {
	templates := make(map[string]*template.Template)

	mainDir, _ := filepath.Glob("view/pages/*.html")
	subDirs, _ := filepath.Glob("view/pages/*/*.html")

	if mainDir == nil && subDirs == nil {
		a.logger.Fatal("no template files found")
	}

	allFiles := append(mainDir, subDirs...)

	//	Parse view
	for _, tmplPath := range allFiles {
		tmpl, err := template.ParseFiles(
			"view/index.html", // Root template file
			"view/components/header.html",
			tmplPath,
		)
		if err != nil {
			a.logger.Fatal("failed to parse template", zap.Error(err))
		}

		name := strings.TrimSuffix(strings.TrimPrefix(
			tmplPath, "view/pages/",
		), ".html")

		templates[name] = tmpl
	}

	a.renderer = renderer.New(templates)
}

func (a *App) initRouter() {
	r := chi.NewRouter()

	//	Global middlewares
	loggerMW := mw.NewLoggerMiddleware(a.logger)

	r.Use(
		middleware.Recoverer,
		middleware.CleanPath,
		middleware.RequestID,
		loggerMW.Log,
	)

	//	404-Page custom setup
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		err := a.renderer.RenderNotFound(w, r)
		if err != nil {
			a.logger.ErrorRenderPage(err)
		}
	})

	a.router = r
}

func (a *App) initValidator() {
	v := gpv.New(gpv.WithRequiredStructEnabled())

	//	Register custom validation rules
	_ = v.RegisterValidation("login", func(fl gpv.FieldLevel) bool {
		var re = regexp.MustCompile("^[a-zA-Z0-9._-]+$")

		login := fl.Field().String()

		//	Check whether it's a valid email, if yes - return true
		if _, err := mail.ParseAddress(login); err == nil {
			return true
		}

		// Otherwise, it's a username: reject if it contains forbidden chars
		return re.MatchString(login)
	})

	a.validator = validator.New(v)
}

func (a *App) initServices() {
	//	Repositories
	userRepo := user.NewDBRepository(a.db)

	//	Services
	a.services.User = user.NewService(userRepo, a.config.Session)
}

func (a *App) registerRoutes() {
	//	Static files
	a.router.Handle(
		"/static/*",
		http.StripPrefix("/static/", static("view/static")),
	)

	//	Handlers
	homeH := home.NewHandler(a.renderer, a.logger)
	userH := user.NewHandler(a.services.User, a.validator, a.renderer, a.config.Session, a.logger)

	//	Middlewares
	authMW := mw.NewAuthMiddleware(a.services.User, a.renderer, a.config.Session, a.logger)
	guestMW := mw.NewGuestMiddleware(a.services.User, a.renderer, a.config.Session, a.logger)
	limiterMW := mw.NewRateLimitMiddleware(a.renderer)

	//	Routes
	a.router.Route("/", func(r chi.Router) {
		r.Use(authMW.RequireAuth)

		r.Get("/", homeH.Index)

		//	Sign out
		a.router.With(
			middleware.ClientIPFromRemoteAddr,
			limiterMW.Limit(30, time.Hour),
		).Post("/auth/sign-out", userH.LogoutUser)
	})

	//	Auth (only guests allowed)
	a.router.Route("/auth", func(r chi.Router) {
		r.Use(guestMW.RequireGuest)

		r.Get("/sign-up", userH.ShowRegister)
		r.Get("/sign-in", userH.ShowLogin)

		//	Rate-limited routes
		r.Route("/", func(r chi.Router) {
			r.Use(
				middleware.ClientIPFromRemoteAddr,
				limiterMW.Limit(30, time.Hour),
			)

			r.Post("/sign-up", userH.RegisterUser)
			r.Post("/sign-in", userH.LoginUser)
		})
	})
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

func (a *App) startServer() {
	server := &http.Server{
		Handler:      a.router,
		Addr:         a.config.HTTPServer.Address,
		ReadTimeout:  a.config.HTTPServer.Timeout,
		WriteTimeout: a.config.HTTPServer.Timeout,
		IdleTimeout:  a.config.HTTPServer.IdleTimeout,
	}

	a.logger.Info("Starting server", zap.String("address", server.Addr))
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.logger.Fatal("App finished unexpectedly", zap.Error(err))
	}
}
