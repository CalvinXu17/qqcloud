package api

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
	"verification/controllers/common"
	"verification/models"
	"verification/validation/api"

	"github.com/astaxie/beego/validation"
	"github.com/beego/beego/v2/core/logs"
)

type Func struct {
	IndexController
}

func (p *IndexController) GetTime() {
	timeUnix := time.Now().Unix()
	p.CallJson("成功", timeUnix)
}

type SoftInfo struct {
	AppName             string `json:"appName"`
	AppStatus           string `json:"appStatus"`
	AppMode             string `json:"appMode"`
	AppNotice           string `json:"appNotice"`
	Version             string `json:"version"`
	VersionNotice       string `json:"versionNotice"`
	VersionIsMustUpdate int    `json:"VersionIsMustUpdate"`
	AppDownUrl          string `json:"appDownUrl"`
	VersionUp           string `json:"versionUp"`
	VersionUpNotice     string `json:"versionUpNotice"`
	LoginRules          string `json:"loginRules"`
}

func (p *IndexController) GetSoftInfo() {
	var appStatus = []string{0: "收费运营", 1: "停止运营", 2: "免费运营"}
	var appMode = []string{0: "单码", 1: "账号"}
	var mustUpdate = 0
	if p.ProjectVersion.IsMustUpdate == 0 {
		mustUpdate = 1
	}
	if p.ProjectVersion.IsActive > 0 {
		mustUpdate = 1
	}
	s := SoftInfo{
		AppName:             p.Project.Name,
		AppStatus:           appStatus[p.Project.StatusType],
		AppMode:             appMode[p.Project.Type],
		AppNotice:           p.Project.Notice,
		Version:             common.GetVersionString(p.ProjectVersion.Version),
		VersionNotice:       p.ProjectVersion.Notice,
		VersionIsMustUpdate: mustUpdate,
		AppDownUrl:          p.ProjectUpVersion.WgtUrl,
		VersionUp:           common.GetVersionString(p.ProjectUpVersion.Version),
		VersionUpNotice:     p.ProjectUpVersion.Notice,
		LoginRules:          p.ProjectLogin.Title,
	}
	p.CallJson("获取成功", s)
}

func (p *IndexController) Register() {
	valid := validation.Validation{}
	if p.Project.Type == 0 {
		p.CallErrorJson("单码状态下无需注册账号", nil)
	}
	param := api.RegisterParam{}
	// 如果没有加密则从json 表单 url参数中提取到param结构体
	if p.IsEncrypt == false {
		if p.IsPost == true && p.ParamContentType == ApplicationJSON {
			_ = p.Ctx.BindJSON(&param)
		} else {
			param.User = p.GetString("user", "")
			param.Pwd = p.GetString("pwd", "")
			param.Pwd2 = p.GetString("pwd2", "")
			param.Email = p.GetString("email", "")
			param.Key = p.GetString("key", "")
			param.Recommender = p.GetString("recommender", "")
			param.Code = p.GetString("code", "")
			param.Captcha = p.GetString("captcha", "")
		}
	}
	// 如果加密了则从解密的json中提取到param结构体
	if p.IsEncrypt == true {
		_ = json.Unmarshal([]byte(p.decryptJsonString), &param)
	}
	b, err := valid.Valid(&param)
	if err != nil {
		logs.Error("验证器错误", err)
	}
	if !b {
		for _, err := range valid.Errors {
			logs.Error(err.Key, err.Message)
			p.Error(9999, err.Message)
		}
	}

	if p.Project.Type == 1 {
		p.CallErrorJson("已停止运营", nil)
	}
	member := models.Member{
		ManagerId:     p.Project.ManagerId,
		ProjectId:     p.Project.ID,
		Email:         param.Email,
		Name:          param.User,
		NickName:      param.User,
		Password:      param.Pwd,
		SafePassword:  param.Pwd2,
		ActiveTime:    time.Now().Unix(),
		LastLoginTime: time.Now().Unix(),
		LastLoginIp:   p.Ip,
	}
	k := models.Keys{LongKeys: param.Key}
	if p.ProjectLogin.RegMode == 0 && param.Key == "" {
		p.CallErrorJson("请填写激活码进行注册", nil)
	}
	// 判断充值卡有没有填写
	if param.Key != "" {
		status, cards := k.GetOneByKeys()
		if status == false {
			p.CallErrorJson("激活码无效", status)
		}
		if k.ProjectId != 0 && k.ProjectId != p.Project.ID {
			p.CallErrorJson("激活码无效", status)
		}
		if k.IsLock > 0 {
			p.CallErrorJson("激活码已锁定", status)
		}
		if k.UseTime > 0 {
			p.CallErrorJson("激活码已被使用", status)
		}
		tag := ""
		if k.Tag == "" {
			tag = cards.Tag
		}
		member.Days = cards.Days
		member.Points = cards.Points
		member.Tag = tag
		member.KeyExtAttr = cards.KeyExtAttr
		member.EndTime = time.Now().Unix() + int64(cards.Days*24*60*60)
		member.ManagerId = k.ManagerId
	}
	if p.Project.StatusType == 2 {
		member.Days = 999
		member.Points = 999
		member.Tag = "免费用户"
		member.KeyExtAttr = "免费用户"
		member.ManagerId = p.Project.ManagerId
		member.EndTime = time.Now().Unix() + int64(24*60*60)
	}
	status, msg := member.Register(k)
	if status == true {
		p.CallJson(msg, status)
	}
	p.CallErrorJson(msg, status)
}

