package models

import (
	"time"

	"github.com/beego/beego/v2/adapter/orm"
	"github.com/beego/beego/v2/core/logs"
)

type Cards struct {
	ID         int       `orm:"column(id)" json:"id"`
	ProjectId  int       `orm:"index;default(0)" json:"project_id"`
	ManagerId  int       `orm:"index;default(0)" json:"manager_id"`
	Title      string    `orm:"size(28);null;description(激活码名称)" json:"title" valid:"Required,AlphaNumeric" valid:"MaxSize(28)"`
	Price      float64   `orm:"digits(12);decimals(2);default(0);description(激活码定价)" json:"price" valid:"Max(99999);MaxSize(6)"`
	KeyPrefix  string    `orm:"index;size(4);null;description(激活码前缀)" json:"key_prefix"`
	LevelId    int       `orm:"index;default(0);description(关联会员套餐)" json:"level_id"`
	Days       float64   `orm:"digits(12);decimals(2);default(0);description(天数)" json:"days" valid:"Max(99999);MaxSize(5)"`
	Points     int       `orm:"default(0);description(点数)" json:"points" valid:"Max(9999);MaxSize(4)"`
	KeyExtAttr string    `orm:"size(200);null;description(激活码扩展属性)" json:"key_ext_attr" valid:"MaxSize(199)"`
	Tag        string    `orm:"size(200);null;description(标签)" json:"tag" valid:"MaxSize(199)"`
	IsLock     int       `orm:"default(0);description(0正常,1锁定);index" json:"is_lock"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);index" json:"create_time"`
}

func (c *Cards) GetById() (bool, Cards) {
	o := orm.NewOrm()
	cards := Cards{ID: c.ID}
	err := o.Read(&cards)
	if err != nil {
		return false, cards
	}
	return true, cards
}

func (c *Cards) Add(ProjectId int, Title string, Price float64, KeyPrefix string, LevelId int, Days float64, Points int, KeyExtAttr string, Tag string, IsLock int, managerId int) int64 {

	o := orm.NewOrm()
	if ProjectId != 0 {
		project := Project{ID: ProjectId}
		err := o.Read(&project)
		if err != nil {
			return 0
		}
	}

	p := Cards{
		ProjectId:  ProjectId,
		Title:      Title,
		Price:      Price,
		KeyPrefix:  KeyPrefix,
		LevelId:    LevelId,
		Days:       Days,
		Points:     Points,
		KeyExtAttr: KeyExtAttr,
		Tag:        Tag,
		IsLock:     IsLock,
		ManagerId:  managerId,
	}
	id, err := o.Insert(&p)
	if err != nil {
		return 0
	}
	return id
}

func (c *Cards) Update(ID int, Title string, Price float64, KeyPrefix string, LevelId int, Days float64, Points int, KeyExtAttr string, Tag string, IsLock int) bool {
	o := orm.NewOrm()
	card := Cards{ID: ID}
	err := o.Read(&card)
	if err != nil {
		return false
	}
	if Title != "" {
		card.Title = Title
	}
	card.Price = Price
	card.KeyPrefix = KeyPrefix
	card.LevelId = LevelId
	card.Days = Days
	card.Points = Points
	if KeyExtAttr != "" {
		card.KeyExtAttr = KeyExtAttr
	}
	if Tag != "" {
		card.Tag = Tag
	}
	card.IsLock = IsLock
	row, err := o.Update(&card)
	if row > 0 {
		return true
	}
	return false
}

func (p *Cards) GetCardList(projectId int, pageSize int64, page int64) (status bool, pager Pager) {
	var data []Cards
	o := orm.NewOrm()
	var count int64
	var err error
	if projectId == 0 {
		count, err = o.QueryTable(&p).Filter("ProjectId", projectId).Count()
	} else {
		count, err = o.QueryTable(&p).Count()
	}
	if err != nil {
		logs.Error(err)
		return false, Pager{}
	}
	totalPage, offset, currentPage := PageCount(count, pageSize, page)
	if projectId == 0 {
		_, err = o.QueryTable(&p).Limit(pageSize, offset).All(&data)
	} else {
		_, err = o.QueryTable(&p).Filter("ProjectId", projectId).Limit(pageSize, offset).All(&data)
	}
	return true, Pager{
		Count:       count,
		CurrentPage: currentPage,
		Data:        data,
		PageSize:    pageSize,
		TotalPages:  totalPage,
	}
}

func (c *Cards) CardList() (status bool, pager Pager) {
	var data []Cards
	o := orm.NewOrm()
	var err error
	_, err = o.QueryTable(&c).All(&data)
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

func (c *Cards) GetAgentCardList(list []int) (status bool, pager Pager) {
	var data []Cards
	o := orm.NewOrm()
	var err error
	_, err = o.QueryTable(&c).Filter("ID__in", list).All(&data)
	if err != nil {
		logs.Error("获取代理授权卡类失败", err)
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

func (c *Cards) Delete(id int) bool {
	p := Cards{
		ID: id,
	}
	o := orm.NewOrm()
	err := o.Read(&p)
	if err != nil {
		logs.Error("查询激活码类型失败")
		return false
	}
	rows, err := o.Delete(&p)
	if rows > 0 {
		qs := o.QueryTable("Keys")
		row, _ := qs.Filter("CardsId", p.ID).Delete()
		if row > 0 {
			logs.Info("删除激活码数据：", row)
		}
		qs = o.QueryTable("ManagerCards")
		row, _ = qs.Filter("CardsId", p.ID).Delete()
		if row > 0 {
			logs.Info("删除授权激活码类型数据：", row)
		}
		return true
	}
	return false
}
