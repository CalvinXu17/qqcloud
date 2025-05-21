package models

import (
	"encoding/json"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"strconv"
	"strings"
	"time"
	"verification/controllers/common"
)

type ProjectLogin struct {
	ID                 int       `orm:"column(id)" json:"id"`
	ManagerId          int       `orm:"index;default(0)" json:"manager_id"`
	Title              string    `orm:"size(40)" json:"title" valid:"MaxSize(40)"`
	Mode               int       `orm:"default(0);description(0绑定登录,1普通登录,2点数登录)" json:"mode" valid:"MaxSize(1)"`
	RegMode            int       `orm:"default(0);description(0带卡注册,1普通注册)" json:"reg_mode" valid:"MaxSize(1)"`
	EmailReg           int       `orm:"default(0);description(0关闭邮箱注册,1开启邮箱注册)" json:"email_reg" valid:"MaxSize(1)"`
	UnbindMode         int       `orm:"default(0);description(0 不允许解绑 1 原机解绑 2 自动解绑 3 任意解绑)" json:"unbind_mode" valid:"MaxSize(1)"`
	UnbindWeakenMode   int       `orm:"default(0);description(0 不扣时 1 解绑就扣时 2 超出扣时 3 超出不扣时)" json:"unbind_weaken_mode" valid:"MaxSize(1)"`
	UnbindWeaken       float64   `orm:"digits(12);decimals(2);" json:"unbind_weaken" valid:"MaxSize(5)"`
	UnbindWeakenPoints int       `orm:"default(1);description(解绑扣点)" json:"unbind_weaken_points" valid:"MaxSize(4)"`
	UnbindTimes        int       `orm:"default(3);description(最大可解绑次数)" json:"unbind_times" valid:"MaxSize(3)"`
	UnbindDate         int       `orm:"default(3);description(0 天 1 月)" json:"unbind_date" valid:"MaxSize(1)"`
	UnbindBefore       int       `orm:"default(0);description(解绑前置时间)" json:"unbind_before" valid:"MaxSize(4)"`
	NumberMode         int       `orm:"default(0);description(0 机器码控制 1 IP控制)" json:"number_mode" valid:"MaxSize(1)"`
	NumberMore         int       `orm:"default(10);description(最大限制)" json:"number_more" valid:"MaxSize(5)"`
	NumberWeaken       int       `orm:"default(1);description(扣除点数)" json:"number_weaken" valid:"MaxSize(4)"`
	NumberWeakenTime   int       `orm:"default(0);description(0 时 1天)" json:"number_weaken_time" valid:"MaxSize(1)"`
	PcMore             int       `orm:"default(0);description(单机最大可在线)" json:"pc_more" valid:"MaxSize(5)"`
	PcCodeMore         int       `orm:"default(0);description(最大多开上限)" json:"pc_code_more" valid:"MaxSize(5)"`
	CreateTime         time.Time `orm:"auto_now_add;type(datetime);index" json:"create_time"`
}

func (p *ProjectLogin) Add() int64 {
	o := orm.NewOrm()
	id, err := o.Insert(p)
	if err != nil {
		logs.Error(err)
		return 0
	}
	if id > 0 {
		_, ac := common.GetCacheAC()
		data, _ := json.Marshal(&p)
		var l strings.Builder
		l.WriteString("cache-project-login-")
		l.WriteString(strconv.FormatInt(id, 10))
		_ = ac.Put(l.String(), string(data), 3*365*24*60*60*time.Second)
		return id
	}
	return 0
}

func (p *ProjectLogin) LoginRuleList() (bool, Pager) {
	var data []ProjectLogin
	o := orm.NewOrm()
	_, err := o.QueryTable(&p).All(&data)
	if err != nil {
		return false, Pager{}
	}
	return true, Pager{
		Count:       0,
		CurrentPage: 1,
		Data:        data,
		PageSize:    0,
		TotalPages:  0,
	}
}

func (p *ProjectLogin) GetLoginRuleList(pageSize int64, page int64) (status bool, pager Pager) {
	var data []ProjectLogin
	o := orm.NewOrm()
	count, err := o.QueryTable(&p).Count()
	if err != nil {
		logs.Error(err)
		return false, Pager{}
	}
	totalPage, offset, currentPage := PageCount(count, pageSize, page)
	_, err = o.QueryTable(&p).Limit(pageSize, offset).All(&data)
	return true, Pager{
		Count:       count,
		CurrentPage: currentPage,
		Data:        data,
		PageSize:    pageSize,
		TotalPages:  totalPage,
	}
}

func (p *ProjectLogin) Update(m ProjectLogin) int {
	o := orm.NewOrm()
	err := o.Read(p)
	if err != nil {
		return 0
	}
	p.Mode = m.Mode
	if m.Title != "" {
		p.Title = m.Title
	}
	p.RegMode = m.RegMode
	p.EmailReg = m.EmailReg
	p.UnbindMode = m.UnbindMode
	p.UnbindWeakenMode = m.UnbindWeakenMode
	p.UnbindWeaken = m.UnbindWeaken
	p.UnbindWeakenPoints = m.UnbindWeakenPoints
	p.UnbindTimes = m.UnbindTimes
	p.UnbindDate = m.UnbindDate
	p.UnbindBefore = m.UnbindBefore
	p.NumberMode = m.NumberMode
	p.NumberMore = m.NumberWeaken
	p.NumberWeakenTime = m.NumberWeakenTime
	p.PcMore = m.PcMore
	p.PcCodeMore = m.PcCodeMore
	row, err := o.Update(p)
	if row > 0 {
		_, ac := common.GetCacheAC()
		ID := p.ID
		data, _ := json.Marshal(&p)
		var l strings.Builder
		l.WriteString("cache-project-login-")
		l.WriteString(strconv.Itoa(ID))
		_ = ac.Put(l.String(), string(data), 3*365*24*60*60*time.Second)
		logs.Error("登录规则缓存：", l.String(), string(data))
		return 1
	}
	return 0
}

func (p *ProjectLogin) Delete() int64 {
	o := orm.NewOrm()
	err := o.Read(p)
	if err != nil {
		return 0
	}
	project := Project{LoginType: p.ID}
	err = o.Read(&project, "LoginType")
	_, ac := common.GetCacheAC()
	if err == nil {
		project.LoginType = 0
		row, _ := o.Update(&project)
		if row > 0 {
			data, _ := json.Marshal(&project)
			_ = ac.Put(common.GetStringMd5(project.AppKey), string(data), 365*24*60*60*time.Second)
			_ = ac.Put(project.AppKey, string(data), 365*24*60*60*time.Second)
			logs.Info("已清除登录规则绑定", project.Name)
		}
	}
	row, err := o.Delete(p)
	if row > 0 {

		ID := p.ID
		var l strings.Builder
		l.WriteString("cache-project-login-")
		l.WriteString(strconv.Itoa(ID))
		_ = ac.Delete(l.String())
		return row
	}
	return 0
}

func (p *ProjectLogin) AllLogin() []ProjectLogin  {
	data := make([]ProjectLogin, 0)
	o := orm.NewOrm()
	_, err := o.QueryTable(&p).All(&data)
	if err != nil {
		return data
	}
	return data
}