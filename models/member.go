package models

import (
	"strconv"
	"time"
	"verification/controllers/common"
	"verification/validation/api"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type Member struct {
	ID            int       `orm:"column(id)" json:"id"`
	ManagerId     int       `orm:"index;default(0)" json:"manager_id"`
	ProjectId     int       `orm:"index;default(0)" json:"project_id"`
	ParentId      int       `orm:"index;default(0)" json:"parent_id"`
	Email         string    `orm:"index;size(40);null" json:"email" valid:"MaxSize(40);"`
	Name          string    `orm:"index;size(32)" json:"name" valid:"MaxSize(32)"`
	NickName      string    `orm:"index;size(32);null" json:"nick_name" valid:"MaxSize(32)"`
	Password      string    `orm:"size(32);null" json:"password" valid:"MaxSize(32)"`
	SafePassword  string    `orm:"size(32);null" json:"safe_password" valid:"MaxSize(32)"`
	Money         float64   `orm:"digits(12);decimals(2);default(0)" json:"money" valid:"MinSize(0);MaxSize(8)"`
	Days          float64   `orm:"digits(12);decimals(2);default(0)" json:"days" valid:"MaxSize(6)"`
	Points        int       `orm:"default(0)" json:"points" valid:"MaxSize(5)"`
	Mac           string    `orm:"null;size(32);" json:"mac" valid:"MaxSize(32)"`
	Phone         int64     `orm:"default(0)" json:"phone" valid:"MaxSize(11)"`
	Scope         int       `orm:"default(0)" json:"scope" valid:"MaxSize(6)"`
	KeyExtAttr    string    `orm:"size(200);null" json:"key_ext_attr" valid:"MaxSize(200)"`
	Tag           string    `orm:"size(200);null" json:"tag" valid:"MaxSize(200)"`
	ActiveTime    int64     `orm:"index;default(0)" json:"active_time"`
	EndTime       int64     `orm:"index;default(0)" json:"end_time"`
	PauseTime     int64     `orm:"index;default(0)" json:"pause_time"`
	LastLoginTime int64     `orm:"index;default(0)" json:"last_login_time"`
	LastLoginIp   string    `orm:"index;size(30);null" json:"last_login_ip"`
	IsLock        int       `orm:"default(0);description(0正常,1锁定);index" json:"is_lock"`
	CreateTime    time.Time `orm:"auto_now_add;type(datetime);index" json:"create_time"`
}

func (m *Member) GetMemberList(managerIdArr []int, projectId int, pageSize int64, page int64, mac string, member string, isLock int) (status bool, pager Pager) {
	var data []Member
	o := orm.NewOrm()
	var count int64
	var err error
	qs := o.QueryTable(&m)
	if mac != "" {
		qs = qs.Filter("Mac", mac)
	}
	if member != "" {
		qs = qs.Filter("Name", member)
	}
	if isLock > -1 {
		if isLock == 0 {
			qs = qs.Filter("IsLock", 0)
		} else {
			qs = qs.Filter("IsLock", 1)
		}
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
	qs = o.QueryTable(&m)
	if mac != "" {
		qs = qs.Filter("Mac", mac)
	}
	if member != "" {
		qs = qs.Filter("Name", member)
	}
	logs.Error("isLock", isLock)
	if isLock > -1 {
		if isLock == 0 {
			qs = qs.Filter("IsLock", 0)
		} else {
			qs = qs.Filter("IsLock", 1)
		}
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

func (m *Member) Add() bool {
	o := orm.NewOrm()
	id, err := o.Insert(m)
	if err != nil {
		logs.Error("会员写入失败", err)
		return false
	}
	if id > 0 {
		return true
	}
	return false
}

func (m *Member) CheckMember() (bool, string) {
	password := m.Password
	o := orm.NewOrm()
	var err error
	if m.Email != "" {
		if m.Email == "" {
			return false, "请填写邮箱账号"
		}
		err = o.Read(m, "Email")
	}
	if m.Email == "" {
		if m.Name == "" {
			return false, "请填写会员账号"
		}
		err = o.Read(m, "Name")
	}
	if err != nil {
		return false, "会员不存在"
	}
	if password != "" && m.Password != password {
		return false, "密码错误"
	}
	return true, "检查成功"
}

func (m *Member) ForgetCheckMember() (bool, string) {
	safePassword := m.SafePassword
	o := orm.NewOrm()
	var err error
	if m.Email != "" {
		err = o.Read(m, "Email")
	}
	if m.Email == "" {
		err = o.Read(m, "Name")
	}
	if err != nil {
		return false, "会员不存在"
	}
	if safePassword != "" && m.SafePassword != safePassword {
		return false, "安全密码错误"
	}
	return true, "检查成功"
}

func (m *Member) UpdateInfo() (bool, string) {
	o := orm.NewOrm()
	row, err := o.Update(m)
	if err != nil {
		logs.Error("更新会员信息失败", err)
	}
	if row == 0 {
		return false, "更新会员信息失败"
	}
	return true, "更新会员信息成功"
}

func (m *Member) KeyRegister(longkeys string, p Project, param api.UnEncrypt, ip string) (bool, string) {
	o := orm.NewOrm()
	to, err := o.Begin()
	if err != nil {
		logs.Error("事务创建失败", err)
		return false, "内部错误"
	}
	if longkeys == "" {
		return false, "激活码不能为空"
	}
	// 免费模式下自动注册单码账户
	if p.StatusType == 2 {
		u := Member{
			ManagerId:     p.ManagerId,
			ProjectId:     p.ID,
			Name:          longkeys,
			NickName:      longkeys,
			Password:      "",
			SafePassword:  "",
			Days:          999,
			Points:        999,
			Mac:           param.Mac,
			KeyExtAttr:    "免费用户",
			Tag:           "免费用户",
			ActiveTime:    time.Now().Unix(),
			EndTime:       time.Now().Unix() + int64(999*(24*60*60)),
			LastLoginTime: time.Now().Unix(),
			LastLoginIp:   ip,
		}
		_, err = to.Insert(&u)
		if err != nil {
			err = to.Rollback()
			if err != nil {
				logs.Error("事务回滚失败", err)
			}
			return false, "会员创建失败"
		}
		err = to.Commit()
		if err != nil {
			logs.Error("提交事务失败", err)
			return false, "激活码注册失败"
		}
		return true, "会员创建成功"
	}
	k := Keys{LongKeys: longkeys}
	err = to.Read(&k, "LongKeys")
	if err != nil {
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, "激活码无效"
	}
	if k.ProjectId != 0 && k.ProjectId != p.ID {
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, "激活码无效"
	}
	if k.UseTime > 0 {
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, "激活码已被使用"
	}
	if k.IsLock > 0 {
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, "激活码已被锁定"
	}
	c := Cards{ID: k.CardsId}
	err = to.Read(&c)
	if err != nil {
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, "激活码类型不存在"
	}
	tag := k.Tag
	if tag == "" {
		tag = c.Tag
	}
	endtime := time.Now().Unix() + int64(c.Days*(24*60*60))
	u := Member{
		ManagerId:     k.ManagerId,
		ProjectId:     p.ID,
		Name:          longkeys,
		NickName:      longkeys,
		Password:      "",
		SafePassword:  "",
		Days:          c.Days,
		Points:        c.Points,
		Mac:           param.Mac,
		KeyExtAttr:    c.KeyExtAttr,
		Tag:           tag,
		ActiveTime:    time.Now().Unix(),
		EndTime:       endtime,
		LastLoginTime: time.Now().Unix(),
		LastLoginIp:   ip,
	}
	id, err := to.Insert(&u)
	if err != nil {
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, "会员创建失败"
	}
	k.UseTime = int(time.Now().Unix())
	k.Member = longkeys
	k.MemberId = int(id)
	_, err = to.Update(&k)
	if err != nil {
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, "更新激活码信息失败"
	}
	err = to.Commit()
	if err != nil {
		logs.Error("提交事务失败", err)
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, "激活码注册失败"
	}
	return true, "激活码注册成功"
}

func (m *Member) Register(k Keys) (bool, string) {
	o := orm.NewOrm()
	to, err := o.Begin()
	if err != nil {
		logs.Error("事务创建失败", err)
	}
	if m.Email != "" {
		err = to.Read(m, "Email")
	}
	if m.Name != "" {
		err = to.Read(m, "Name")
	}
	if m.ID > 0 {
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, "该会员已注册"
	}
	id, err := to.Insert(m)
	if err != nil {
		logs.Error("会员写入失败", err)
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, "注册失败"
	}
	if id == 0 {
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, "注册失败"
	}
	if k.ID > 0 {
		k.Member = m.Name
		k.MemberId = int(id)
		k.UseTime = int(time.Now().Unix())
		kId, err := to.Update(&k)
		if err != nil {
			logs.Error("激活码状态更新失败", err)
			err = to.Rollback()
			if err != nil {
				logs.Error("事务回滚失败", err)
			}
			return false, "注册失败"
		}
		if kId == 0 {
			err = to.Rollback()
			if err != nil {
				logs.Error("事务回滚失败", err)
			}
			return false, "注册失败"
		}
	}

	err = to.Commit()
	if err != nil {
		logs.Error("提交事务失败", err)
		err = to.Rollback()
		if err != nil {
			logs.Error("事务回滚失败", err)
		}
		return false, "注册失败"
	}
	return true, "注册成功"
}

func (u *Member) Unbind(p ProjectLogin, mac string) (status bool, msg string) {
	o := orm.NewOrm()
	isCut := false
	log := MemberLogin{}
	count, _ := log.FetchLog(p.UnbindDate, 2, u.ID)
	switch p.UnbindWeakenMode {
	case 0:
		// 不扣时
		u.Mac = mac
	case 1:
		// 解绑就扣时间
		isCut = true
	case 2:
		// 超出才扣时间
		if int(count) > p.UnbindTimes {
			isCut = true
		}
	case 3:
		// 超出无法解绑
		if int(count) > p.UnbindTimes {
			return false, "已达到最大可解绑次数"
		}
	}
	if p.UnbindWeaken == 0 && p.UnbindWeakenPoints == 0 {
		isCut = false
	}
	if isCut == true {
		// 扣除天数
		if u.Days >= p.UnbindWeaken {
			u.Days = u.Days - p.UnbindWeaken
			u.EndTime = u.EndTime - int64(p.UnbindWeaken*(24*60*60))
		} else {
			return false, "剩余天数不满足解绑扣除需求，无法解绑"
		}
		// 扣除点数
		if u.Points >= p.UnbindWeakenPoints {
			u.Points = u.Points - p.UnbindWeakenPoints
		} else {
			return false, "剩余点数不满足解绑扣除需求，无法解绑"
		}
	}
	u.Mac = mac
	_, err := o.Update(u)
	if err != nil {
		logs.Error("解绑失败：", err)
		return false, "解绑失败"
	}
	return true, "解绑成功"

}

func (m *Member) QueryById() (bool, string) {
	o := orm.NewOrm()
	err := o.Read(m)
	if err != nil {
		return false, "会员不存在"
	}
	return true, "查询成功"
}

func (m *Member) DeleteById() (bool, string) {
	o := orm.NewOrm()
	err := o.Read(m)
	if err != nil {
		return false, "会员不存在"
	}
	row, err := o.Delete(m)
	if err != nil {
		return false, "删除失败"
	}
	if row > 0 {
		_, _ = o.QueryTable("MemberLogin").Filter("MemberId", m.ID).Delete()
		return true, "删除成功"
	}
	return false, "删除失败"
}

func (m *Member) LockById() (bool, string) {
	o := orm.NewOrm()
	err := o.Read(m)
	if err != nil {
		return false, "会员不存在"
	}
	var msg string
	if m.IsLock == 0 {
		m.IsLock = 1
		msg = "锁定成功"
	} else {
		m.IsLock = 0
		msg = "解锁成功"
	}
	row, err := o.Update(m)
	if err != nil {
		return false, "更新失败"
	}
	if row > 0 {
		_, ac := common.GetCacheAC()
		_ = ac.Delete("member-key-" + strconv.Itoa(m.ProjectId) + "-" + m.Name)
		_ = ac.Delete("member-key-" + strconv.Itoa(m.ProjectId) + "-" + m.Email)
		return true, msg
	}
	return false, "修改失败"
}

func (m *Member) UpdateById(member Member) (bool, string) {
	o := orm.NewOrm()
	err := o.Read(m)
	if err != nil {
		return false, "会员不存在"
	}

	if member.Days != m.Days {
		newDays := member.Days - m.Days
		m.EndTime = m.EndTime + int64(newDays*24*60*60)
	}
	if member.EndTime > 0 {
		m.EndTime = member.EndTime
	}
	m.Password = member.Password
	m.SafePassword = member.SafePassword
	m.Mac = member.Mac
	m.Points = member.Points
	m.IsLock = member.IsLock
	m.Email = member.Email
	m.Tag = member.Tag
	m.KeyExtAttr = member.KeyExtAttr
	row, err := o.Update(m)
	if err != nil {
		return false, "更新失败"
	}
	if row > 0 {
		_, ac := common.GetCacheAC()
		_ = ac.Delete("member-key-" + strconv.Itoa(m.ProjectId) + "-" + m.Name)
		_ = ac.Delete("member-key-" + strconv.Itoa(m.ProjectId) + "-" + m.Email)
		return true, "更新成功"
	}
	return false, "更新失败"
}

func (m *Member) UnbindById() (bool, string) {
	o := orm.NewOrm()
	err := o.Read(m)
	if err != nil {
		return false, "会员不存在"
	}

	m.Mac = ""
	row, err := o.Update(m)
	if err != nil {
		return false, "解绑失败"
	}
	if row > 0 {
		_, ac := common.GetCacheAC()
		_ = ac.Delete("member-key-" + strconv.Itoa(m.ProjectId) + "-" + m.Name)
		_ = ac.Delete("member-key-" + strconv.Itoa(m.ProjectId) + "-" + m.Email)
		return true, "解绑成功"
	}
	return false, "解绑失败"
}

func (m *Member) GetCount(mangerIdArr []int) int64 {
	o := orm.NewOrm()
	qs := o.QueryTable("Member")
	count, err := qs.Filter("ManagerId__in", mangerIdArr).Count()
	if err != nil {
		logs.Error("查询激活码总数错误：", err)
		return 0
	}
	return count
}

type Count struct {
	CountNum int64 `json:"count_num"`
}

func (m *Member) GetRangeCount(dateRange []common.DateObj, projectId int, managerIdArr []int) []int64 {
	var countArray []int64
	o := orm.NewOrm()
	qs := o.QueryTable(&m)
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

func (m *Member) LockSelect(lock int, id []int, managerIdArr []int) int64 {
	o := orm.NewOrm()
	qs := o.QueryTable(&m)
	var list []Member
	_, err := qs.Filter("Id__in", id).Filter("ManagerId__in", managerIdArr).All(&list, "Name", "ProjectId")
	_, ac := common.GetCacheAC()
	for _, i := range list {
		_ = ac.Delete("member-key-" + strconv.Itoa(i.ProjectId) + "-" + i.Name)
		_ = ac.Delete("member-key-" + strconv.Itoa(i.ProjectId) + "-" + i.Email)
	}
	rows, err := qs.Filter("Id__in", id).Filter("ManagerId__in", managerIdArr).Update(orm.Params{"IsLock": lock})
	if err != nil {
		return 0
	}
	return rows
}

func (m *Member) UnBindSelect(id []int, managerIdArr []int) int64 {
	o := orm.NewOrm()
	qs := o.QueryTable(&m)
	var list []Member
	_, err := qs.Filter("Id__in", id).Filter("ManagerId__in", managerIdArr).All(&list, "Name", "ProjectId")
	_, ac := common.GetCacheAC()
	for _, i := range list {
		_ = ac.Delete("member-key-" + strconv.Itoa(i.ProjectId) + "-" + i.Name)
		_ = ac.Delete("member-key-" + strconv.Itoa(i.ProjectId) + "-" + i.Email)
	}
	rows, err := qs.Filter("Id__in", id).Filter("ManagerId__in", managerIdArr).Update(orm.Params{"Mac": ""})
	if err != nil {
		return 0
	}
	return rows
}

func (m *Member) DeleteSelect(id []int, managerIdArr []int) int64 {
	o := orm.NewOrm()
	qs := o.QueryTable(&m)
	var list []Member
	_, err := qs.Filter("Id__in", id).Filter("ManagerId__in", managerIdArr).All(&list, "Name", "ProjectId")
	_, ac := common.GetCacheAC()
	for _, i := range list {
		_ = ac.Delete("member-key-" + strconv.Itoa(i.ProjectId) + "-" + i.Name)
		_ = ac.Delete("member-key-" + strconv.Itoa(i.ProjectId) + "-" + i.Email)
	}
	rows, err := qs.Filter("Id__in", id).Filter("ManagerId__in", managerIdArr).Delete()
	if err != nil {
		return 0
	}
	return rows
}

func (m *Member) LockQuery(lock int, managerIdArr []int, projectId int, mac string, member string, isLock int) int64 {
	var count int64
	var err error
	o := orm.NewOrm()
	qs := o.QueryTable(&m)
	if mac != "" {
		qs = qs.Filter("Mac", mac)
	}
	if member != "" {
		qs = qs.Filter("Name", member)
	}
	if isLock > -1 {
		if isLock == 0 {
			qs = qs.Filter("IsLock", 0)
		} else {
			qs = qs.Filter("IsLock", 1)
		}
	}

	qs = qs.Filter("ManagerId__in", managerIdArr)
	if projectId > -1 {
		qs = qs.Filter("ProjectId", projectId)
	}
	var list []Member
	_, err = qs.All(&list, "Name", "ProjectId")
	_, ac := common.GetCacheAC()
	for _, i := range list {
		_ = ac.Delete("member-key-" + strconv.Itoa(i.ProjectId) + "-" + i.Name)
		_ = ac.Delete("member-key-" + strconv.Itoa(i.ProjectId) + "-" + i.Email)
	}
	qs = o.QueryTable(&m)
	if mac != "" {
		qs = qs.Filter("Mac", mac)
	}
	if member != "" {
		qs = qs.Filter("Name", member)
	}
	if isLock > -1 {
		if isLock == 0 {
			qs = qs.Filter("IsLock", 0)
		} else {
			qs = qs.Filter("IsLock", 1)
		}
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

func (m *Member) UnBindQuery(managerIdArr []int, projectId int, mac string, member string, isLock int) int64 {
	var count int64
	var err error
	o := orm.NewOrm()
	qs := o.QueryTable(&m)
	if mac != "" {
		qs = qs.Filter("Mac", mac)
	}
	if member != "" {
		qs = qs.Filter("Name", member)
	}
	if isLock > -1 {
		if isLock == 0 {
			qs = qs.Filter("IsLock", 0)
		} else {
			qs = qs.Filter("IsLock", 1)
		}
	}

	qs = qs.Filter("ManagerId__in", managerIdArr)
	if projectId > -1 {
		qs = qs.Filter("ProjectId", projectId)
	}
	var list []Member
	_, err = qs.All(&list, "Name", "ProjectId")
	_, ac := common.GetCacheAC()
	for _, i := range list {
		_ = ac.Delete("member-key-" + strconv.Itoa(i.ProjectId) + "-" + i.Name)
		_ = ac.Delete("member-key-" + strconv.Itoa(i.ProjectId) + "-" + i.Email)
	}
	qs = o.QueryTable(&m)
	if mac != "" {
		qs = qs.Filter("Mac", mac)
	}
	if member != "" {
		qs = qs.Filter("Name", member)
	}
	if isLock > -1 {
		if isLock == 0 {
			qs = qs.Filter("IsLock", 0)
		} else {
			qs = qs.Filter("IsLock", 1)
		}
	}

	qs = qs.Filter("ManagerId__in", managerIdArr)
	if projectId == -1 {
		count, err = qs.Update(orm.Params{"Mac": ""})
	} else {
		count, err = qs.Filter("ProjectId", projectId).Update(orm.Params{"Mac": ""})
	}
	if err != nil {
		return 0
	}
	return count
}

func (m *Member) DeleteQuery(managerIdArr []int, projectId int, mac string, member string, isLock int) int64 {
	var count int64
	var err error
	o := orm.NewOrm()
	qs := o.QueryTable(&m)
	if mac != "" {
		qs = qs.Filter("Mac", mac)
	}
	if member != "" {
		qs = qs.Filter("Name", member)
	}
	if isLock > -1 {
		if isLock == 0 {
			qs = qs.Filter("IsLock", 0)
		} else {
			qs = qs.Filter("IsLock", 1)
		}
	}

	qs = qs.Filter("ManagerId__in", managerIdArr)
	if projectId > -1 {
		qs = qs.Filter("ProjectId", projectId)
	}
	var list []Member
	_, err = qs.All(&list, "Name", "ProjectId")
	_, ac := common.GetCacheAC()
	for _, i := range list {
		_ = ac.Delete("member-key-" + strconv.Itoa(i.ProjectId) + "-" + i.Name)
		_ = ac.Delete("member-key-" + strconv.Itoa(i.ProjectId) + "-" + i.Email)
	}
	qs = o.QueryTable(&m)
	if mac != "" {
		qs = qs.Filter("Mac", mac)
	}
	if member != "" {
		qs = qs.Filter("Name", member)
	}
	if isLock > -1 {
		if isLock == 0 {
			qs = qs.Filter("IsLock", 0)
		} else {
			qs = qs.Filter("IsLock", 1)
		}
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
