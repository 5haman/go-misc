package api

import (
	"../../config/app"
	"../../config/db"
	"../../model"
	"github.com/qor/admin"
	"github.com/qor/qor"
)

func New(config *Config) *App {
	if config.Prefix == "" {
		config.Prefix = "/api"
	}
	return &App{Config: config}
}

type App struct {
	Config *Config
}
type Config struct {
	Prefix string
}

func (app App) ConfigureApplication(application *app.Application) {
	API := admin.New(&qor.Config{DB: db.DB})

	API.AddResource(&model.Game{})

	application.Router.Mount(app.Config.Prefix, API.NewServeMux(app.Config.Prefix))
}
