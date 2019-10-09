package config

import (
	"log"
	"os"

	"github.com/jinzhu/configor"
	//"github.com/qor/oss/filesystem"
	"github.com/qor/redirect_back"
	"github.com/qor/session/manager"
	"github.com/unrolled/render"
)

var Config = struct {
	Port uint `default:"80" env:"PORT"`
	DB   struct {
		Name     string `env:"DBName"`
		Adapter  string `env:"DBAdapter" default:"mysql"`
		Host     string `env:"DBHost" default:"localhost"`
		Port     string `env:"DBPort" default:"3306"`
		User     string `env:"DBUser"`
		Password string `env:"DBPassword"`
	}
}{}

var (
	Root         = os.Getenv("HOME") + "/Projects/casino"
	Render       = render.New()
	RedirectBack = redirect_back.New(&redirect_back.Config{
		SessionManager:  manager.SessionManager,
		IgnoredPrefixes: []string{"/auth"},
	})
)

func init() {
	if err := configor.Load(&Config, "config/database.yml"); err != nil {
		log.Fatal(err)
	}
}
