package models

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"time"
	"verification/controllers/common"
)

type Project struct {
	ID         int       `orm:"column(id)" json:"id"`
	ManagerId  int       `orm:"index;default(0)" json:"manager_id"`
	Name       string    `orm:"size(40)" json:"name" valid:"MaxSize(40)"`
	AppKey     string    `orm:"index;size(32)" json:"app_key" valid:"MaxSize(32)"`
	SecretKey  string    `orm:"size(32)" json:"secret_key" valid:"MaxSize(32)"`
	Type       int       `orm:"default(0);description(0单码,1账号)" json:"type" valid:"Max(1);MaxSize(1)"`
	StatusType int       `orm:"default(0);description(0收费,1停止,2免费)" json:"status_type"`
	PublicKey  string    `orm:"type(text);size(4000);description(RSA公钥);null" json:"public_key" valid:"MaxSize(4000)"`
	PrivateKey string    `orm:"type(text);size(4000);description(RSA私钥);null" json:"private_key" valid:"MaxSize(4000)"`
	KeyA       string    `orm:"size(16);description(AES加密密匙A);null" json:"key_a" valid:"MaxSize(16)"`
	KeyB       string    `orm:"size(16);description(AES加密密匙B);null" json:"key_b" valid:"MaxSize(16)"`
	Sign       int       `orm:"default(0);description(0 MD5 1 SHA1 2 SHA224 3 SHA256 4 SHA384 5 SHA512)" json:"sign" valid:"MaxSize(1)"`
	Encrypt    int       `orm:"default(0);description(0开放签名API,1AES)" json:"encrypt" valid:"MaxSize(1)"`
	LoginType  int       `orm:"default(0)" json:"login_type" valid:"MaxSize(1)"`
	GiftType   int       `orm:"default(0)" json:"gift_type" valid:"MaxSize(1)"`
	Notice     string    `orm:"type(text);description(公告);size(9000);null" json:"notice" valid:"MaxSize(9000)"`
	Api        string    `orm:"type(text);description(公告);size(5000);null" json:"api" valid:"MaxSize(5000)"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);index" json:"create_time"`
}