type MemberInfo struct {
	Client             string  `json:"client"`
	Endtime            string  `json:"endtime"`
	EndtimeTimestamp   int64   `json:"endtimeTimestamp"`
	Starttime          string  `json:"starttime"`
	StarttimeTimestamp int64   `json:"StarttimeTimestamp"`
	RealCdays          float64 `json:"realCdays"`
	CountCdays         float64 `json:"countCdays"`
	CountPoints        int     `json:"countPoints"`
	Username           string  `json:"username"`
	Nickname           string  `json:"nickname"`
	Ip                 string  `json:"ip"`
	Mac                string  `json:"mac"`
	KeyExtAttr         string  `json:"keyExtAttr"`
	Tag                string  `json:"tag"`
}

type HeartList struct {
	List     []string         `json:"list"`
	Online   []api.OnlineData `json:"online"`
	PcOnline []api.OnlineData `json:"pcOnline"`
}

func (p *IndexController) Login() {
	var status bool
	var msg string
	client := common.GetToken()
	param := api.LoginParam{}
	// 如果没有加密则从json 表单 url参数中提取到param结构体
	if p.IsEncrypt == false {
		if p.IsPost == true && p.ParamContentType == ApplicationJSON {
			_ = p.Ctx.BindJSON(&param)
		} else {
			param.User = p.GetString("user", "")
			param.Pwd = p.GetString("pwd", "")
		}
	}
	// 如果加密了则从解密的json中提取到param结构体
	if p.IsEncrypt == true {
		_ = json.Unmarshal([]byte(p.decryptJsonString), &param)
	}
	isEmail := common.VerifyEmailFormat(param.User)
	u := models.Member{}
	if p.Project.StatusType == 2 {
		u.Name = p.Param.Mac
		param.User = p.Param.Mac
	}
	if p.Project.Type == 0 {
		u.Name = param.User
		param.Pwd = ""
		// 检查激活码用户是否存在
		status, msg = u.CheckMember()
		if status == false {
			// 如果不存在则调用激活码注册会员
			status, msg = u.KeyRegister(param.User, p.Project, p.Param, p.Ip)
			if status == false {
				p.CallErrorJson(msg, nil)
			}
		}
	} else {
		u.Password = param.Pwd
		if isEmail == true {
			u.Email = param.User
		} else {
			u.Name = param.User
		}

	}
	// 检查用户是否存在
	status, msg = u.CheckMember()
	if status == false && p.Project.Type == 1 {
		p.CallErrorJson(msg, nil)
	}
	if p.Project.ID != u.ProjectId {
		p.CallErrorJson("项目不匹配，请重新选择", nil)
	}
	if p.Project.StatusType == 0 {
		if u.EndTime < time.Now().Unix() {
			p.CallErrorJson("会员账号已过期", nil)
		}
	}
	if p.Project.StatusType != 2 && u.Tag == "免费用户" && p.Project.Type == 0 {
		p.CallErrorJson("免费用户已不可用，请使用激活码进行登录", nil)
	}
	if p.Project.StatusType != 2 && u.Tag == "免费用户" && p.Project.Type == 1 {
		p.CallErrorJson("免费用户已不可用，请注册会员之后再登录", nil)
	}
	if u.IsLock > 0 {
		p.CallErrorJson("会员账号已被锁定", nil)
	}

	heartList := HeartList{}
	heartList = p.FetchHeartOnline(u)
	switch p.ProjectLogin.Mode {
	case 0:
		// 绑定登录
		// 自动解绑
		if u.Mac == "" {
			u.Mac = p.Param.Mac
		}
		if p.ProjectLogin.UnbindMode == 2 && p.Param.Mac != u.Mac {
			// 执行解绑
			status, msg = u.Unbind(p.ProjectLogin, p.Param.Mac)
			if status == false {
				p.CallErrorJson(msg, nil)
			}
			// 添加解绑记录
			_ = models.AddLog(p.Ip, int(time.Now().Unix()), p.Param.Mac, 1, client, u)

		}
		// 如果不是自动解绑，则提示需要解绑
		if p.Param.Mac != u.Mac && p.ProjectLogin.UnbindMode != 2 {
			p.CallErrorJson("请先解绑，再进行登录", nil)
		}
		// 如果设置允许顶号登录，则先下线之前所有的在线
		if p.ProjectLogin.PcMore == 0 {
			p.RemoveHeartAll(u, "您的会员账号当前已在其它机器登录")
		}
		// 判断全局单机多开
		if p.ProjectLogin.PcMore > 0 {
			if len(heartList.PcOnline) > p.ProjectLogin.PcMore {
				p.CallErrorJson("单机在线已超出限制："+strconv.Itoa(p.ProjectLogin.PcMore), p.ProjectLogin.PcMore)
			}
		}

	case 1:
		// 普通登录
		// 如果不允许单机多开 则下线单机的所有在线
		if p.ProjectLogin.PcMore == 0 && p.ProjectLogin.PcCodeMore > 0 {
			p.RemoveHeartAllByMac(u, p.Param.Mac, "您的会员账号当前已在其它机器登录")
		}
		// 如果不允许多机器多开 则下线所有在线
		if p.ProjectLogin.PcMore == 0 && p.ProjectLogin.PcCodeMore == 0 {
			p.RemoveHeartAll(u, "您的会员账号当前已在其它机器登录")
		}

	case 2:
		// 点数登录
		// 判断全局单机多开
		if p.ProjectLogin.PcMore > 0 {
			if len(heartList.PcOnline) > p.ProjectLogin.PcMore {
				p.CallErrorJson("单机在线已超出限制："+strconv.Itoa(p.ProjectLogin.PcMore), p.ProjectLogin.PcMore)
			}
		}
		// 判断点数全局多开
		if p.ProjectLogin.NumberMore > 0 {
			if len(heartList.Online) > p.ProjectLogin.NumberMore {
				p.CallErrorJson("在线已超出限制："+strconv.Itoa(p.ProjectLogin.NumberMore), p.ProjectLogin.NumberMore)
			}
		}
	}
	u.LastLoginIp = p.Ip
	u.LastLoginTime = time.Now().Unix()
	status, msg = u.UpdateInfo()
	//if status == false {
	//	p.CallErrorJson(msg, nil)
	//}
	// 将用户基本信息缓存起来，方便其它查询接口调用
	memberJson, _ := json.Marshal(u)
	_ = p.Ac.Put("member-key-"+strconv.Itoa(p.Project.ID)+"-"+u.Name, string(memberJson), 20*60*60*time.Second)
	_ = p.Ac.Put("member-key-"+strconv.Itoa(p.Project.ID)+"-"+u.Email, string(memberJson), 20*60*60*time.Second)
	// 写入登录日志
	status = models.AddLog(p.Ip, int(time.Now().Unix()), p.Param.Mac, 0, client, u)
	status = models.AddLog(p.Ip, int(time.Now().Unix()), p.Param.Mac, 3, client, u)
	if status == false {
		p.CallErrorJson("请稍后再尝试登录", nil)
	}
	// 写入心跳在线缓存
	o := api.OnlineData{
		Mac:        p.Param.Mac,
		Addtime:    int(time.Now().Unix()),
		ClientTime: int(time.Now().Unix()),
		Ip:         p.Ip,
		Email:      u.Email,
		Name:       u.Name,
		NickName:   u.NickName,
		MemberId:   u.ID,
		ManagerId:  u.ManagerId,
		ProjectId:  u.ProjectId,
		Client:     client,
	}
	p.InsertHeartList(u, o)

	realCday := float64((u.EndTime - time.Now().Unix()) / 24 * 60 * 60)

	i := MemberInfo{
		Client:             client,
		Endtime:            time.Unix(int64(u.EndTime), 0).Format("2006-01-02 15:04:05"),
		EndtimeTimestamp:   u.EndTime,
		Starttime:          time.Unix(int64(u.ActiveTime), 0).Format("2006-01-02 15:04:05"),
		StarttimeTimestamp: u.ActiveTime,
		RealCdays:          realCday,
		CountCdays:         u.Days,
		CountPoints:        u.Points,
		Username:           u.Name,
		Nickname:           u.NickName,
		Ip:                 p.Ip,
		Mac:                p.Param.Mac,
		KeyExtAttr:         u.KeyExtAttr,
		Tag:                u.Tag,
	}
	p.CallJson("登录成功", i)
}

