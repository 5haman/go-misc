package app

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/middlewares"
	"github.com/qor/wildcard_router"
)

type MicroAppInterface interface {
	ConfigureApplication(*Application)
}

type Application struct {
	*Config
}

type Config struct {
	Router   *chi.Mux
	Handlers []http.Handler
	Admin    *admin.Admin
	DB       *gorm.DB
}

func New(cfg *Config) *Application {
	if cfg == nil {
		cfg = &Config{}
	}
	if cfg.Router == nil {
		cfg.Router = chi.NewRouter()
	}

	return &Application{
		Config: cfg,
	}
}

func (application *Application) Use(app MicroAppInterface) {
	app.ConfigureApplication(application)
}

func (application *Application) NewServeMux() http.Handler {
	if len(application.Config.Handlers) == 0 {
		return middlewares.Apply(application.Config.Router)
	}

	wildcardRouter := wildcard_router.New()
	for _, handler := range application.Config.Handlers {
		wildcardRouter.AddHandler(handler)
	}
	wildcardRouter.AddHandler(application.Config.Router)

	return middlewares.Apply(wildcardRouter)
}
