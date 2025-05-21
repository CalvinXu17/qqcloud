package models

import (
	"github.com/beego/beego/v2/adapter/orm"
	"github.com/beego/beego/v2/core/logs"
	"time"
)

type ManagerCards struct {
	ID         int       `orm:"column(id)" json:"id"`
	ProjectId  int       `orm:"index;default(0)" json:"project_id"`
	ManagerId  int       `orm:"index;default(0)" json:"manager_id"`
	CardsId    int       `orm:"index;default(0);description(关联激活码类型ID)" json:"cards_id"`
	Price      float64   `orm:"digits(12);decimals(2);default(0);description(激活码定价)" json:"price" valid:"Max(9999);MaxSize(4)"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);index" json:"create_time"`
}

func (m *ManagerCards) Update() (bool, string) {
	o := orm.NewOrm()
	price := m.Price
	err := o.QueryTable("ManagerCards").Filter("ManagerId", m.ManagerId).Filter("ID", m.ID).One(m)
	if err != nil {
		return false, "授权激活码类型不存在"
	}
	m.Price = price
	row, err := o.Update(m)
	if row > 0 {
		return true, "更新成功"
	}
	return false, "更新失败"
}

func (m *ManagerCards) GetAgentProjectId(managerId int) []int {
	var data []ManagerCards
	var list []int
	o := orm.NewOrm()
	var err error
	_, err = o.QueryTable(&m).Filter("ManagerId", managerId).All(&data)
	if err != nil {
		logs.Error(err)
		return list
	}
	for _, i := range data {
		list = append(list, i.ProjectId)
	}
	return list
}

func (m *ManagerCards) GetCardList(managerId int) (bool, []ManagerCards) {
	var data []ManagerCards
	o := orm.NewOrm()
	var err error
	_, err = o.QueryTable(&m).Filter("ManagerId", managerId).All(&data)
	if err != nil {
		logs.Error(err)
		return false, data
	}
	return true, data
}

func (m *ManagerCards) GetGroupCardList(managerId int) (bool, []ManagerCards) {
	var data []ManagerCards
	o := orm.NewOrm()
	var err error
	_, err = o.QueryTable(&m).Filter("ManagerId", managerId).GroupBy("ProjectId").All(&data)
	if err != nil {
		logs.Error(err)
		return false, data
	}
	return true, data
}

func (m *ManagerCards) CardList(managerId int) (status bool, pager Pager) {
	var data []ManagerCards
	o := orm.NewOrm()
	var err error
	_, err = o.QueryTable(&m).Filter("ManagerId", managerId).All(&data)
	if err != nil {
		logs.Error(err)
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

func (m *ManagerCards) Add() (int64, string) {
	o := orm.NewOrm()
	row, err := o.QueryTable("ManagerCards").Filter("CardsId", m.CardsId).Filter("ManagerId", m.ManagerId).Count()
	if row > 0 {
		return 0, "已添加该授权"
	}
	id, err := o.Insert(m)
	if err != nil {
		return 0, "添加授权失败"
	}
	return id, "添加授权成功"
}

func (m *ManagerCards) Delete() (int64, string) {
	o := orm.NewOrm()
	row, err := o.QueryTable("ManagerCards").Filter("ID", m.ID).Filter("ManagerId", m.ManagerId).Delete()
	if err != nil {
		return 0, "移除失败"
	}
	if row > 0 {
		return row, "授权移除成功"
	}
	return 0, "移除失败"
}
