package static

import (
	"net/http"
	"strings"

	"../../config/app"
)

func New(config *Config) *App {
	return &App{Config: config}
}

type App struct {
	Config *Config
}

type Config struct {
	Prefixs []string
	Handler http.Handler
}

func (app App) ConfigureApplication(application *app.Application) {
	for _, prefix := range app.Config.Prefixs {
		application.Router.Mount("/"+strings.TrimPrefix(prefix, "/"), app.Config.Handler)
	}
}