func (p *IndexController) UnBind() {
	var status bool
	var msg string
	if p.ProjectLogin.Mode != 0 {
		p.CallErrorJson("无需解绑", nil)
	}
	if p.ProjectLogin.UnbindMode == 0 {
		p.CallErrorJson("不允许解绑", nil)
	}
	if p.ProjectLogin.UnbindMode == 2 {
		p.CallErrorJson("已应用自动解绑，无需手动解绑", nil)
	}
	param := api.LoginParam{}
	// 如果没有加密则从json 表单 url参数中提取到param结构体
	if p.IsEncrypt == false {
		if p.IsPost == true && p.ParamContentType == ApplicationJSON {
			_ = p.Ctx.BindJSON(&param)
		} else {
			param.User = p.GetString("user", "")
			param.Pwd = p.GetString("pwd", "")
		}
	}
	// 如果加密了则从解密的json中提取到param结构体
	if p.IsEncrypt == true {
		_ = json.Unmarshal([]byte(p.decryptJsonString), &param)
	}
	var k strings.Builder
	u := models.Member{}
	k.WriteString("member-key-" + strconv.Itoa(p.Project.ID) + "-")
	k.WriteString(param.User)
	memberJson := common.Strval(p.Ac.Get(k.String()))
	if memberJson != "" {
		err := json.Unmarshal([]byte(common.Strval(memberJson)), &u)
		if err == nil {
			if p.Project.Type == 1 && u.Password != param.Pwd {
				p.CallErrorJson("密码错误", nil)
			}
		}
	}
	isEmail := common.VerifyEmailFormat(param.User)
	if memberJson == "" {
		if p.Project.Type == 0 {
			u.Name = param.User

		} else {
			u.Password = param.Pwd
			if isEmail == true {
				u.Email = param.User
			} else {
				u.Name = param.User
			}
		}
		// 检查用户是否存在
		status, msg = u.CheckMember()
		if status == false && p.Project.Type == 1 {
			p.CallErrorJson(msg, nil)
		}
	}
	if u.ID == 0 {
		p.CallErrorJson("会员不存在", nil)
	}
	if u.IsLock > 0 {
		p.CallErrorJson("会员账号已被锁定", nil)
	}
	if u.Mac == "" {
		p.CallErrorJson("您已经解绑过了", nil)
	}
	if u.Mac == p.Param.Mac {
		p.CallErrorJson("无需解绑", nil)
	}
	logs.Error("提取的用户数据", u.Mac, p.Param.Mac)
	// 执行解绑
	status, msg = u.Unbind(p.ProjectLogin, p.Param.Mac)
	if status == false {
		p.CallErrorJson(msg, nil)
	}
	// 添加解绑记录
	_ = models.AddLog(p.Ip, int(time.Now().Unix()), p.Param.Mac, 1, "", u)
	memberJ, _ := json.Marshal(u)
	_ = p.Ac.Put("member-key-"+strconv.Itoa(p.Project.ID)+"-"+u.Name, string(memberJ), 20*60*60*time.Second)
	_ = p.Ac.Put("member-key-"+strconv.Itoa(p.Project.ID)+"-"+u.Email, string(memberJ), 20*60*60*time.Second)
	p.CallJson("解绑成功", nil)
}

