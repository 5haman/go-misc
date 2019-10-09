package model

import (
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/qor/media"
	"github.com/qor/media/oss"
	"github.com/qor/publish2"
	qor_seo "github.com/qor/seo"
	"github.com/qor/sorting"
	"github.com/qor/validations"
)

type Game struct {
	gorm.Model

	Type            string
	Name            string `gorm:"index:name;unique;not null"`
	Vendor          string `gorm:"index:vendor;not null"`
	RatingValue     float32
	RatingCount     int
	Paylines        int
	Reels           int
	MinCoinsPerLine int
	MaxCoinsPerLine int
	MinCoinsSize    float32
	MaxCoinsSize    float32
	Jackpot         float32
	RTP             float32
	Object          string `sql:"type:text"`
	Description     string `sql:"type:text"`
	Preview         Image  `sql:"size:4294967295;" media_library:"url:/system_new/{{class}}/{{primary_key}}/{{column}}.{{extension}}"`
	Seo             qor_seo.Setting

	sorting.SortingDESC
	publish2.Schedule
	publish2.Visible
}

type Image struct {
	oss.OSS
}

func (Image) GetSizes() map[string]*media.Size {
	return map[string]*media.Size{
		"sd": {720, 480, false},
		"preview": {480, 240, false},
	}
}

func (game Game) DefaultPath() string {
	defaultPath := "/"
	return defaultPath
}

func (game Game) ImageURL() string {
	return game.Preview.URL("thumb")
}

func (game Game) Validate(db *gorm.DB) {
	if strings.TrimSpace(game.Name) == "" {
		db.AddError(validations.NewError(game, "Name", "Name can not be empty"))
	}
}

func (game Game) GetSEO() *qor_seo.SEO {
	return SEOCollection.GetSEO("Game Page")
}
