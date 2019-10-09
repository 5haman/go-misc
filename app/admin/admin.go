package admin

import (
	"../../config/app"
	"github.com/qor/action_bar"
	"github.com/qor/admin"
)

type App struct {
	Config *Config
}

type Config struct {
	Prefix string
}

var (
	ActionBar    *action_bar.ActionBar
	AssetManager *admin.Resource
)

func New(config *Config) *App {
	if config.Prefix == "" {
		config.Prefix = "/admin"
	}
	return &App{Config: config}
}

func (app App) ConfigureApplication(application *app.Application) {
	Admin := application.Admin

	ActionBar = action_bar.New(Admin)
	ActionBar.RegisterAction(&action_bar.Action{Name: "Admin Dashboard", Link: "/admin"})

	SetupSEO(Admin)

	application.Router.Mount(app.Config.Prefix, Admin.NewServeMux(app.Config.Prefix))
}