func (p *IndexController) Forget() {
	if p.Project.Type == 0 {
		p.CallErrorJson("单码模式下不需要找回账号", nil)
	}
	var status bool
	var msg string
	param := api.ForgetParam{}
	// 如果没有加密则从json 表单 url参数中提取到param结构体
	if p.IsEncrypt == false {
		if p.IsPost == true && p.ParamContentType == ApplicationJSON {
			_ = p.Ctx.BindJSON(&param)
		} else {
			param.User = p.GetString("user", "")
			param.Pwd = p.GetString("pwd", "")
			param.Pwd2 = p.GetString("pwd2", "")
			param.Code = p.GetString("code", "")
			param.Captcha = p.GetString("captcha", "")
		}
	}
	// 如果加密了则从解密的json中提取到param结构体
	if p.IsEncrypt == true {
		_ = json.Unmarshal([]byte(p.decryptJsonString), &param)
	}

	// 检查会员是否合法
	var k strings.Builder
	u := models.Member{}
	k.WriteString("member-key-" + strconv.Itoa(p.Project.ID) + "-")
	k.WriteString(param.User)
	memberJson := common.Strval(p.Ac.Get(k.String()))
	if memberJson != "" {
		err := json.Unmarshal([]byte(common.Strval(memberJson)), &u)
		if err == nil {
			if u.SafePassword != param.Pwd2 {
				p.CallErrorJson("安全密码错误", nil)
			}
		}
	}
	isEmail := common.VerifyEmailFormat(param.User)
	if memberJson == "" {
		u.SafePassword = param.Pwd2
		if isEmail == true {
			u.Email = param.User
		} else {
			u.Name = param.User
		}
		// 检查用户是否存在
		status, msg = u.ForgetCheckMember()
		if status == false && p.Project.Type == 1 {
			p.CallErrorJson(msg, nil)
		}
	}
	if u.ID == 0 {
		p.CallErrorJson("会员不存在", nil)
	}
	if u.IsLock > 0 {
		p.CallErrorJson("会员账号已被锁定", nil)
	}
	u.Password = param.Pwd
	status, msg = u.UpdateInfo()
	if status == true {
		memberJ, _ := json.Marshal(u)
		_ = p.Ac.Put("member-key-"+strconv.Itoa(p.Project.ID)+"-"+u.Name, string(memberJ), 20*60*60*time.Second)
		_ = p.Ac.Put("member-key-"+strconv.Itoa(p.Project.ID)+"-"+u.Email, string(memberJ), 20*60*60*time.Second)
		p.CallJson(msg, nil)
	}
	p.CallErrorJson(msg, nil)
}

