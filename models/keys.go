package models

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"strings"
	"time"
	"verification/controllers/common"
)

type Keys struct {
	ID         int       `orm:"column(id)" json:"id"`
	ProjectId  int       `orm:"index;default(0)" json:"project_id"`
	ManagerId  int       `orm:"index;default(0)" json:"manager_id"`
	LongKeys   string    `orm:"index;size(32);null;description(激活码卡密);" json:"long_keys"`
	CardsId    int       `orm:"index;default(0);description(关联激活码类型ID)" json:"cards_id"`
	LevelId    int       `orm:"index;default(0);description(套餐VIP级别,0普通)" json:"level_id"`
	OrderId    int       `orm:"index;default(0);description(关联激活码订单ID)" json:"order_id"`
	Member     string    `orm:"size(60);null;description(关联用户)" json:"member"`
	MemberId   int       `orm:"default(0);index;description(关联用户ID)" json:"member_id"`
	UseTime    int       `orm:"index;default(0);description(使用时间)" json:"use_time"`
	Tag        string    `orm:"size(200);null;description(标签)" json:"tag"`
	IsLock     int       `orm:"default(0);description(0正常,1锁定);index" json:"is_lock"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);index" json:"create_time"`
}
type KeysResult struct {
	Keys         []Keys  `json:"keys"`
	Count        int     `json:"count"`
	Price        float64 `json:"price"`
	ManagerMoney float64 `json:"manager_money"`
}

func (k *Keys) QueryCreate(cardsId int, count int64, managerId int, agentPrice float64) (status bool, cost float64, msg string) {
	o := orm.NewOrm()
	c := Cards{ID: cardsId}
	err := o.Read(&c)
	if err != nil {
		return false, 0, "激活码类型不存在"
	}
	var price float64
	if agentPrice > 0 {
		price = agentPrice
	} else {
		price = c.Price
	}
	u := Manager{ID: managerId}
	err = o.Read(&u)
	if err != nil {
		return false, 0, "管理账号不存在"
	}
	if u.Level == 0 {
		price = 0
	}
	countPrice := float64(count) * price
	if countPrice > u.Money {
		return false, countPrice, "账号余额不足"
	}
	return true, countPrice, "可创建"
}

func (k *Keys) Create(cardsId int, count int64, length int, createType int, managerId int, projectId int, tag string, agentPrice float64) (status bool, result KeysResult, msg string) {
	_, ac := common.GetCacheAC()
	var res = KeysResult{}
	var successCount = 0
	length = length - 4
	o := orm.NewOrm()
	to, err := o.Begin()
	if err != nil {
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, res, "请稍后尝试"
	}
	c := Cards{ID: cardsId}
	err = to.Read(&c)
	if err != nil {
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, res, "激活码类型不存在"
	}
	if tag == "" {
		tag = c.Tag
	}
	length = length + (4 - len(c.KeyPrefix))
	var price float64
	if agentPrice > 0 {
		price = agentPrice
	} else {
		price = c.Price
	}
	u := Manager{ID: managerId}
	err = to.Read(&u)
	if err != nil {
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, res, "管理账号不存在"
	}
	if u.Level == 0 {
		price = 0
	}
	countPrice := float64(count) * price
	if countPrice > u.Money {
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, res, "账号余额不足"
	}
	order := Order{Count: int(count), Price: countPrice, CardsId: cardsId, ManagerId: managerId, ProjectId: projectId}
	orderId, err := to.Insert(&order)
	if err != nil {
		return false, res, "订单创建失败"
	}
	keyList := make([]Keys, 0)
	keyStrList := make([]string, 0)
	for {
		var longKey strings.Builder
		switch createType {
		// 纯数字组合
		case 0:
			longKey.WriteString(c.KeyPrefix)
			longKey.WriteString(RandNumStr(length))
		// 大写字母数字组合
		case 1:
			longKey.WriteString(strings.ToUpper(c.KeyPrefix))
			longKey.WriteString(RandUpperStr(length))
		// 小写字母数字组合
		case 2:
			longKey.WriteString(strings.ToLower(c.KeyPrefix))
			longKey.WriteString(RandLowerStr(length))
		// 随机字母数字组合
		case 3:
			longKey.WriteString(c.KeyPrefix)
			longKey.WriteString(RandStr(length))
		}
		spliceKeys := Keys{
			ProjectId: projectId,
			ManagerId: managerId,
			LongKeys:  longKey.String(),
			CardsId:   cardsId,
			LevelId:   0,
			OrderId:   int(orderId),
			Tag:       tag,
		}
		logs.Error("keys:", spliceKeys.LongKeys)
		cache := ac.Get(spliceKeys.LongKeys)
		if cache == "" && In(spliceKeys.LongKeys, keyStrList) != true {
			keyList = append(keyList, spliceKeys)
			count = count - 1
			successCount = successCount + 1
			keyStrList = append(keyStrList, spliceKeys.LongKeys)
		}
		if count == 0 {
			break
		}
	}
	if price > 0 {
		u.Money = u.Money - countPrice
		_, err := to.Update(&u)
		if err != nil {
			logs.Error("更新代理余额失败", err)
			err = to.Rollback()
			if err != nil {
				logs.Error("事务回滚失败", err)
			}
			return false, res, "更新代理余额失败"
		}
	}
	_, err = to.InsertMulti(1, &keyList)
	if err != nil {
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, res, "创建激活码失败"
	}
	err = to.Commit()
	if err != nil {
		logs.Error("提交事务失败", err)
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, res, "事务提交失败"
	}
	for _, i := range keyList {
		_ = ac.Put(i.LongKeys, i.LongKeys, 5*365*24*60*60*time.Second)
	}
	return true, KeysResult{Count: successCount, Price: price, Keys: keyList, ManagerMoney: u.Money}, "创建成功"
}

