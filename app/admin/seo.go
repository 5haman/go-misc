package admin

import (
	"../../model"
	"github.com/qor/admin"
	qor_seo "github.com/qor/seo"
)

func SetupSEO(Admin *admin.Admin) {
	model.SEOCollection = qor_seo.New("Common SEO")
	model.SEOCollection.RegisterGlobalVaribles(&model.SEOGlobalSetting{SiteName: "Casino"})
	model.SEOCollection.SettingResource = Admin.AddResource(&model.MySEOSetting{}, &admin.Config{Invisible: true})
	model.SEOCollection.RegisterSEO(&qor_seo.SEO{
		Name: "Home Page",
	})

	model.SEOCollection.RegisterSEO(&qor_seo.SEO{
		Name:     "Games",
		Varibles: []string{"Name"},
		Context: func(objects ...interface{}) map[string]string {
			game := objects[0].(model.Game)
			context := make(map[string]string)
			context["Name"] = game.Name
			return context
		},
	})

	Admin.AddResource(model.SEOCollection, &admin.Config{Name: "SEO Settings", Menu: []string{"Site Management"}, Singleton: true, Priority: 2})
}
