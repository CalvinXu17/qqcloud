// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"verification/controllers"
	"verification/controllers/admin"
	"verification/controllers/api"
)

func init() {
	beego.Include(&controllers.MainController{})
	nsAdmin := beego.NewNamespace("/admin",
		beego.NSNamespace("/user",
			beego.NSInclude(&admin.UserController{}),
		),
		beego.NSNamespace("/project",
			beego.NSInclude(&admin.ProjectController{}),
		),
	)
	beego.AddNamespace(nsAdmin)
	nsApi := beego.NewNamespace("/api",
		beego.NSNamespace("/index",
			beego.NSInclude(&api.IndexController{}),
		),
	)
	beego.AddNamespace(nsApi)
}