type RechargeResult struct {
	AddCdays    float64 `json:"addCdays"`
	AddPoints   int     `json:"addPoints"`
	KeyCdays    float64 `json:"keyCdays"`
	KeyPoints   int     `json:"keyPoints"`
	CountCdays  float64 `json:"countCdays"`
	CountPoints int     `json:"countPoints"`
}

func (p *IndexController) Recharge() {
	var status bool
	var msg string
	if p.Project.StatusType == 2 {
		p.CallErrorJson("免费运营不需要充值", nil)
	}
	param := api.RechargeParam{}
	// 如果没有加密则从json 表单 url参数中提取到param结构体
	if p.IsEncrypt == false {
		if p.IsPost == true && p.ParamContentType == ApplicationJSON {
			_ = p.Ctx.BindJSON(&param)
		} else {
			param.User = p.GetString("user", "")
			param.Key = p.GetString("key", "")
		}
	}
	// 如果加密了则从解密的json中提取到param结构体
	if p.IsEncrypt == true {
		_ = json.Unmarshal([]byte(p.decryptJsonString), &param)
	}
	k := models.Keys{}
	status, msg, rechargeResult := k.Recharge(param.User, param.Key, p.Project)
	if status == false {
		p.CallErrorJson(msg, nil)
	}
	_ = p.Ac.Delete("member-key-" + strconv.Itoa(p.Project.ID) + "-" + param.User)
	p.CallJson(msg, rechargeResult)
}

type PointsResult struct {
	CountPoints int `json:"countPoints"`
}

