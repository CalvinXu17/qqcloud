package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

	beego.GlobalControllerRouter["verification/controllers/api:IndexController"] = append(beego.GlobalControllerRouter["verification/controllers/api:IndexController"],
		beego.ControllerComments{
			Method:           "Index",
			Router:           "/",
			AllowHTTPMethods: []string{"post", "get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
