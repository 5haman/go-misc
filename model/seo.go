package model

import (
	"github.com/qor/seo"
)

type MySEOSetting struct {
	seo.QorSEOSetting
}

type SEOGlobalSetting struct {
	SiteName string
}

var SEOCollection *seo.Collection