func (p *Keys) GetKeysList(managerIdArr []int, projectId int, pageSize int64, page int64, longKeys string, cardsId int, isActive int, isLock int, member string, orderId int) (status bool, pager Pager) {
	var data []Keys
	o := orm.NewOrm()
	var count int64
	var err error
	qs := o.QueryTable(&p)
	if longKeys != "" {
		qs = qs.Filter("LongKeys", longKeys)
	}
	if cardsId != 0 {
		qs = qs.Filter("CardsId", cardsId)
	}
	if isLock == 0 {
		qs = qs.Filter("IsLock", 0)
	}
	if isLock == 1 {
		qs = qs.Filter("IsLock__gt", 0)
	}
	if orderId > 0 {
		qs = qs.Filter("OrderId", orderId)
	}
	if member != "" {
		qs = qs.Filter("Member", member)
	}
	if isActive == 0 {
		qs = qs.Filter("UseTime", 0)
	}
	if isActive == 1 {
		qs = qs.Filter("UseTime__gt", 0)
	}
	qs = qs.Filter("ManagerId__in", managerIdArr)

	if projectId == -1 {
		count, err = qs.Count()
	} else {
		count, err = qs.Filter("ProjectId", projectId).Count()
	}

	if err != nil {
		logs.Error(err)
		return false, Pager{}
	}
	qs = o.QueryTable(&p)
	if longKeys != "" {
		qs = qs.Filter("LongKeys", longKeys)
	}
	if cardsId != 0 {
		qs = qs.Filter("CardsId", cardsId)
	}
	if isLock == 0 {
		qs = qs.Filter("IsLock", 0)
	}
	if isLock == 1 {
		qs = qs.Filter("IsLock__gt", 0)
	}
	if orderId > 0 {
		qs = qs.Filter("OrderId", orderId)
	}
	if member != "" {
		qs = qs.Filter("Member", member)
	}
	if isActive == 0 {
		qs = qs.Filter("UseTime", 0)
	}
	if isActive == 1 {
		qs = qs.Filter("UseTime__gt", 0)
	}
	qs = qs.Filter("ManagerId__in", managerIdArr)
	totalPage, offset, currentPage := PageCount(count, pageSize, page)
	if projectId == -1 {
		_, err = qs.Limit(pageSize, offset).All(&data)
	} else {
		_, err = qs.Filter("ProjectId", projectId).Limit(pageSize, offset).All(&data)
	}
	return true, Pager{
		Count:       count,
		CurrentPage: currentPage,
		Data:        data,
		PageSize:    pageSize,
		TotalPages:  totalPage,
	}
}

func (p *Keys) Lock(id int) (status bool, msg string) {
	o := orm.NewOrm()
	k := Keys{ID: id}
	err := o.Read(&k)
	if err != nil {
		return false, "激活码不存在"
	}
	if k.IsLock == 0 {
		k.IsLock = 1
		msg = "锁定成功"
	} else {
		k.IsLock = 0
		msg = "解锁成功"
	}
	row, err := o.Update(&k)
	if row > 0 {
		return true, msg
	}
	return false, "操作失败"
}

func (p *Keys) Delete(id int) bool {
	o := orm.NewOrm()
	k := Keys{ID: id}
	err := o.Read(&k)
	if err != nil {
		return false
	}
	row, err := o.Delete(&k)
	if row > 0 {
		return true
	}
	return false
}

func (p *Keys) GetOneByKeys() (bool, Cards) {
	c := Cards{}
	if p.LongKeys == "" {
		return false, c
	}
	o := orm.NewOrm()
	err := o.Read(p, "LongKeys")
	if err != nil {
		return false, c
	}
	c.ID = p.CardsId
	err = o.Read(&c)
	if err != nil {
		return false, c
	}
	return true, c
}

