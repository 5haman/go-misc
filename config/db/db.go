package db

import (
	"errors"
	"fmt"
	"log"
	"os"

	"../../config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/qor/media"
	"github.com/qor/publish2"
	//"github.com/qor/sorting"
	//"github.com/qor/validations"
)

var DB *gorm.DB

func init() {
	var err error

	dbConfig := config.Config.DB
	if config.Config.DB.Adapter == "mysql" {
		DB, err = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name))
		DB = DB.Set("gorm:table_options", "CHARSET=utf8")
	} else {
		log.Fatal(errors.New("not supported database adapter"))
	}

	if err == nil {
		if os.Getenv("DEBUG") != "" {
			DB.LogMode(true)
		}

		//sorting.RegisterCallbacks(DB)
		//validations.RegisterCallbacks(DB)
		media.RegisterCallbacks(DB)
		publish2.RegisterCallbacks(DB)

		DB = DB.Set(publish2.VisibleMode, publish2.ModeOff).Set(publish2.ScheduleMode, publish2.ModeOff)
	} else {
		log.Fatal(err)
	}
}
