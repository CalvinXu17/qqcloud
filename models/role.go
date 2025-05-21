package models

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"time"
)

type Role struct {
	ID          int       `orm:"column(id)" json:"id"`
	ManagerId   int       `orm:"index;default(0);description(归属开发者)" json:"manager_id"`
	Title       string    `orm:"null;description(角色名称);size(20)" json:"title"`
	Description string    `orm:"null;description(角色功能描述);size(300)" json:"description"`
	CreateTime  time.Time `orm:"auto_now_add;type(datetime);index" json:"create_time"`
}

func GetAllRole() []role {
	var roleList []role
	RoleWithArr(RoleList, &roleList)
	return roleList
}

func (r *Role) GetRoleList(pageSize int64, page int64, managerId int) (status bool, pager Pager) {
	var data []Role
	o := orm.NewOrm()
	count, err := o.QueryTable(&r).Filter("ManagerId", managerId).Count()
	if err != nil {
		logs.Error(err)
		return false, Pager{}
	}
	totalPage, offset, currentPage := PageCount(count, pageSize, page)
	_, err = o.QueryTable(&r).Filter("ManagerId", managerId).Limit(pageSize, offset).All(&data)
	return true, Pager{
		Count:       count,
		CurrentPage: currentPage,
		Data:        data,
		PageSize:    pageSize,
		TotalPages:  totalPage,
	}
}

func (r *Role) GetRoleAll(managerId int) []Role {
	var data []Role
	o := orm.NewOrm()
	_, err := o.QueryTable(&r).Filter("ManagerId", managerId).All(&data)
	if err != nil {
		return data
	}
	return data
}

func (r *Role) Add(name []string) int64 {
	var list []RoleItem
	o := orm.NewOrm()
	id, err := o.Insert(r)
	if err != nil {
		return 0
	}
	var roleList []role
	RoleWithArr(RoleList, &roleList)
	for _, i := range roleList {
		for _, x := range name {
			if i.Name == x {
				list = append(list, RoleItem{
					RoleId: int(id),
					Path:   i.Path,
					Index:  i.Index,
					Value:  i.Value,
					Name:   i.Name,
				})
			}
		}
	}
	if len(list) > 0 {
		_, _ = o.InsertMulti(1, list)
	}
	return id
}

func (r *Role) Update(name []string) int64 {

	var list []RoleItem
	o := orm.NewOrm()
	old := Role{ID: r.ID}
	err := o.Read(&old)
	if err != nil {
		return 0
	}
	old.Title = r.Title
	old.Description = r.Description
	var roleList []role
	RoleWithArr(RoleList, &roleList)
	for _, i := range roleList {
		for _, x := range name {
			if i.Name == x {
				list = append(list, RoleItem{
					RoleId: r.ID,
					Path:   i.Path,
					Index:  i.Index,
					Value:  i.Value,
					Name:   i.Name,
				})
			}
		}
	}

	row, err := o.Update(&old)
	if err != nil {
		return 0
	}
	qs := o.QueryTable("RoleItem")
	_, _ = qs.Filter("RoleId", r.ID).Delete()
	if len(list) > 0 {
		_, _ = o.InsertMulti(1, list)
	}
	return row
}

func (r *Role) Delete() int64 {
	o := orm.NewOrm()
	row, err := o.Delete(r)
	if err != nil {
		return 0
	}
	qs := o.QueryTable("RoleItem")
	_, _ = qs.Filter("RoleId", r.ID).Delete()
	qs = o.QueryTable("Manager")
	_, _ = qs.Filter("PowerId", r.ID).Update(orm.Params{"PowerId": 0})
	return row
}
