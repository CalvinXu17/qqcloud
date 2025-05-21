package models

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"time"
)

type Order struct {
	ID         int       `orm:"column(id)" json:"id"`
	ProjectId  int       `orm:"index;default(0)" json:"project_id"`
	ManagerId  int       `orm:"index;default(0)" json:"manager_id"`
	CardsId    int       `orm:"index;default(0);description(关联激活码类型ID)" json:"cards_id"`
	Count      int       `orm:"default(0);description(激活码个数)" json:"count"`
	Price      float64   `orm:"digits(12);decimals(2);default(0);description(总金额)" json:"price"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime)" json:"create_time"`
}

func (order *Order) GetOrderList(pageSize int64, page int64, managerId []int) (bool,Pager) {
	var data []Order
	o := orm.NewOrm()
	count, err := o.QueryTable(&order).Filter("ManagerId__in", managerId).Count()
	if err != nil {
		logs.Error(err)
		return false, Pager{}
	}
	totalPage, offset, currentPage := PageCount(count, pageSize, page)
	_, err = o.QueryTable(&order).Filter("ManagerId__in", managerId).Limit(pageSize, offset).All(&data)
	return true, Pager{
		Count:       count,
		CurrentPage: currentPage,
		Data:        data,
		PageSize:    pageSize,
		TotalPages:  totalPage,
	}
}