func (p *IndexController) Points() {
	var status bool
	var msg string
	param := api.PointsParam{}
	// 如果没有加密则从json 表单 url参数中提取到param结构体
	if p.IsEncrypt == false {
		if p.IsPost == true && p.ParamContentType == ApplicationJSON {
			_ = p.Ctx.BindJSON(&param)
		} else {
			param.User = p.GetString("user", "")
			param.Pwd = p.GetString("pwd", "")
			param.Number, _ = p.GetInt("number", 0)
		}
	}
	// 如果加密了则从解密的json中提取到param结构体
	if p.IsEncrypt == true {
		_ = json.Unmarshal([]byte(p.decryptJsonString), &param)
	}
	if param.Number <= 0 {
		p.CallErrorJson("扣点必须大于0", nil)
	}
	var k strings.Builder
	u := models.Member{}
	k.WriteString("member-key-" + strconv.Itoa(p.Project.ID) + "-")
	k.WriteString(param.User)
	memberJson := common.Strval(p.Ac.Get(k.String()))
	if memberJson != "" {
		err := json.Unmarshal([]byte(common.Strval(memberJson)), &u)
		if err == nil {
			if p.Project.Type == 1 && u.Password != param.Pwd {
				p.CallErrorJson("密码错误", nil)
			}
		}
	}
	isEmail := common.VerifyEmailFormat(param.User)
	if memberJson == "" {
		if p.Project.Type == 0 {
			u.Name = param.User

		} else {
			u.Password = param.Pwd
			if isEmail == true {
				u.Email = param.User
			} else {
				u.Name = param.User
			}
		}
		// 检查用户是否存在
		status, msg = u.CheckMember()
		if status == false && p.Project.Type == 1 {
			p.CallErrorJson(msg, nil)
		}
	}
	if u.ID == 0 {
		p.CallErrorJson("会员不存在", nil)
	}
	if u.IsLock > 0 {
		p.CallErrorJson("会员账号已被锁定", nil)
	}
	if u.Points <= 0 {
		p.CallErrorJson("剩余点数不足以扣除", nil)
	}
	if param.Number > u.Points {
		p.CallErrorJson("剩余点数不足以扣除", nil)
	}
	u.Points = u.Points - param.Number
	status, msg = u.UpdateInfo()
	if status == false {
		p.CallErrorJson(msg, nil)
	}
	memberJ, _ := json.Marshal(u)
	_ = p.Ac.Put("member-key-"+strconv.Itoa(p.Project.ID)+"-"+u.Name, string(memberJ), 20*60*60*time.Second)
	_ = p.Ac.Put("member-key-"+strconv.Itoa(p.Project.ID)+"-"+u.Email, string(memberJ), 20*60*60*time.Second)
	p.CallJson("扣除成功", PointsResult{CountPoints: u.Points})

}

type HeartInfo struct {
	Client                      string           `json:"client"`
	Endtime                     string           `json:"endtime"`
	EndtimeTimestamp            int64            `json:"endtimeTimestamp"`
	Starttime                   string           `json:"starttime"`
	StarttimeTimestamp          int64            `json:"StarttimeTimestamp"`
	RealCdays                   float64          `json:"realCdays"`
	CountCdays                  float64          `json:"countCdays"`
	CountPoints                 int              `json:"countPoints"`
	Username                    string           `json:"username"`
	Nickname                    string           `json:"nickname"`
	Ip                          string           `json:"ip"`
	Mac                         string           `json:"mac"`
	KeyExtAttr                  string           `json:"keyExtAttr"`
	Tag                         string           `json:"tag"`
	OnlineAllList               []api.OnlineData `json:"onlineAllList"`
	StandAloneAllList           []api.OnlineData `json:"standAloneAllList"`
	StandAloneCurrentLogin      int              `json:"standAloneCurrentLogin"`
	StandAloneAllowCurrentLogin int              `json:"standAloneAllowCurrentLogin"`
	OnlineCurrentLogin          int              `json:"onlineCurrentLogin"`
	OnlineAllowCurrentLogin     int              `json:"onlineAllowCurrentLogin"`
}