func (k *Keys) Update() bool {
	o := orm.NewOrm()
	row, err := o.Update(k)
	if err != nil {
		return false
	}
	if row > 0 {
		return true
	}
	return false
}

type RechargeResult struct {
	AddCdays    float64 `json:"addCdays"`
	AddPoints   int     `json:"addPoints"`
	KeyCdays    float64 `json:"keyCdays"`
	KeyPoints   int     `json:"keyPoints"`
	CountCdays  float64 `json:"countCdays"`
	CountPoints int     `json:"countPoints"`
}
type RechargeParam struct {
	User string `json:"user"`
	Key  string `json:"key"`
}

func (this *Keys) Recharge(user string, key string, p Project) (bool, string, RechargeResult) {
	result := RechargeResult{}
	o := orm.NewOrm()
	to, err := o.Begin()
	if err != nil {
		logs.Error("事务开启失败", err)
		return false, "充值失败，请稍后再试", result
	}
	u := Member{}
	isEmail := common.VerifyEmailFormat(user)
	if isEmail == true {
		u.Email = user
		err = o.Read(&u, "Email")
	} else {
		u.Name = user
		err = o.Read(&u, "Name")
	}

	if err != nil {
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, "会员不存在", result
	}
	k := Keys{LongKeys: key}
	err = to.Read(&k, "LongKeys")
	if err != nil {
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, "激活码无效", result
	}
	if k.ProjectId != 0 && k.ProjectId != p.ID {
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, "激活码无效", result
	}
	if k.UseTime > 0 {
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, "激活码已被使用", result
	}
	c := Cards{
		ID: k.CardsId,
	}
	err = to.Read(&c)
	if err != nil {
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, "激活码类型不存在", result
	}
	days := u.Days
	activeTime := u.ActiveTime
	endTime := u.EndTime

	if u.EndTime < int(time.Now().Unix()) {
		activeTime = int(time.Now().Unix())
		days = c.Days
		endTime = activeTime + int(c.Days*24*60*60)
	} else {
		endTime = endTime + int(c.Days*24*60*60)
		days = c.Days + u.Days
	}
	if endTime > 2147483647 {
		endTime = 2147483647
	}

	u.ActiveTime = activeTime
	u.EndTime = endTime
	u.Days = days
	u.Points = u.Points + c.Points
	u.ManagerId = k.ManagerId
	if k.Tag != "" {
		u.Tag = k.Tag
	}
	if c.Tag != "" && k.Tag == "" {
		u.Tag = c.Tag
	}
	if c.KeyExtAttr != "" {
		u.KeyExtAttr = c.KeyExtAttr
	}
	row, err := to.Update(&u)
	if err != nil {
		logs.Error("充值失败：", err)
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, "充值失败", result
	}
	if row == 0 {
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, "充值失败", result
	}
	k.UseTime = int(time.Now().Unix())
	k.Member = u.Name
	k.MemberId = u.ID
	row, err = to.Update(&k)
	if err != nil {
		logs.Error("更新激活码失败", err)
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, "充值失败,请稍后再试", result
	}
	if row == 0 {
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, "充值失败", result
	}
	err = to.Commit()
	if err != nil {
		logs.Error("提交事务失败", err)
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, "充值失败，请稍后再试", result
	}
	result.AddCdays = c.Days
	result.AddPoints = c.Points
	result.KeyCdays = c.Days
	result.KeyPoints = c.Points
	result.CountPoints = u.Points
	result.CountCdays = u.Days
	return true, "充值成功", result
}

func (p *Keys) GetCount(managerIdArr []int) int64 {
	o := orm.NewOrm()
	qs := o.QueryTable("Keys")
	count, err := qs.Filter("ManagerId__in", managerIdArr).Count()
	if err != nil {
		logs.Error("查询激活码总数错误：", err)
		return 0
	}
	return count
}

func (p *Keys) GetAddRangeCount(dateRange []common.DateObj, projectId int, managerIdArr []int) []int64 {
	var countArray []int64
	o := orm.NewOrm()
	qs := o.QueryTable(&p)
	var count int64
	var err error
	for _, i := range dateRange {
		if projectId == 0 {
			count, err = qs.Filter("ManagerId__in", managerIdArr).Filter("CreateTime__gte", i.StartDate).Filter("CreateTime__lte", i.EndDate).Count()
		} else {
			count, err = qs.Filter("ManagerId__in", managerIdArr).Filter("CreateTime__gte", i.StartDate).Filter("ProjectId", projectId).Filter("CreateTime__lte", i.EndDate).Count()
		}
		//err = o.Raw("select count(*) as count_num from member where create_time between datetime(?) and  datetime(?)", i.StartDate, i.EndDate).QueryRow(&c)
		//logs.Info("count:", c, i.StartDate, i.EndDate)
		if err != nil {
			logs.Error("查询总数错误：", err)
		}
		countArray = append(countArray, count)
	}
	return countArray
}

