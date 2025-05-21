package admin

import (
	"encoding/json"
	"strconv"
	"time"
	"verification/controllers/common"
	"verification/models"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type UserController struct {
	BaseController
}

// Login @Title 登录
// @router /login [post]
func (u *UserController) Login() {
	logs.Info(u.Ctx.Request)
	username := u.GetString("user")
	password := u.GetString("pwd")
	var user *models.Manager
	err, user := user.Login(username, password)
	if err != "" {
		u.Error(400, err)
	}
	token := common.GetToken()
	user.Pwd = ""
	result := common.LoginResult{Info: user, Token: token}
	status, cacheClient := common.GetCacheAC()
	if status == false {
		logs.Error("缓存初始化错误")
		return
	}
	_ = u.SetSession("token", strconv.Itoa(user.ID))
	_ = cacheClient.Put(token, strconv.Itoa(user.ID), 2*60*60*time.Second)
	data, _ := json.Marshal(user)
	_ = cacheClient.Put("manager-info-"+strconv.Itoa(user.ID), string(data), 365*60*60*time.Second)
	u.Success(0, result, "登录成功")
}

type versionData struct {
	Version    string `json:"version"`
	Update     bool   `json:"update"`
	NewVersion string `json:"new_version"`
	Notice     string `json:"notice"`
}

type VersionRes struct {
	Errno  int         `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   versionData `json:"data"`
}

// Info @Title 获取登录用户信息
// @router /info [post]
func (u *UserController) Info() {
	var update bool
	var newVersion string
	var notice string

	managerId := common.GetManagerId(u.GetSession("token"))
	if managerId == 0 {
		u.Error(400, "请先登录")
	}
	var m *models.Manager
	user, err := m.GetInfoById(managerId)
	if err != nil {
		u.Error(400, "获取失败")
	}
	user.Pwd = ""

	var Permissions []models.Permission
	if user.Pid == 0 {
		Permissions = models.DeveloperPermission
	} else {
		Permissions = models.AgentPermission
	}
	version := "1.11"
	versionRes := VersionRes{}
	err = json.Unmarshal([]byte(""), &versionRes)
	if err != nil {
		notice = ""
		update = false
		newVersion = "1.00"
	} else {
		notice = versionRes.Data.Notice
		update = versionRes.Data.Update
		newVersion = versionRes.Data.NewVersion
	}
	info := models.ManagerInfo{
		UserId:      managerId,
		Username:    user.User,
		RealName:    user.User,
		Avatar:      user.Avatar,
		Desc:        "",
		Token:       common.GetTokenString(u.Ctx),
		Permissions: Permissions,
		User:        user,
		Version:     version,
		Update:      update,
		NewVersion:  newVersion,
		Notice:      notice,
	}
	u.Success(0, info, "获取成功")
}

// Logout @Title 退出
// @router /logout [post]
func (u *UserController) Logout() {
	managerId := common.GetManagerId(u.GetSession("token"))
	if managerId > 0 {
		_ = u.SetSession("token", "0")
	}
	u.Success(0, nil, "退出成功")
}

type ManagerStatic struct {
	Project int64 `json:"project"`
	Card    int64 `json:"card"`
	Member  int64 `json:"member"`
}

// GetInfo @Title 获取项目统计信息
// @router /getInfo [post]
func (u *UserController) GetInfo() {
	var c ManagerStatic
	managerId := common.GetManagerId(u.GetSession("token"))
	if u.ManagerInfo.Pid == 0 {
		p := models.Project{ManagerId: managerId}
		k := models.Keys{ManagerId: managerId}
		m := models.Member{ManagerId: managerId}
		c = ManagerStatic{
			Project: p.GetCount(),
			Card:    k.GetCount(u.ManagerIdArr),
			Member:  m.GetCount(u.ManagerIdArr),
		}
	} else {
		m := models.ManagerCards{}
		project := 0
		status, list := m.GetGroupCardList(u.ManagerInfo.ID)
		logs.Error("list：", list)
		if status == false {
			project = 0
		} else {
			project = len(list)
		}
		s := models.Member{ManagerId: managerId}
		k := models.Keys{ManagerId: managerId}
		c = ManagerStatic{
			Project: int64(project),
			Card:    k.GetCount(u.ManagerIdArr),
			Member:  s.GetCount(u.ManagerIdArr),
		}
	}

	u.Success(0, c, "获取成功")
}

type Notice struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CreateTime time.Time `json:"create_time"`
}

type NoticeRes struct {
	Errno  int      `json:"errno"`
	Errmsg string   `json:"errmsg"`
	Data   []Notice `json:"data"`
}

// GetSysNotice @Title 获取公告
// @router /getSysNotice [post]
func (u *UserController) GetSysNotice() {
	var data NoticeRes
	_ = json.Unmarshal([]byte(""), &data)
	u.Success(0, data.Data, "获取成功")
}

type EchartsRes struct {
	Login    []int64  `json:"login"`
	Register []int64  `json:"register"`
	Range    []string `json:"range"`
}

// GetMemberEcharts @Title 获取统计图表数据
// @router /getMemberEcharts [post]
// 登录 注册
func (u *UserController) GetMemberEcharts() {
	projectId, _ := u.GetInt("project_id", 0)
	var res EchartsRes
	dateRange := common.GetDateRange(7)
	m := models.Member{}
	res.Register = m.GetRangeCount(dateRange, projectId, u.ManagerIdArr)

	l := models.MemberLogin{}
	res.Login = l.GetRangeCount(dateRange, projectId, u.ManagerIdArr)
	var date []string
	for _, i := range dateRange {
		date = append(date, i.Date)
	}
	res.Range = date
	u.Success(0, res, "获取成功")
}

type EchartsKeysRes struct {
	Active []int64  `json:"active"`
	Add    []int64  `json:"add"`
	Range  []string `json:"range"`
}

// GetKeysEcharts @Title 获取统计图表数据
// @router /getKeysEcharts [post]
// 创建 激活
func (u *UserController) GetKeysEcharts() {
	projectId, _ := u.GetInt("project_id", 0)
	var res EchartsKeysRes
	dateRange := common.GetDateRange(7)
	m := models.Keys{}
	res.Add = m.GetAddRangeCount(dateRange, projectId, u.ManagerIdArr)
	res.Active = m.GetActiveRangeCount(dateRange, projectId, u.ManagerIdArr)
	var date []string
	for _, i := range dateRange {
		date = append(date, i.Date)
	}
	res.Range = date
	u.Success(0, res, "获取成功")
}

// Update @Title 更新信息
// @router /update [post]
func (u *UserController) Update() {
	runmode, _ := beego.AppConfig.String("runmode")
	if runmode == "dev" {
		u.Error(400, "调试模式下不允许修改超管账号")
	}
	managerId := common.GetManagerId(u.GetSession("token"))
	username := u.GetString("user", "")
	pwd := u.GetString("pwd", "")
	email := u.GetString("email", "")
	user := models.Manager{ID: managerId}
	n := models.Manager{
		User:  username,
		Pwd:   pwd,
		Email: email,
	}
	row, msg := user.UpdateById(n)
	if row > 0 {
		u.Success(0, nil, msg)
	}
	u.Error(400, msg)
}
