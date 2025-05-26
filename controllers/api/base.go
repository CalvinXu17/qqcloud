package api

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"verification/controllers/common"
	"verification/models"
	"verification/validation/api"

	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/validation"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/dop251/goja"
	uuid "github.com/satori/go.uuid"
	"github.com/yuchenfw/gocrypt"
)

const (
	heartPrefix = "member-heart-"
)

type BaseController struct {
	web.Controller
	Project           models.Project
	IsPost            bool
	EnParam           api.Encrypt
	Param             api.UnEncrypt
	ProjectVersion    models.ProjectVersion
	ProjectUpVersion  models.ProjectVersion
	ProjectLogin      models.ProjectLogin
	ParamContentType  string
	IsEncrypt         bool
	Encrypt           int
	decryptJsonString string
	Ip                string
	Ac                cache.Cache
}

type UnEncryptSuccessResult struct {
	Errno     int         `json:"errno"`
	Data      interface{} `json:"data"`
	Errmsg    string      `json:"errmsg"`
	Uid       string      `json:"uid"`
	Timestamp int64       `json:"timestamp"`
	Sign      string      `json:"sign"`
	Signal    string      `json:"signal"`
}
type UnEncryptErrorResult struct {
	Errno  int    `json:"errno"`
	Errmsg string `json:"errmsg"`
}

type EncryptSuccessResult struct {
	Errno  int    `json:"errno"`
	Data   string `json:"data"`
	Signal string `json:"signal"`
}

func (b *BaseController) Success(code int, data any, msg string, uid string, timestamp int64, sign string, signal string) {
	res := UnEncryptSuccessResult{Errno: code, Data: data, Errmsg: msg, Uid: uid, Timestamp: timestamp, Sign: sign, Signal: signal}
	b.Data["json"] = res
	_ = b.ServeJSON()
	b.StopRun()
}
func (b *BaseController) EncryptSuccess(code int, data string, signal string) {
	res := EncryptSuccessResult{
		Errno:  code,
		Data:   data,
		Signal: signal,
	}
	b.Data["json"] = res
	_ = b.ServeJSON()
	b.StopRun()
}

func (b *BaseController) Error(code int, msg string) {
	res := UnEncryptErrorResult{Errno: code, Errmsg: msg}
	b.Data["json"] = res
	_ = b.ServeJSON()
	b.StopRun()
}

func (b *BaseController) CallJson(msg string, result any) {
	var sign string
	timestamp := time.Now().UnixNano() / 1e6
	uid := uuid.NewV4().String()
	var t strings.Builder
	t.WriteString(b.Project.SecretKey)
	t.WriteString(uid)
	t.WriteString(strconv.FormatInt(timestamp, 10))
	re := md5.Sum([]byte(t.String()))
	sign = fmt.Sprintf("%x", re)
	r := common.RSACrypt{
		PublicKey:  b.Project.PublicKey,
		PrivateKey: b.Project.PrivateKey,
		Type:       gocrypt.Base64,
	}
	hashObj := UnEncryptSuccessResult{Errno: 0, Data: result, Errmsg: msg, Uid: uid, Timestamp: timestamp, Sign: sign, Signal: ""}
	hashStr, _ := json.Marshal(hashObj)
	md5HashStr := common.GetStringMd5(string(hashStr))
	signal := r.RSASign(md5HashStr, b.Project.Sign)
	if b.Encrypt == 1 {
		script := aes
		vm := goja.New()
		_, err := vm.RunString(script)
		if err != nil {
			logs.Error("加密脚本出现问题")
			b.Error(10001, "加密脚本出现问题,请联系开发者")
		}
		var encrypt func(string, string, string) string
		err = vm.ExportTo(vm.Get("AesEncrypt"), &encrypt)
		if err != nil {
			logs.Error("加密脚本映射到Go函数失败")
			b.Error(10002, "加密脚本映射到Go函数失败")
		}
		encryptResultStr := encrypt(string(hashStr), b.Project.KeyA, b.Project.KeyB)
		b.EncryptSuccess(0, encryptResultStr, signal)
	} else {
		b.Success(0, result, msg, uid, timestamp, sign, signal)
	}

}