func (p *Keys) GetActiveRangeCount(dateRange []common.DateObj, projectId int, managerIdArr []int) []int64 {
	var countArray []int64
	o := orm.NewOrm()
	qs := o.QueryTable(&p)
	var count int64
	var err error
	for _, i := range dateRange {
		if projectId == 0 {
			count, err = qs.Filter("ManagerId__in", managerIdArr).Filter("UseTime__gte", i.Start).Filter("UseTime__lte", i.End).Count()
		} else {
			count, err = qs.Filter("ManagerId__in", managerIdArr).Filter("UseTime__gte", i.Start).Filter("ProjectId", projectId).Filter("UseTime__lte", i.End).Count()
		}
		//err = o.Raw("select count(*) as count_num from member where create_time between datetime(?) and  datetime(?)", i.StartDate, i.EndDate).QueryRow(&c)
		//logs.Info("count:", c, i.StartDate, i.EndDate)
		if err != nil {
			logs.Error("查询总数错误：", err)
		}
		countArray = append(countArray, count)
	}
	return countArray
}

func (p *Keys) LockSelect(lock int, id []int, managerIdArr []int) int64 {
	o := orm.NewOrm()
	qs := o.QueryTable(&p)
	rows, err := qs.Filter("Id__in", id).Filter("ManagerId__in", managerIdArr).Update(orm.Params{"IsLock": lock})
	if err != nil {
		return 0
	}
	return rows
}

func (p *Keys) DeleteSelect(id []int, managerIdArr []int) int64 {
	o := orm.NewOrm()
	qs := o.QueryTable(&p)
	rows, err := qs.Filter("Id__in", id).Filter("ManagerId__in", managerIdArr).Delete()
	if err != nil {
		return 0
	}
	return rows
}

func (p *Keys) LockQuery(lock int, managerIdArr []int, projectId int, longKeys string, cardsId int, isActive int, isLock int, member string, orderId int) int64 {
	var count int64
	var err error
	o := orm.NewOrm()
	qs := o.QueryTable(&p)
	if longKeys != "" {
		qs = qs.Filter("LongKeys", longKeys)
	}
	if cardsId != 0 {
		qs = qs.Filter("CardsId", cardsId)
	}
	if isLock == 0 {
		qs = qs.Filter("IsLock", 0)
	}
	if isLock == 1 {
		qs = qs.Filter("IsLock__gt", 0)
	}
	if orderId > 0 {
		qs = qs.Filter("OrderId", orderId)
	}
	if member != "" {
		qs = qs.Filter("Member", member)
	}
	if isActive == 0 {
		qs = qs.Filter("UseTime", 0)
	}
	if isActive == 1 {
		qs = qs.Filter("UseTime__gt", 0)
	}
	qs = qs.Filter("ManagerId__in", managerIdArr)
	if projectId == -1 {
		count, err = qs.Update(orm.Params{"IsLock": lock})
	} else {
		count, err = qs.Filter("ProjectId", projectId).Update(orm.Params{"IsLock": lock})
	}
	if err != nil {
		return 0
	}
	return count
}

func (p *Keys) DeleteQuery(managerIdArr []int, projectId int, longKeys string, cardsId int, isActive int, isLock int, member string, orderId int) int64 {
	var count int64
	var err error
	o := orm.NewOrm()
	qs := o.QueryTable(&p)
	if longKeys != "" {
		qs = qs.Filter("LongKeys", longKeys)
	}
	if cardsId != 0 {
		qs = qs.Filter("CardsId", cardsId)
	}
	if isLock == 0 {
		qs = qs.Filter("IsLock", 0)
	}
	if isLock == 1 {
		qs = qs.Filter("IsLock__gt", 0)
	}
	if orderId > 0 {
		qs = qs.Filter("OrderId", orderId)
	}
	if member != "" {
		qs = qs.Filter("Member", member)
	}
	if isActive == 0 {
		qs = qs.Filter("UseTime", 0)
	}
	if isActive == 1 {
		qs = qs.Filter("UseTime__gt", 0)
	}
	qs = qs.Filter("ManagerId__in", managerIdArr)
	if projectId == -1 {
		count, err = qs.Delete()
	} else {
		count, err = qs.Filter("ProjectId", projectId).Delete()
	}
	if err != nil {
		return 0
	}
	return count
}
