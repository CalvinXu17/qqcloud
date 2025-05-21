package models

import (
	"github.com/beego/beego/v2/client/orm"
	"time"
)

type RoleItem struct {
	ID         int       `orm:"column(id)" json:"id"`
	RoleId     int       `orm:"index;default(0);description(归属角色)" json:"role_id"`
	Path       string    `orm:"null;description(请求路由);size(50)" json:"path"`
	Index      string    `orm:"null;description(细粒化请求Key);size(50)" json:"index"`
	Value      string    `orm:"null;description(细粒化请求value);size(50)" json:"value"`
	Name       string    `orm:"null;description(权限名称);size(50)" json:"description"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);index" json:"create_time"`
}

type role struct {
	Path  string `json:"path"`
	Index string `json:"index"`
	Value string `json:"value"`
	Name  string `json:"name"`
	Child []role `json:"children"`
}

func (r *RoleItem) GetUserRole() []RoleItem {
	var list []RoleItem
	o := orm.NewOrm()
	qs := o.QueryTable(r)
	_, err := qs.Filter("RoleId", r.RoleId).All(&list)
	if err != nil {
		return []RoleItem{}
	}
	return list
}

var RoleList = []role{
	{
		Path:  "",
		Index: "",
		Value: "",
		Name:  "激活码管理",
		Child: []role{
			{
				Path:  "/admin/project/createKeys",
				Index: "",
				Value: "",
				Name:  "创建激活码",
			},
			{
				Path:  "/admin/project/lockKey",
				Index: "",
				Value: "",
				Name:  "锁定/解锁激活码",
			},
			{
				Path:  "/admin/project/deleteKeys",
				Index: "",
				Value: "",
				Name:  "删除激活码",
			},
			{
				Path:  "/admin/project/batchKeys",
				Index: "type",
				Value: "lock",
				Name:  "批量锁定/解锁激活码",
			},
			{
				Path:  "/admin/project/batchKeys",
				Index: "type",
				Value: "delete",
				Name:  "批量删除激活码",
			},
		},
	},
	{
		Path:  "",
		Index: "",
		Value: "",
		Name:  "会员管理",
		Child: []role{
			{
				Path:  "/admin/project/lockMember",
				Index: "",
				Value: "",
				Name:  "锁定/解锁会员",
			},
			{
				Path:  "/admin/project/unbindMember",
				Index: "",
				Value: "",
				Name:  "解绑会员",
			},
			{
				Path:  "/admin/project/updateMember",
				Index: "",
				Value: "",
				Name:  "修改会员基础资料",
			},
		},
	},
	{
		Path:  "",
		Index: "",
		Value: "",
		Name:  "最近在线",
		Child: []role{
			{
				Path:  "/admin/project/memberLogout",
				Index: "",
				Value: "",
				Name:  "强制会员下线",
			},
		},
	},
}

func RoleWithArr(roleItem []role, roleList *[]role) []role {
	var list []role
	for _, i := range roleItem {
		list = append(list, i)
		if i.Child != nil {
			list = append(list, RoleWithArr(i.Child, roleList)...)
		}
	}
	*roleList = list
	return list
}