func (b *BaseController) CallErrorJson(msg string, result any) {
	var sign string
	timestamp := time.Now().UnixNano() / 1e6
	uid := uuid.NewV4().String()
	var t strings.Builder
	t.WriteString(b.Project.SecretKey)
	t.WriteString(uid)
	t.WriteString(strconv.FormatInt(timestamp, 10))
	re := md5.Sum([]byte(t.String()))
	sign = fmt.Sprintf("%x", re)
	r := common.RSACrypt{
		PublicKey:  b.Project.PublicKey,
		PrivateKey: b.Project.PrivateKey,
		Type:       gocrypt.Base64,
	}
	hashObj := UnEncryptSuccessResult{Errno: 400, Data: result, Errmsg: msg, Uid: uid, Timestamp: timestamp, Sign: sign, Signal: ""}
	hashStr, _ := json.Marshal(hashObj)
	md5HashStr := common.GetStringMd5(string(hashStr))
	signal := r.RSASign(md5HashStr, b.Project.Sign)

	b.Success(400, result, msg, uid, timestamp, sign, signal)
}

func (b *BaseController) GetContentType() string {
	ct, exist := b.Ctx.Request.Header["Content-Type"]
	if !exist || len(ct) == 0 {
		return "application/json"
	}
	i, l := 0, len(ct[0])
	for i < l && ct[0][i] != ';' {
		i++
	}
	return ct[0][0:i]
}
func In(index string, array []string) bool {
	for _, e := range array {
		if e == index {
			return true
		}
	}
	return false
}

