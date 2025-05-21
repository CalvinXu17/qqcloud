package admin

import (
	"encoding/json"
	"github.com/beego/beego/v2/server/web"
	"strconv"
	"strings"
	"time"
	"verification/controllers/common"
	"verification/models"
)

type BaseController struct {
	web.Controller
	apiUrl       string
	ManagerIdArr []int
	ManagerId    int
	ManagerInfo  models.Manager
}

type SuccessResult struct {
	Errno  int         `json:"errno"`
	Data   interface{} `json:"data"`
	Errmsg string      `json:"errmsg"`
}

type ErrorResult struct {
	Errno  int    `json:"errno"`
	Errmsg string `json:"errmsg"`
}

// 判断是否开发者
func (b *BaseController) IsDeveloper() {
	if b.ManagerInfo.Pid > 0 {
		b.Error(400, "暂无权限")
	}
}

// 判断是否操作非自身的代理
func (b *BaseController) IsInMangaerList(id int) bool {
	for _, i := range b.ManagerIdArr {
		if id == i {
			return true
		}
	}
	b.Error(400, "无法操作不属于自己的代理")
	return false
}

func (b *BaseController) Success(code int, data any, msg string) {
	res := SuccessResult{Errno: code, Data: data, Errmsg: msg}
	b.Data["json"] = res
	_ = b.ServeJSON()
	b.StopRun()
}

func (b *BaseController) Error(code int, msg string) {
	res := ErrorResult{Errno: code, Errmsg: msg}
	b.Data["json"] = res
	_ = b.ServeJSON()
	b.StopRun()
}

func (b *BaseController) Prepare() {
	var managerIdList []int
	b.apiUrl = "http://www.developer.qqcloudcom.cn"
	Uri := b.Ctx.Request.RequestURI
	// 只有登录路由才不需要检测是否登录
	if strings.ToLower(Uri) != "/admin/user/login" {
		token := b.Ctx.Request.Header["Token"]
		status, ac := common.GetCacheAC()
		if status == false {
			b.Error(400, "缓存初始化失败")
		}
		id := b.GetSession("token")
		Id := common.GetInterfaceToInt(id)
		if id == "" || Id <= 0 {
			b.Error(401, "请重新登录")
		}
		_ = ac.Put(token[0], strconv.Itoa(Id), 2*60*60*time.Second)
		m := models.Manager{ID: Id}
		b.ManagerId = Id
		m.GetManagerIdList(Id, &managerIdList)
		b.ManagerIdArr = managerIdList
		// 获取缓存的用户信息
		userCache := ac.Get("manager-info-" + strconv.Itoa(Id))
		var user = models.Manager{}
		if userCache == "" {
			b.Error(401, "请重新登录")
		}
		_ = json.Unmarshal([]byte(common.Strval(userCache)), &user)
		b.ManagerInfo = user
		// 如果不是开发者则需要检测权限
		if b.ManagerInfo.Pid > 0 {
			r := models.RoleItem{RoleId: b.ManagerInfo.PowerId}
			// 获取预制权限列表
			list := r.GetUserRole()
			roleList := models.GetAllRole()
			var needCheck = false
			errMsg := ""
			// 匹配是否需要权限的路由
			for _, i := range roleList {
				if strings.ToLower(i.Path) == strings.ToLower(Uri) {
					needCheck = true
					errMsg = i.Name
				}
			}
			// 如果需要检查权限则开始检查权限
			if needCheck == true {
				var hasPermission = false
				for _, i := range list {
					if strings.ToLower(i.Path) == strings.ToLower(Uri) {
						if i.Index != "" && i.Value != "" {
							value := b.GetString(i.Index)
							if strings.Contains(value,i.Value) {
								hasPermission = true
							}
						} else {
							hasPermission = true
						}
					}
				}
				if hasPermission == false {
					if errMsg != "" {
						b.Error(400, "您暂无["+errMsg+"]的操作权限")
					}
					b.Error(400, "暂无权限")
				}
			}
		}
	}

}