func (p *Project) GetProjectList(pageSize int64, page int64) (status bool, pager Pager) {
	var data []Project
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

func (p *Project) ProjectList() (status bool, pager Pager) {
	var data []Project
	o := orm.NewOrm()
	_, err := o.QueryTable(&p).All(&data, "ID", "Name")
	if err != nil {
		return false, Pager{}
	}
	return true, Pager{
		Count:       0,
		CurrentPage: 0,
		Data:        data,
		PageSize:    0,
		TotalPages:  0,
	}
}

func (p *Project) GetAgentProjectList(list []int) (status bool, pager Pager) {
	var data []Project
	o := orm.NewOrm()
	_, err := o.QueryTable(&p).Filter("ID__in", list).All(&data, "ID", "Name")
	if err != nil {
		return false, Pager{}
	}
	return true, Pager{
		Count:       0,
		CurrentPage: 0,
		Data:        data,
		PageSize:    0,
		TotalPages:  0,
	}
}

func (c *Project) Add(name string, projectType int, statusType int, encrypt int, notice string, api string, managerId int, sign int) (projectId int64) {
	status, publicKey, privateKey := GetRsaKey()
	if status == false {
		return 0
	}
	p := Project{
		ManagerId:  managerId,
		Name:       name,
		AppKey:     RandStr(32),
		SecretKey:  RandStr(32),
		Type:       projectType,
		StatusType: statusType,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		KeyA:       RandStr(16),
		KeyB:       RandStr(16),
		Encrypt:    encrypt,
		LoginType:  0,
		GiftType:   0,
		Notice:     notice,
		Api:        api,
		Sign:       sign,
	}
	o := orm.NewOrm()
	id, err := o.Insert(&p)
	if err != nil {
		return 0
	}
	_, ac := common.GetCacheAC()
	data, err := json.Marshal(&p)
	_ = ac.Put(common.GetStringMd5(p.AppKey), string(data), 365*24*60*60*time.Second)
	_ = ac.Put(p.AppKey, string(data), 365*24*60*60*time.Second)
	return id
}

func (c *Project) Update(id int, name string, projectType int, statusType int, encrypt int, notice string, api string, updateRsa int, updateKey int, updateAppKey int, updateSecretKey int, sign int) bool {
	p := Project{
		ID: id,
	}
	_, ac := common.GetCacheAC()
	o := orm.NewOrm()
	err := o.Read(&p)
	if err != nil {
		logs.Error("查询项目失败")
		return false
	}
	p.Name = name
	p.Type = projectType
	p.StatusType = statusType
	p.Encrypt = encrypt
	p.Sign = sign
	if notice != "" {
		p.Notice = notice
	}
	if api != "" {
		p.Api = api
	}
	if updateRsa > 0 {
		status, publicKey, privateKey := GetRsaKey()
		if status == false {
			return false
		}
		p.PublicKey = publicKey
		p.PrivateKey = privateKey
	}
	if updateKey > 0 {
		p.KeyA = RandStr(16)
		p.KeyB = RandStr(16)
	}
	if updateAppKey > 0 {
		_ = ac.Delete(common.GetStringMd5(p.AppKey))
		_ = ac.Delete(p.AppKey)
		p.AppKey = RandStr(32)
	}
	if updateSecretKey > 0 {
		p.SecretKey = RandStr(32)
	}

	row, err := o.Update(&p)
	if err != nil {
		return false
	}
	if row > 0 {

		data, _ := json.Marshal(&p)
		_ = ac.Put(common.GetStringMd5(p.AppKey), string(data), 3*365*24*60*60*time.Second)
		_ = ac.Put(p.AppKey, string(data), 3*365*24*60*60*time.Second)
		return true
	}
	return false
}

func (c *Project) Delete(id int) bool {
	p := Project{
		ID: id,
	}
	o := orm.NewOrm()
	err := o.Read(&p)
	if err != nil {
		logs.Error("查询项目失败")
		return false
	}
	rows, err := o.Delete(&p)
	if rows > 0 {
		_, err = o.Raw("delete from cards where project_id = ?", p.ID).Exec()
		_, err = o.Raw("delete from keys where project_id = ?", p.ID).Exec()
		_, err = o.Raw("delete from project_version where project_id = ?", p.ID).Exec()
		_, err = o.Raw("delete from order where project_id = ?", p.ID).Exec()
		_, err = o.Raw("delete from member where project_id = ?", p.ID).Exec()
		_, err = o.Raw("delete from member_login where project_id = ?", p.ID).Exec()
		_, err = o.Raw("delete from manager_cards where project_id = ?", p.ID).Exec()
		_, ac := common.GetCacheAC()
		_ = ac.Delete(fmt.Sprintf(cacheIndex[0], id))
		_ = ac.Delete(common.GetStringMd5(p.AppKey))
		_ = ac.Delete(p.AppKey)
		return true
	}
	return false
}

func (c *Project) Bind(id int, projectLoginId int) (status bool, msg string) {
	o := orm.NewOrm()
	if projectLoginId > 0 {

		projectLogin := ProjectLogin{ID: projectLoginId}
		err := o.Read(&projectLogin)
		if err != nil {
			return false, "登录规则不存在"
		}
	}
	project := Project{ID: id}
	err := o.Read(&project)
	if err != nil {
		return false, "项目不存在"
	}
	if project.LoginType > 0 {
		project.LoginType = 0
		msg = "解绑登录规则成功"
	} else {
		project.LoginType = projectLoginId
		msg = "绑定登录规则成功"
	}

	row, err := o.Update(&project)
	if row > 0 {
		_, ac := common.GetCacheAC()
		data, _ := json.Marshal(&project)
		_ = ac.Put(fmt.Sprintf(cacheIndex[0], id), string(data), 3*365*24*60*60*time.Second)
		_ = ac.Put(common.GetStringMd5(project.AppKey), string(data), 3*365*24*60*60*time.Second)
		_ = ac.Put(project.AppKey, string(data), 3*365*24*60*60*time.Second)
		return true, msg
	}
	return false, "操作失败"
}

func (p *Project) GetCount() int64 {
	managerId := p.ManagerId
	o := orm.NewOrm()
	qs := o.QueryTable("Project")
	count, err := qs.Filter("ManagerId", managerId).Count()
	if err != nil {
		logs.Error("查询项目总数错误：", err)
		return 0
	}
	return count
}

func (p *Project) AllProject() []Project {
	data := make([]Project, 0)
	o := orm.NewOrm()
	_, err := o.QueryTable(&p).All(&data)
	if err != nil {
		return data
	}
	return data
}