func (p *BaseController) Prepare() {
	valid := validation.Validation{}
	realIp := p.Ctx.Request.Header.Get("X-Real-ip")
	addr := common.GetAddrIp(p.Ctx.Request.RemoteAddr)
	logs.Info("ip:", realIp, addr)
	if realIp == "" {
		p.Ip = addr
	} else {
		p.Ip = realIp
	}
	// 初始化缓存
	status, ac := common.GetCacheAC()
	p.Ac = ac
	if status == false {
		p.Error(10000, "缓存初始化失败")
	}
	// 读取公共请求参数到结构体
	param := api.UnEncrypt{}
	EnParam := api.Encrypt{}
	if p.Ctx.Input.IsGet() {
		p.IsPost = false
	} else {
		p.IsPost = true
	}

	contentType := p.GetContentType()
	p.ParamContentType = contentType
	if p.IsPost == true && contentType == ApplicationJSON {
		_ = p.Ctx.BindJSON(&param)
		_ = p.Ctx.BindJSON(&EnParam)
	} else {
		param.Sign = p.GetString("sign", "")
		param.Appkey = p.GetString("appkey", "")
		param.Timestamp, _ = p.GetInt64("timestamp", 0)
		param.Action = p.GetString("action", "")
		param.Version = p.GetString("version", "")
		param.Mac = p.GetString("mac", "")
		EnParam.Sign = p.GetString("sign", "")
		EnParam.Signal = p.GetString("signal", "")
		EnParam.Encrypt = p.GetString("encrypt", "")
		EnParam.Ciphertext = p.GetString("ciphertext", "")
		EnParam.Timestamp, _ = p.GetInt64("timestamp", 0)

	}
	// 判断是否加密了
	p.IsEncrypt = false
	if EnParam.Signal != "" {
		param.Appkey = EnParam.Signal
		p.IsEncrypt = true
	}
	p.Encrypt = 0
	p.Param = param
	p.EnParam = EnParam
	for _, i := range []int{0: 0, 1: 1} {
		var e strings.Builder
		e.WriteString(strconv.Itoa(i))
		e.WriteString(p.EnParam.Signal)
		if strings.ToUpper(common.GetStringMd5(e.String())) == strings.ToUpper(p.EnParam.Encrypt) {
			p.Encrypt = i
			break
		}
	}
	if p.IsEncrypt == false && p.Project.Encrypt == 1 {
		p.Error(10004, "请使用加密访问")
	}

	// 如果使用了AES加密则调用AES解密
	if p.Encrypt == 1 {
		// 读取项目缓存信息
		projectString := ac.Get(p.Param.Appkey)
		if projectString != "" {
			err := json.Unmarshal([]byte(common.Strval(projectString)), &p.Project)
			if err != nil {
				p.Error(10004, "项目缓存出现错误")
			}
		} else {
			p.Error(10005, "项目不存在")
		}
		if p.Project.StatusType == 1 {
			p.Error(10006, "已停止运营")
		}
		b, err := valid.Valid(&p.EnParam)
		if err != nil {
			logs.Error("验证器错误", err)
		}
		if !b {
			for _, err := range valid.Errors {
				logs.Error(err.Key, err.Message)
				p.Error(9999, err.Message)
			}
		}
		script := aes
		vm := goja.New()
		_, err = vm.RunString(script)
		if err != nil {
			logs.Error("加密脚本出现问题")
			p.Error(10001, "加密脚本出现问题")
		}
		var decrypt func(string, string, string) string
		err = vm.ExportTo(vm.Get("AesDecrypt"), &decrypt)
		if err != nil {
			logs.Error("加密脚本映射到Go函数失败")
			p.Error(10002, "加密脚本映射到Go函数失败")
		}

		p.decryptJsonString = decrypt(p.EnParam.Ciphertext, p.Project.KeyA, p.Project.KeyB)
		err = json.Unmarshal([]byte(p.decryptJsonString), &p.Param)
		if err != nil {
			p.Error(10003, "解密失败,请检查密匙是否正确")
		}
		//p.decryptJsonString, err = common.AesEncrypt([]byte(p.EnParam.Ciphertext), []byte(p.Project.KeyA), []byte(p.Project.KeyB))
		//if err != nil {
		//	logs.Error("解密失败:", err, p.Project.KeyA, p.Project.KeyB)
		//	p.Error(10003, "解密失败,请检查密匙是否正确")
		//}
		//logs.Info("解密字符串：", p.decryptJsonString)
		//err = json.Unmarshal([]byte(p.decryptJsonString), &p.Param)
		//if err != nil {
		//	p.Error(10003, "解密失败,请检查密匙是否正确")
		//}
	}

	// 验证请求参数是否合法
	b, err := valid.Valid(&p.Param)
	if err != nil {
		logs.Error("验证器错误", err)
	}
	if !b {
		for _, err := range valid.Errors {
			logs.Error(err.Key, err.Message)
			p.Error(9999, err.Message)
		}
	}
	// 读取项目缓存信息
	projectString := ac.Get(p.Param.Appkey)
	if projectString != "" {
		err = json.Unmarshal([]byte(common.Strval(projectString)), &p.Project)
		if err != nil {
			p.Error(10004, "项目缓存出现错误")
		}
	} else {
		p.Error(10005, "项目不存在")
	}
	if p.Project.StatusType == 1 {
		p.Error(10006, "已停止运营")
	}
	// 检查签名是否正确
	var t strings.Builder
	t.WriteString(p.Project.AppKey)
	t.WriteString(p.Project.SecretKey)
	t.WriteString(p.Param.Version)
	t.WriteString(strconv.FormatInt(p.Param.Timestamp, 10))
	t.WriteString(p.Param.Mac)
	re := md5.Sum([]byte(t.String()))
	sign := fmt.Sprintf("%x", re)
	if strings.ToUpper(sign) != strings.ToUpper(p.Param.Sign) {
		p.Error(10007, "签名效验不通过")
	}
	//logs.Info("缓存的软件信息：", p.Project)
	// 检查当前版本是否存在
	var s strings.Builder
	s.WriteString(strconv.Itoa(p.Project.ID))
	s.WriteString("-")
	s.WriteString(p.Param.Version)
	versionString := ac.Get(s.String())
	if versionString != "" {
		err = json.Unmarshal([]byte(common.Strval(versionString)), &p.ProjectVersion)
		if err != nil {
			p.Error(10008, "版本号缓存提取失败")
		}
	} else {
		p.Error(10009, "请创建对应可用的版本号")
	}
	if p.Param.Action != "soft.init" && p.ProjectVersion.IsActive > 0 {
		p.Error(10017, "当前版本已禁用,请更新")
	}
	// 检查最新版本是否存在
	var u strings.Builder
	u.WriteString("id-")
	u.WriteString(strconv.Itoa(p.Project.ID))
	u.WriteString("-up-version")
	versionUpString := ac.Get(u.String())
	if versionUpString != "" {
		err = json.Unmarshal([]byte(common.Strval(versionString)), &p.ProjectUpVersion)
		if err != nil {
			p.Error(10010, "最新版本号缓存提取失败")
		}
	} else {
		ProjectUpVersion := models.ProjectVersion{}
		status, p.ProjectUpVersion = ProjectUpVersion.GetUpVersion()
	}

	// 检查是否绑定了登录规则
	if p.Project.LoginType == 0 {
		p.Error(10011, "请绑定登录规则")
	}
	var l strings.Builder
	l.WriteString("cache-project-login-")
	l.WriteString(strconv.Itoa(p.Project.LoginType))
	loginString := ac.Get(l.String())
	if loginString != "" {
		err = json.Unmarshal([]byte(common.Strval(loginString)), &p.ProjectLogin)
		if err != nil {
			p.Error(10012, "登录规则缓存提取失败")
		}
	} else {
		p.Error(10013, "登录规则不存在")
	}
	// 检查是否某些操作是否需要验证client字段
	Uri := p.Ctx.Request.RequestURI
	a := []string{"/api/user/login"}
	if In(Uri, a) {
		token := p.GetString("client")
		id := ac.Get(token)
		Id := common.GetInterfaceToInt(id)
		logs.Info("当前登录的ID：", Id)
		if id == "" || Id <= 0 {
			p.Error(10014, "请重新登录")
		}
	}
}

