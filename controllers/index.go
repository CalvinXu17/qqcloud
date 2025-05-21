package controllers

import (
	"github.com/beego/beego/v2/server/web"
)

type MainController struct {
	web.Controller
}

// Index @Title
// @router / [get]
func (this *MainController) Index() {
	this.TplName = "index.html"
}
