package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	adminapp "./app/admin"
	"./app/api"
	"./app/game"
	"./app/static"
	"./config"
	"./config/app"
	"./config/db"
	//"./config/bindatafs"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/qor/admin"
	"github.com/qor/qor/utils"
	//"github.com/qor/assetfs"
)

var (
	Router = chi.NewRouter()
	Admin  = admin.New(&admin.AdminConfig{
		SiteName: "Casino Admin",
		//Auth:     auth.AdminAuth{},
		DB: db.DB,
	})
	Application = app.New(&app.Config{
		Router: Router,
		Admin:  Admin,
		DB:     db.DB,
	})
)

func main() {
	Router.Use(middleware.RealIP)
	Router.Use(middleware.Logger)
	Router.Use(middleware.Recoverer)

	Application.Use(game.New(&game.Config{}))
	Application.Use(api.New(&api.Config{}))
	Application.Use(adminapp.New(&adminapp.Config{}))
	Application.Use(static.New(&static.Config{
		Prefixs: []string{"/system"},
		Handler: utils.FileServer(http.Dir(filepath.Join(config.Root, "public"))),
	}))
	/*
	Application.Use(static.New(&static.Config{
		Prefixs: []string{"javascripts", "stylesheets", "images", "dist", "fonts", "vendors", "favicon.ico"},
		Handler: bindatafs.AssetFS.FileServer(http.Dir("public"), "javascripts", "stylesheets", "images", "dist", "fonts", "vendors", "favicon.ico"),
	}))
	*/
	log.Printf("Listening on: %v\n", config.Config.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), Router); err != nil {
		log.Fatal(err)
	}
}