func (p *IndexController) Heart() {
	var msg string
	var status bool
	param := api.HeartParm{}
	// 如果没有加密则从json 表单 url参数中提取到param结构体
	if p.IsEncrypt == false {
		if p.IsPost == true && p.ParamContentType == ApplicationJSON {
			_ = p.Ctx.BindJSON(&param)
		} else {
			param.Client = p.GetString("client", "")
		}
	}
	// 如果加密了则从解密的json中提取到param结构体
	if p.IsEncrypt == true {
		_ = json.Unmarshal([]byte(p.decryptJsonString), &param)
	}
	outMsg := common.Strval(p.Ac.Get("h-o-" + param.Client))
	if outMsg != "" {
		msg = common.Strval(outMsg)
		_ = p.Ac.Delete("h-" + param.Client)
		p.CallErrorJson(msg, nil)
	}
	var s strings.Builder
	s.WriteString("h-")
	s.WriteString(param.Client)
	clientJson := common.Strval(p.Ac.Get(s.String()))
	if clientJson == "" {
		logs.Error("已下线1")
		p.CallErrorJson("已下线", nil)
	}
	heartData := api.OnlineData{}
	err := json.Unmarshal([]byte(common.Strval(clientJson)), &heartData)
	if err != nil {
		logs.Error("已下线2")
		p.CallErrorJson("已下线", nil)
	}
	if heartData.Mac != p.Param.Mac {
		p.CallErrorJson("机器码出现变动", nil)
	}

	heartData.ClientTime = int(time.Now().Unix())
	if p.ProjectLogin.Mode == 0 && heartData.Mac != p.Param.Mac {
		var out strings.Builder
		out.WriteString("h-o-")
		out.WriteString(param.Client)
		_ = p.Ac.Put(out.String(), msg, 1*60*60*time.Second)
		_ = p.Ac.Delete("h-" + param.Client)
		p.CallErrorJson("机器码出现变动，请重新登录", nil)
	}
	var heartMacList []string
	var h strings.Builder
	h.WriteString(heartPrefix)
	h.WriteString(strconv.Itoa(p.Project.ID))
	h.WriteString("-")
	h.WriteString(strconv.Itoa(heartData.MemberId))
	heartJson := common.Strval(p.Ac.Get(h.String()))
	_ = json.Unmarshal([]byte(common.Strval(heartJson)), &heartMacList)
	data, _ := json.Marshal(heartMacList)
	_ = p.Ac.Put(h.String(), string(data), 2*60*60*time.Second)
	memberJson := common.Strval(p.Ac.Get("member-key-" + strconv.Itoa(p.Project.ID) + "-" + heartData.Name))
	u := models.Member{ID: heartData.MemberId}
	if memberJson == "" {
		status, msg = u.QueryById()
		if status == false {
			_ = p.Ac.Delete("h-" + param.Client)
			p.CallErrorJson(msg, nil)
		}
	} else {
		_ = json.Unmarshal([]byte(common.Strval(memberJson)), &u)
	}
	if p.Project.StatusType == 0 {
		if u.EndTime < time.Now().Unix() {
			p.CallErrorJson("会员账号已过期", nil)
		}
	}
	if u.IsLock > 0 {
		p.CallErrorJson("会员账号已被锁定", nil)
	}
	realCday := float64((u.EndTime - time.Now().Unix()) / 24 * 60 * 60)
	onlineAllowCurrentLogin := 0
	if p.ProjectLogin.Mode == 2 {
		onlineAllowCurrentLogin = p.ProjectLogin.NumberMore
	} else {
		onlineAllowCurrentLogin = p.ProjectLogin.PcCodeMore
	}
	heartList := HeartList{}
	heartList = p.FetchHeartOnline(u)
	heartInfo := HeartInfo{
		Client:                      param.Client,
		Endtime:                     time.Unix(int64(u.EndTime), 0).Format("2006-01-02 15:04:05"),
		EndtimeTimestamp:            u.EndTime,
		Starttime:                   time.Unix(int64(u.ActiveTime), 0).Format("2006-01-02 15:04:05"),
		StarttimeTimestamp:          u.ActiveTime,
		RealCdays:                   realCday,
		CountCdays:                  u.Days,
		CountPoints:                 u.Points,
		Username:                    u.Name,
		Nickname:                    u.NickName,
		Ip:                          p.Ip,
		Mac:                         p.Param.Mac,
		KeyExtAttr:                  u.KeyExtAttr,
		Tag:                         u.Tag,
		OnlineAllList:               heartList.Online,
		StandAloneAllList:           heartList.PcOnline,
		StandAloneCurrentLogin:      len(heartList.PcOnline),
		StandAloneAllowCurrentLogin: p.ProjectLogin.PcMore,
		OnlineCurrentLogin:          len(heartList.Online),
		OnlineAllowCurrentLogin:     onlineAllowCurrentLogin,
	}
	status = models.AddLog(p.Ip, int(time.Now().Unix()), heartData.Mac, 3, param.Client, u)
	logs.Info("写入情况", status)
	_ = p.Ac.Put(s.String(), heartData, 30*60*time.Second)
	p.CallJson("心跳成功", heartInfo)
}

type ClientInfo struct {
	Mac                         string           `json:"mac"`
	OnlineAllList               []api.OnlineData `json:"onlineAllList"`
	StandAloneAllList           []api.OnlineData `json:"standAloneAllList"`
	StandAloneCurrentLogin      int              `json:"standAloneCurrentLogin"`
	StandAloneAllowCurrentLogin int              `json:"standAloneAllowCurrentLogin"`
	OnlineCurrentLogin          int              `json:"onlineCurrentLogin"`
	OnlineAllowCurrentLogin     int              `json:"onlineAllowCurrentLogin"`
}

