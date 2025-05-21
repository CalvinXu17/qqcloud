package admin

import (
	"github.com/beego/beego/v2/core/validation"
	"regexp"
)

type CreateCards struct {
	Title      string  `json:"title" valid:"Required;MaxSize(28)" `
	Price      float64 `json:"price"`
	KeyPrefix  string  `json:"key_prefix" valid:"MaxSize(4)"`
	LevelId    int     `json:"level_id"`
	Days       float64 `json:"days"`
	Points     int     `json:"points" valid:""`
	KeyExtAttr string  `json:"key_ext_attr" valid:"MaxSize(199)"`
	Tag        string  `json:"tag" valid:"MaxSize(199)"`
	IsLock     int     `json:"is_lock"`
}

func (c *CreateCards) Valid(v *validation.Validation) {
	status, _ := regexp.MatchString("^[a-zA-Z0-9]{0,4}$", c.KeyPrefix)
	if status == false {
		_ = v.SetError("KeyPrefix", "前缀格式为数字字母组合最大4字符")
	}
	if c.Price > 9999 {
		_ = v.SetError("Price", "定价最大9999")
	}
	if c.Days > 9999 {
		_ = v.SetError("Price", "天数最大9999")
	}
	if c.Points > 9999 {
		_ = v.SetError("Price", "点数最大9999")
	}
}

type CreateKeys struct {
	CardsId    int    `json:"cards_id"`
	Count      int64  `json:"count"`
	Length     int    `json:"length"`
	CreateType int    `json:"create_type"`
	ProjectId  int    `json:"project_id"`
	Tag        string `json:"tag" valid:"MaxSize(200)""`
}

func (c *CreateKeys) Valid(v *validation.Validation) {
	if c.Count > 500 {
		_ = v.SetError("Count", "一次最多提取500张激活码")
	}
	if c.Length > 32 {
		_ = v.SetError("Length", "激活码长度最大32")
	}
	if c.Length < 8 {
		_ = v.SetError("Length", "激活码长度最小32")
	}
	if c.CreateType < 0 || c.CreateType > 3 {
		_ = v.SetError("Length", "激活码组合形式不存在")
	}
}
