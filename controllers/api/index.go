package api

const (
	ApplicationJSON = "application/json"
)

type IndexController struct {
	BaseController
}

// Index @Title
// @router / [post,get]
func (p *IndexController) Index() {
	action := p.Param.Action
	switch action {
	case "soft.timestamp":
		p.GetTime()
	case "soft.init":
		p.GetSoftInfo()
	case "user.register":
		p.Register()
	case "user.login":
		p.Login()
	case "user.bind":
		p.UnBind()
	case "user.forget":
		p.Forget()
	case "user.recharge":
		p.Recharge()
	case "user.points":
		p.Points()
	case "user.heart":
		p.Heart()
	case "user.client":
		p.Client()
	case "user.logout":
		p.Logout()
	}
}