func (p *IndexController) Client() {
	var status bool
	var msg string
	param := api.ClientParam{}
	// 如果没有加密则从json 表单 url参数中提取到param结构体
	if p.IsEncrypt == false {
		if p.IsPost == true && p.ParamContentType == ApplicationJSON {
			_ = p.Ctx.BindJSON(&param)
		} else {
			param.User = p.GetString("user", "")
			param.Pwd = p.GetString("pwd", "")
		}
	}
	// 如果加密了则从解密的json中提取到param结构体
	if p.IsEncrypt == true {
		_ = json.Unmarshal([]byte(p.decryptJsonString), &param)
	}
	var k strings.Builder
	u := models.Member{}
	k.WriteString("member-key-" + strconv.Itoa(p.Project.ID) + "-")
	k.WriteString(param.User)
	memberJson := common.Strval(p.Ac.Get(k.String()))
	if memberJson != "" {
		err := json.Unmarshal([]byte(common.Strval(memberJson)), &u)
		if err == nil {
			if p.Project.Type == 1 && u.Password != param.Pwd {
				p.CallErrorJson("密码错误", nil)
			}
		}
	}
	isEmail := common.VerifyEmailFormat(param.User)
	if memberJson == "" {
		if p.Project.Type == 0 {
			u.Name = param.User
		} else {
			u.Password = param.Pwd
			if isEmail == true {
				u.Email = param.User
			} else {
				u.Name = param.User
			}
		}
		// 检查用户是否存在
		status, msg = u.CheckMember()
		if status == false && p.Project.Type == 1 {
			p.CallErrorJson(msg, nil)
		}
	}
	if u.ID == 0 {
		p.CallJson("获取成功", ClientInfo{})
	}
	onlineAllowCurrentLogin := 0
	if p.ProjectLogin.Mode == 2 {
		onlineAllowCurrentLogin = p.ProjectLogin.NumberMore
	} else {
		onlineAllowCurrentLogin = p.ProjectLogin.PcCodeMore
	}
	heartList := HeartList{}
	heartList = p.FetchHeartOnline(u)
	clientInfo := ClientInfo{
		Mac:                         p.Param.Mac,
		OnlineAllList:               heartList.Online,
		StandAloneAllList:           heartList.PcOnline,
		StandAloneCurrentLogin:      len(heartList.PcOnline),
		StandAloneAllowCurrentLogin: p.ProjectLogin.PcMore,
		OnlineCurrentLogin:          len(heartList.Online),
		OnlineAllowCurrentLogin:     onlineAllowCurrentLogin,
	}
	p.CallJson("获取成功", clientInfo)
}

func (p *IndexController) Logout() {
	param := api.LogoutParm{}
	// 如果没有加密则从json 表单 url参数中提取到param结构体
	if p.IsEncrypt == false {
		if p.IsPost == true && p.ParamContentType == ApplicationJSON {
			_ = p.Ctx.BindJSON(&param)
		} else {
			param.Client = p.GetString("client", "")
			param.Type, _ = p.GetInt("client", 0)
		}
	}
	// 如果加密了则从解密的json中提取到param结构体
	if p.IsEncrypt == true {
		_ = json.Unmarshal([]byte(p.decryptJsonString), &param)
	}
	var s strings.Builder
	s.WriteString("h-")
	s.WriteString(param.Client)
	clientJson := common.Strval(p.Ac.Get(s.String()))
	if clientJson == "" {
		p.CallErrorJson("已下线", nil)
	}
	heartData := api.OnlineData{}
	err := json.Unmarshal([]byte(common.Strval(clientJson)), &heartData)
	if err != nil {
		p.CallErrorJson("已下线", nil)
	}
	u := models.Member{
		ManagerId: heartData.ManagerId,
		ProjectId: heartData.ProjectId,
		Name:      heartData.Name,
		ID:        heartData.MemberId,
	}
	if param.Type == 0 {
		_ = p.Ac.Put("h-o-"+param.Client, "强制下线", 1*60*60*time.Second)
		_ = p.Ac.Delete("h-" + param.Client)

		// 写入登录日志
		_ = models.AddLog(p.Ip, int(time.Now().Unix()), heartData.Mac, 2, param.Client, u)
		p.CallJson("下线成功", nil)
	} else {
		_ = models.AddLog(p.Ip, int(time.Now().Unix()), heartData.Mac, 2, param.Client, u)
		p.RemoveHeartAll(u, "强制下线")
		p.CallJson("全部下线成功", nil)
	}
}
