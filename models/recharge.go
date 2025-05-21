package models

import "time"

type Recharge struct {
	ID         int       `orm:"column(id)" json:"id"`
	Type       int       `orm:"default(0);index;description(0用户激活,1邀请赠送,2注册赠送,3活动赠送,4在线充值,5套餐消费)" json:"type"`
	ProjectId  int       `orm:"index;default(0)" json:"project_id"`
	ManagerId  int       `orm:"index;default(0)" json:"manager_id"`
	CardsId    int       `orm:"index;default(0);description(关联激活码类型ID)" json:"cards_id"`
	CardsTitle string    `orm:"size(28);null;description(关联激活码名称)" json:"title"`
	KeysId     int       `orm:"index;default(0);description(关联激活码ID)" json:"keys_id"`
	LongKeys   string    `orm:"size(32);null;description(激活码卡密)" json:"long_keys"`
	Price      float64   `orm:"digits(12);decimals(2);default(0);description(关联激活码定价)" json:"price"`
	Money      float64   `orm:"digits(12);decimals(2);default(0);description(关联套餐定价)" json:"money"`
	LevelId    int       `orm:"index;default(0);description(套餐VIP级别,0普通)" json:"level_id"`
	LevelTitle string    `orm:"size(28);null;description(套餐名称)" json:"level_title"`
	CoverType  int       `orm:"index;default(0);description(0叠加,1覆盖)" json:"cover_type"`
	OrderId    int       `orm:"index;default(0);description(关联订单ID)" json:"order_id"`
	Days       float64   `orm:"digits(12);decimals(2);default(0);description(关联激活码天数)" json:"days"`
	Points     int       `orm:"default(0);description(关联激活码点数)" json:"points"`
	KeyExtAttr string    `orm:"size(200);null;description(关联激活码扩展属性)" json:"key_ext_attr"`
	Tag        string    `orm:"size(200);null;description(关联激活码标签)" json:"tag"`
	Member     string    `orm:"size(60);null;description(关联用户)" json:"member"`
	MemberId   int       `orm:"default(0);index;description(关联用户ID)" json:"member_id"`
	UseTime    int       `orm:"index;default(0);description(关联使用时间)" json:"use_time"`
	ActiveTime int       `orm:"index;default(0);description(关联激活时间)" json:"active_time"`
	EndTime    int       `orm:"index;default(0);description(关联失效时间)" json:"end_time"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);index" json:"create_time"`
}
