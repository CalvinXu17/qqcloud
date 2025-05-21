package models

import "time"

type Client struct {
	ID         int       `orm:"column(id)" json:"id"`
	Type       int       `orm:"default(0);index;description(0登录,1解绑,2心跳,3下线)" json:"type"`
	ProjectId  int       `orm:"index;default(0)" json:"project_id"`
	ManagerId  int       `orm:"index;default(0)" json:"manager_id"`
	Member     string    `orm:"size(60);null;description(关联用户)" json:"user"`
	MemberId   int       `orm:"index;default(0);description(关联用户ID)" json:"member_id"`
	Token      string    `orm:"size(40);description(关联特征码)" json:"token"`
	Mac        string    `orm:"size(32);index;description(机器码)" json:"mac"`
	Ip         string    `orm:"size(32);description(ip地址)" json:"ip"`
	Location   string    `orm:"size(32);description(登录IP地理位置)" json:"location"`
	LocationId int       `orm:"default(0);index;description(登录地址关联ID)" json:"location_id"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime)" json:"create_time"`
}
