package game

import (
	"../../config/app"
	"../../model"
	"github.com/qor/admin"
)

func New(config *Config) *App {
	return &App{Config: config}
}

type App struct {
	Config *Config
}

type Config struct {
}

func (app App) ConfigureApplication(application *app.Application) {
	application.DB.AutoMigrate(&model.Game{})
	//controller := &Controller{View: render.New(&render.Config{AssetFileSystem: application.AssetFS.NameSpace("blog")}, "app/pages/views")}

	//funcmapmaker.AddFuncMapMaker(controller.View)
	app.ConfigureAdmin(application.Admin)
	//application.Router.Get("/games", controller.Index)
}

func (App) ConfigureAdmin(Admin *admin.Admin) {
	Admin.AddMenu(&admin.Menu{Name: "Games Management", Priority: 4})
	Admin.AddResource(&model.Game{}, &admin.Config{Menu: []string{"Games Management"}})
}
