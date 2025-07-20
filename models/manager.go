package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"
	"verification/controllers/common"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type Manager struct {
	ID         int       `orm:"column(id)" json:"id"`
	Email      string    `orm:"size(60);null;index;description(邮箱)" json:"email" valid:"MaxSize(32)"`
	User       string    `orm:"index;size(60);description(用户名)" json:"user" valid:"MaxSize(60)"`
	Pwd        string    `orm:"size(60);description(密码)" json:"pwd" valid:"MaxSize(60)"`
	Level      int       `orm:"default(0);description(等级)" json:"level" valid:"MaxSize(2)"`
	Avatar     string    `orm:"size(240);null;description(头像)" json:"avatar" valid:"MaxSize(240)"`
	Scope      int       `orm:"default(0);description(积分)" json:"scope" valid:"MaxSize(5)"`
	Pid        int       `orm:"index;default(0);description(关联上级ID)" json:"pid"`
	InviteId   int       `orm:"index;default(0);description(邀请人ID)" json:"invite_id"`
	Invite     string    `orm:"size(6);null;description(邀请码)" json:"invite" valid:"MaxSize(6)"`
	Money      float64   `orm:"digits(12);decimals(2);default(0);description(关联账户余额)" json:"money" valid:"MaxSize(8)"`
	LoginTime  int       `orm:"index;default(0);description(登录时间)" json:"login_time"`
	LoginIp    string    `orm:"size(30);null;description(登录IP)" json:"login_ip"`
	IsLock     int       `orm:"default(0);description(0正常,1锁定);index" json:"is_lock" valid:"MaxSize(1)"`
	IsDel      int       `orm:"default(0);description(0正常,1删除);index" json:"is_del" valid:"MaxSize(1)"`
	PowerId    int       `orm:"default(0);description(关联权限ID)" json:"power_id"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);index" json:"create_time"`
}

type ManagerTree struct {
	ID         int           `orm:"column(id)" json:"id"`
	Email      string        `orm:"size(60);null;index;description(邮箱)" json:"email" valid:"MaxSize(32)"`
	User       string        `orm:"index;size(60);description(用户名)" json:"user" valid:"MaxSize(60)"`
	Pwd        string        `orm:"size(60);description(密码)" json:"pwd" valid:"MaxSize(60)"`
	Level      int           `orm:"default(0);description(等级)" json:"level" valid:"MaxSize(2)"`
	Avatar     string        `orm:"size(240);null;description(头像)" json:"avatar" valid:"MaxSize(240)"`
	Scope      int           `orm:"default(0);description(积分)" json:"scope" valid:"MaxSize(5)"`
	Pid        int           `orm:"index;default(0);description(关联上级ID)" json:"pid"`
	InviteId   int           `orm:"index;default(0);description(邀请人ID)" json:"invite_id"`
	Invite     string        `orm:"size(6);null;description(邀请码)" json:"invite" valid:"MaxSize(6)"`
	Money      float64       `orm:"digits(12);decimals(2);default(0);description(关联账户余额)" json:"money" valid:"MaxSize(8)"`
	LoginTime  int           `orm:"index;default(0);description(登录时间)" json:"login_time"`
	LoginIp    string        `orm:"size(30);null;description(登录IP)" json:"login_ip"`
	IsLock     int           `orm:"default(0);description(0正常,1锁定);index" json:"is_lock" valid:"MaxSize(1)"`
	IsDel      int           `orm:"default(0);description(0正常,1删除);index" json:"is_del" valid:"MaxSize(1)"`
	PowerId    int           `orm:"default(0);description(关联权限ID)" json:"power_id"`
	Children   []ManagerTree `json:"children"`
	CreateTime time.Time     `orm:"auto_now_add;type(datetime);index" json:"create_time"`
}

type Permission struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type ManagerInfo struct {
	UserId      int          `json:"user_id"`
	Username    string       `json:"username"`
	RealName    string       `json:"real_name"`
	Avatar      string       `json:"avatar"`
	Desc        string       `json:"desc"`
	Token       string       `json:"token"`
	Permissions []Permission `json:"permissions"`
	User        interface{}  `json:"user"`
	Version     string       `json:"version"`
	Update      bool         `json:"update"`
	NewVersion  string       `json:"new_version"`
	Notice      string       `json:"notice"`
}

var DeveloperPermission = []Permission{
	{
		Label: "开发者",
		Value: "developer",
	},
}
var AgentPermission = []Permission{
	{
		Label: "代理",
		Value: "agent",
	},
}

func (m *Manager) Login(username string, password string) (error string, user *Manager) {
	o := orm.NewOrm()
	u := Manager{}
	u.User = username
	err := o.Read(&u, "User")
	if err != nil {
		return "账号或者密码错误", nil
	}
	if password != u.Pwd {
		return "密码错误", nil
	}
	return "", &u
}

func GetManagerById(id int) (ManagerTree, bool) {
	res := ManagerTree{}
	manager := Manager{ID: id}
	o := orm.NewOrm()
	err := o.Read(&manager)
	if err != nil {
		return res, false
	}
	data, _ := json.Marshal(manager)
	_ = json.Unmarshal(data, &res)
	return res, true
}

