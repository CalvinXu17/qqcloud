package api

import (
	"regexp"
	"strconv"

	"github.com/astaxie/beego/validation"
)

type Encrypt struct {
	Signal     string `json:"signal" valid:"Required;MaxSize(32);"`
	Sign       string `json:"sign" valid:"Required;MaxSize(32);"`
	Encrypt    string `json:"encrypt" valid:"Required;MaxSize(32);"`
	Timestamp  int64  `json:"timestamp" valid:"Required;"`
	Ciphertext string `json:"ciphertext" valid:"Required"`
}

type UnEncrypt struct {
	Appkey    string `json:"appkey" valid:"Required;MaxSize(32);"`
	Version   string `json:"version" valid:"Required;MaxSize(5);"`
	Sign      string `json:"sign" valid:"Required;MaxSize(32);"`
	Action    string `json:"action" valid:"Required;MaxSize(24);"`
	Timestamp int64  `json:"timestamp" valid:"Required"`
	Mac       string `json:"mac" valid:"Required;MaxSize(32);MinSize(3)"`
}

func (u *UnEncrypt) Valid(v *validation.Validation) {
	status, _ := regexp.MatchString("^[A-Za-z0-9]*$", u.Appkey)
	if status == false {
		_ = v.SetError("Appkey", "appkey只能为数字字母组合")
	}
	status, _ = regexp.MatchString("^[+-]?[0-9]+(\\.[0-9]{1,4})?$", u.Version)
	if status == false {
		_ = v.SetError("Version", "version不合法")
	}
	status, _ = regexp.MatchString("^[A-Za-z0-9]*$", u.Sign)
	if status == false {
		_ = v.SetError("Sign", "sign只能为数字字母组合")
	}
	status, _ = regexp.MatchString("^[0-9]*$", strconv.FormatInt(u.Timestamp, 10))
	if status == false {
		_ = v.SetError("Timestamp", "timestamp只能为数字组合")
	}
	status, _ = regexp.MatchString("^[A-Za-z0-9]*$", u.Mac)
	if status == false {
		_ = v.SetError("Mac", "mac只能为数字字母组合")
	}
}

func (u *Encrypt) Valid(v *validation.Validation) {
	status, _ := regexp.MatchString("^[A-Za-z0-9]*$", u.Signal)
	if status == false {
		_ = v.SetError("Signal", "signal只能为数字字母组合")
	}
	status, _ = regexp.MatchString("^[A-Za-z0-9]*$", u.Sign)
	if status == false {
		_ = v.SetError("Sign", "sign只能为数字字母组合")
	}
	status, _ = regexp.MatchString("^[A-Za-z0-9]*$", u.Encrypt)
	if status == false {
		_ = v.SetError("Encrypt", "encrypt只能为数字字母组合")
	}
	status, _ = regexp.MatchString("^[0-9]*$", strconv.FormatInt(u.Timestamp, 10))
	if status == false {
		_ = v.SetError("Timestamp", "timestamp只能为数字组合")
	}
}

type RegisterParam struct {
	Email       string `json:"email"`
	User        string `json:"user"`
	Pwd         string `json:"pwd"`
	Pwd2        string `json:"pwd2"`
	Key         string `json:"key"`
	Recommender string `json:"recommender"`
	Code        string `json:"code"`
	Captcha     string `json:"captcha"`
}

func (u *RegisterParam) Valid(v *validation.Validation) {
	status, _ := regexp.MatchString("^[a-zA-Z0-9._-]{5,32}$", u.User)
	if status == false {
		_ = v.SetError("User", "账号格式:长度5到32，字母/数字/字母数字组合")
	}
	status, _ = regexp.MatchString("^[a-zA-Z0-9._-]{5,32}$", u.Pwd)
	if status == false {
		_ = v.SetError("Pwd", "密码格式:长度5到32，字母/数字/字母数字组合")
	}
	status, _ = regexp.MatchString("^[a-zA-Z0-9._-]{5,32}$", u.Pwd2)
	if status == false {
		_ = v.SetError("Pwd2", "安全密码格式:长度5到32，字母/数字/字母数字组合")
	}
	status, _ = regexp.MatchString("^([a-zA-Z]|[0-9])(\\w|\\-)+@[a-zA-Z0-9]+\\.([a-zA-Z]{2,4})$", u.Email)
	if status == false {
		_ = v.SetError("Email", "邮箱不正确")
	}
	status, _ = regexp.MatchString("^[A-Za-z0-9]*$", u.Captcha)
	if status == false {
		_ = v.SetError("Captcha", "验证码密码为字母数字组合")
	}
}

type LoginParam struct {
	User string `json:"user" valid:"Required"`
	Pwd  string `json:"pwd" valid:"Required"`
}

type ForgetParam struct {
	User    string `json:"user" valid:"Required"`
	Pwd     string `json:"pwd" valid:"Required"`
	Pwd2    string `json:"pwd2"`
	Code    string `json:"code"`
	Captcha string `json:"captcha"`
}

type RechargeParam struct {
	User string `json:"user" valid:"Required"`
	Key  string `json:"key" valid:"Required"`
}

type PointsParam struct {
	User   string `json:"user" valid:"Required"`
	Pwd    string `json:"pwd" valid:"Required"`
	Number int    `json:"number" valid:"Required"`
}

type HeartParm struct {
	Client string `json:"client" valid:"Required"`
}

type ClientParam struct {
	User string `json:"user" valid:"Required"`
	Pwd  string `json:"pwd" valid:"Required"`
}

type LogoutParm struct {
	Client string `json:"client" valid:"Required"`
	Type   int    `json:"type" valid:"Required"`
}

type OnlineData struct {
	Mac        string `json:"mac"`
	Addtime    int    `json:"addtime"`
	ClientTime int    `json:"clienttime"`
	Ip         string `json:"ip"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	NickName   string `json:"nick_name"`
	MemberId   int    `json:"memberId"`
	ManagerId  int    `json:"managerId"`
	ProjectId  int    `json:"projectId"`
	Client     string `json:"client"`
}
