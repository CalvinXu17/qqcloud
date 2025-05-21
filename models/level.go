package models

import "time"

type Level struct {
	ID         int       `orm:"column(id)" json:"id"`
	ProjectId  int       `orm:"index;default(0)" json:"project_id"`
	ManagerId  int       `orm:"index;default(0)" json:"manager_id"`
	Title      string    `orm:"size(28);null;description(套餐名称)" json:"title"`
	Price      float64   `orm:"digits(12);decimals(2);default(0);description(套餐定价)" json:"price"`
	Level      int       `orm:"index;default(0);description(套餐级别)" json:"level"`
	CoverType  int       `orm:"index;default(0);description(0叠加,1覆盖)" json:"cover_type"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);index" json:"create_time"`
}