func toTree(manager Manager) ManagerTree {
	res := ManagerTree{}
	data, _ := json.Marshal(manager)
	_ = json.Unmarshal(data, &res)
	return res
}

func (m *Manager) InitManager() {
	o := orm.NewOrm()
	u := Manager{}
	u.ID = 1
	err := o.Read(&u)
	if err != nil {
		user := new(Manager)
		user.User = "admin"
		user.Pwd = "112233"
		id, err := o.Insert(user)
		if err != nil {
			fmt.Printf("管理员创建失败")
		}
		if id > 0 {
			fmt.Printf("管理员创建成功")
		}
	}
}

func (m *Manager) Add() (int64, error) {
	o := orm.NewOrm()
	return o.Insert(m)
}

func GetManagerTree(id int) []ManagerTree {
	var idList []ManagerTree
	o := orm.NewOrm()
	qs := o.QueryTable("Manager")
	var list []Manager
	_, err := qs.Filter("Pid", id).All(&list)
	if err != nil {
		return idList
	}
	var child []ManagerTree
	for _, i := range list {
		var cList []Manager
		_, _ = qs.Filter("ID", i.ID).All(&cList)
		for _, item := range cList {
			treeItem := toTree(item)
			next := GetManagerTree(i.ID)
			if len(next) > 0 {
				treeItem.Children = next
			}
			child = append(child, treeItem)
		}
	}
	return child
}

func (m *Manager) GetManagerList(managerId int, managerIdList []int) (status bool, pager Pager) {
	var data []ManagerTree
	o := orm.NewOrm()
	count, _ := o.QueryTable(&m).Filter("ID__in", managerIdList).Count()
	data = GetManagerTree(managerId)
	return true, Pager{
		Count:       count,
		CurrentPage: 1,
		Data:        data,
		PageSize:    int64(len(data)),
		TotalPages:  1,
	}
}

func (m *Manager) UpdateById(n Manager) (int64, string) {
	managerId := m.ID
	o := orm.NewOrm()
	exist := Manager{User: n.User}
	err := o.Read(&exist, "User")
	if err == nil {
		if exist.ID != managerId {
			return 0, "该账号名称已被使用"
		}
	}
	err = o.Read(m)
	if err != nil {
		return 0, "账号不存在"
	}
	status, _ := regexp.MatchString("^[a-zA-Z0-9._-]{5,32}$", n.User)
	if status == false {
		return 0, "账号格式:长度5到32，字母/数字/字母数字组合"
	}
	m.User = n.User
	if n.Pwd != "" {
		status, _ = regexp.MatchString("^[a-zA-Z0-9._-]{5,32}$", n.Pwd)
		if status == false {
			return 0, "密码格式:长度5到32，字母/数字/字母数字组合"
		}
		m.Pwd = n.Pwd
	}
	status, _ = regexp.MatchString("^([a-zA-Z]|[0-9])(\\w|\\-)+@[a-zA-Z0-9]+\\.([a-zA-Z]{2,4})$", n.Email)
	if status == false {
		return 0, "邮箱格式错误"
	}
	m.Email = n.Email
	row, err := o.Update(m)
	if err != nil {
		return 0, "更新失败"
	}
	return row, "更新成功"
}

func (m *Manager) GetInfoById(id int) (Manager, error) {
	u := Manager{ID: id}
	o := orm.NewOrm()
	err := o.Read(&u)
	return u, err
}

func (m *Manager) GetManagerIdList(id int, managerList *[]int) {
	var list []int
	status, ac := common.GetCacheAC()
	if status == false {
		logs.Error("缓存初始化失败")
	}
	c := ac.Get("manager-" + strconv.Itoa(id))
	if c == "" {
		GetManagerList(id, managerList)

		l, err := json.Marshal(managerList)
		if err != nil {
			logs.Error("序列化失败", err)
		} else {
			_ = ac.Put("manager-"+strconv.Itoa(id), string(l), 365*24*60*60*time.Second)
		}
	} else {
		_ = json.Unmarshal([]byte(common.Strval(c)), &list)
		*managerList = list
	}
}

func GetManagerList(id int, managerList *[]int) []int {
	var idList []int
	idList = append(idList, id)
	o := orm.NewOrm()
	qs := o.QueryTable("Manager")
	var list []Manager
	_, err := qs.Filter("Pid", id).All(&list)
	if err != nil {
		return idList
	}
	for _, i := range list {
		if i.IsLock == 0 {
			idList = append(idList, GetManagerList(i.ID, managerList)...)
		}
	}
	*managerList = idList
	return idList
}

func GetUpManagerList(id int, managerList *[]int) []int {
	var idList []int
	idList = append(idList, id)
	o := orm.NewOrm()
	qs := o.QueryTable("Manager")
	var list []Manager
	_, err := qs.Filter("ID", id).All(&list)
	if err != nil {
		return idList
	}
	for _, i := range list {
		if i.IsLock == 0 {
			idList = append(idList, GetManagerList(i.Pid, managerList)...)
		}
	}
	*managerList = idList
	return idList
}