// 移除指定在线
func (p *BaseController) RemoveHeart(u models.Member, client string, msg string) {
	hearList := p.FetchHeartOnline(u)
	for _, i := range hearList.Online {
		if client == i.Client {
			var s strings.Builder
			s.WriteString("h-")
			s.WriteString(i.Client)
			_ = p.Ac.Delete(s.String())
			var out strings.Builder
			out.WriteString("h-o-")
			out.WriteString(i.Client)
			_ = p.Ac.Put(out.String(), msg, 1*60*60*time.Second)
		}
	}
}

// 移除单机所有在线
func (p *BaseController) RemoveHeartAllByMac(u models.Member, mac string, msg string) {
	hearList := p.FetchHeartOnline(u)
	for _, i := range hearList.Online {
		if i.Mac == mac {
			var s strings.Builder
			s.WriteString("h-")
			s.WriteString(i.Client)
			_ = p.Ac.Delete(s.String())
			var out strings.Builder
			out.WriteString("h-o-")
			out.WriteString(i.Client)
			_ = p.Ac.Put(out.String(), msg, 1*60*60*time.Second)
		}

	}
}

// 移除所有在线
func (p *BaseController) RemoveHeartAll(u models.Member, msg string) {
	hearList := p.FetchHeartOnline(u)
	for _, i := range hearList.Online {
		var s strings.Builder
		s.WriteString("h-")
		s.WriteString(i.Client)
		_ = p.Ac.Delete(s.String())
		var out strings.Builder
		out.WriteString("h-o-")
		out.WriteString(i.Client)
		_ = p.Ac.Put(out.String(), msg, 1*60*60*time.Second)
	}
	var h strings.Builder
	h.WriteString(heartPrefix)
	h.WriteString(strconv.Itoa(u.ProjectId))
	h.WriteString("-")
	h.WriteString(strconv.Itoa(u.ID))
	_ = p.Ac.Delete(h.String())
}

// 遍历所有在线
func (p *BaseController) FetchHeartOnline(u models.Member) HeartList {
	var h strings.Builder
	var heartMacList []string
	var onlineArray []api.OnlineData
	var pcOnlineArray []api.OnlineData
	heartList := HeartList{}
	h.WriteString(heartPrefix)
	h.WriteString(strconv.Itoa(u.ProjectId))
	h.WriteString("-")
	h.WriteString(strconv.Itoa(u.ID))
	heartJson := p.Ac.Get(h.String())

	if heartJson != "" {
		err := json.Unmarshal([]byte(common.Strval(heartJson)), &heartMacList)
		if err != nil {
			p.CallErrorJson("读取在线列表失败", nil)
		}
		heartList.List = heartMacList
		for _, i := range heartList.List {
			online := api.OnlineData{}
			var s strings.Builder
			s.WriteString("h-")
			s.WriteString(i)
			onlineJson := p.Ac.Get(s.String())
			if onlineJson != "" {
				err = json.Unmarshal([]byte(common.Strval(onlineJson)), &online)
				if err == nil {
					onlineArray = append(onlineArray, online)
					if online.Mac == p.Param.Mac {
						pcOnlineArray = append(pcOnlineArray, online)
					}
				}
			}
		}
		heartList.Online = onlineArray
		heartList.PcOnline = pcOnlineArray
		return heartList
	} else {
		return HeartList{}
	}
}

// 写到在线列表缓存
func (p *BaseController) InsertHeartList(u models.Member, o api.OnlineData) {
	var h strings.Builder
	var heartMacList []string
	var s strings.Builder
	h.WriteString(heartPrefix)
	h.WriteString(strconv.Itoa(u.ProjectId))
	h.WriteString("-")
	h.WriteString(strconv.Itoa(u.ID))
	heartJson := p.Ac.Get(h.String())
	s.WriteString("h-")
	s.WriteString(o.Client)
	data, _ := json.Marshal(o)
	if heartJson != "" {
		err := json.Unmarshal([]byte(common.Strval(heartJson)), &heartMacList)
		if err != nil {
			p.CallErrorJson("读取在线列表失败", nil)
		}
		heartMacList = append(heartMacList, o.Client)
	} else {
		heartMacList = append(heartMacList, o.Client)
	}
	heartMacJson, _ := json.Marshal(heartMacList)
	_ = p.Ac.Put(h.String(), string(heartMacJson), 2*60*60*time.Second)
	_ = p.Ac.Put(s.String(), string(data), 30*60*time.Second)
}