func (m *Manager) AddManager() (int64, string) {
	o := orm.NewOrm()
	_ = o.Read(m, "User")
	if m.ID > 0 {
		return 0, "账号已被使用"
	}
	status, _ := regexp.MatchString("^[a-zA-Z0-9._-]{5,32}$", m.User)
	if status == false {
		return 0, "账号格式:长度5到32，字母/数字/字母数字组合"
	}
	status, _ = regexp.MatchString("^[a-zA-Z0-9._-]{5,32}$", m.Pwd)
	if status == false {
		return 0, "密码格式:长度5到32，字母/数字/字母数字组合"
	}
	status, _ = regexp.MatchString("^([a-zA-Z]|[0-9])(\\w|\\-)+@[a-zA-Z0-9]+\\.([a-zA-Z]{2,4})$", m.Email)
	if status == false {
		return 0, "邮箱格式错误"
	}
	_ = o.Read(m, "Email")
	if m.ID > 0 {
		return 0, "邮箱已被使用"
	}
	id, err := o.Insert(m)
	if err != nil {
		logs.Error("创建代理失败", err)
		return 0, "创建失败"
	}
	status, ac := common.GetCacheAC()
	if status == false {
		logs.Error("缓存初始化失败")
	}
	_ = ac.Delete("manager-" + strconv.Itoa(m.Pid))
	return id, "创建成功"
}

func (m *Manager) UpdateManager(n Manager) (int64, string) {
	o := orm.NewOrm()
	id := m.ID
	user := n.User
	email := n.Email
	pwd := n.Pwd
	powerId := n.PowerId
	money := n.Money
	_ = o.Read(&n, "User")
	if n.ID != id {
		return 0, "账号已被使用"
	}
	status, _ := regexp.MatchString("^[a-zA-Z0-9._-]{5,32}$", n.User)
	if status == false {
		return 0, "账号格式:长度5到32，字母/数字/字母数字组合"
	}
	status, _ = regexp.MatchString("^[a-zA-Z0-9._-]{5,32}$", n.Pwd)
	if status == false {
		return 0, "密码格式:长度5到32，字母/数字/字母数字组合"
	}
	status, _ = regexp.MatchString("^([a-zA-Z]|[0-9])(\\w|\\-)+@[a-zA-Z0-9]+\\.([a-zA-Z]{2,4})$", n.Email)
	if status == false {
		return 0, "邮箱格式错误"
	}
	_ = o.Read(&n, "Email")
	if n.ID != id {
		return 0, "邮箱已被使用"
	}
	err := o.Read(m)
	if err != nil {
		return 0, "代理账号不存在"
	}
	m.User = user
	m.Email = email
	m.Pwd = pwd
	m.PowerId = powerId
	m.Money = money
	row, err := o.Update(m)
	if err != nil {
		logs.Error("修改代理失败", err)
		return 0, "修改失败"
	}
	status, ac := common.GetCacheAC()
	if status == false {
		logs.Error("缓存初始化失败")
	}
	_ = ac.Delete("manager-" + strconv.Itoa(m.Pid))
	_ = ac.Delete("manager-info-" + strconv.Itoa(m.ID))
	return row, "修改成功"
}

func (m *Manager) Delete() int64 {
	status, ac := common.GetCacheAC()
	if status == false {
		logs.Error("缓存初始化失败")
	}
	var upManagerIdList []int
	var managerIdList []int
	o := orm.NewOrm()
	err := o.Read(m)
	if err != nil {
		return 0
	}

	GetUpManagerList(m.ID, &upManagerIdList)
	for _, i := range upManagerIdList {
		_ = ac.Delete("manager-" + strconv.Itoa(i))
	}
	GetManagerList(m.ID, &managerIdList)
	if m.Pid == 0 {
		return 0
	}
	row, _ := o.Delete(m)
	member := Member{}
	qs := o.QueryTable(&member)
	_, _ = qs.Filter("ManagerId__in", managerIdList).Delete()
	keys := Keys{}
	qs = o.QueryTable(&keys)
	_, _ = qs.Filter("ManagerId__in", managerIdList).Delete()
	memberLogin := MemberLogin{}
	qs = o.QueryTable(&memberLogin)
	_, _ = qs.Filter("ManagerId__in", managerIdList).Delete()
	order := Order{}
	qs = o.QueryTable(&order)
	_, _ = qs.Filter("ManagerId__in", managerIdList).Delete()
	managerCards := ManagerCards{}
	qs = o.QueryTable(&managerCards)
	_, _ = qs.Filter("ManagerId__in", managerIdList).Delete()
	return row
}

func (m *Manager) GetAllAgent(managerIdArr []int) []Manager {
	data := make([]Manager, 0)
	o := orm.NewOrm()
	qs := o.QueryTable(&m)
	_, err := qs.Filter("ID__in", managerIdArr).All(&data, "User", "Email", "ID")
	if err != nil {
		return data
	}
	return data
}
