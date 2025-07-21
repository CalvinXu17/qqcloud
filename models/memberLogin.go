package models

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
	"verification/controllers/common"
	"verification/validation/api"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

const (
	heartPrefix = "member-heart-"
)

type MemberLogin struct {
	ID            int       `orm:"column(id)" json:"id"`
	ManagerId     int       `orm:"index;default(0)" json:"manager_id"`
	ProjectId     int       `orm:"index;default(0)" json:"project_id"`
	Member        string    `orm:"size(60);null" json:"member"`
	MemberId      int       `orm:"default(0);index" json:"member_id"`
	Mac           string    `orm:"null;size(32);" json:"mac"`
	LastLoginTime int       `orm:"index;default(0)" json:"last_login_time"`
	LastLoginIp   string    `orm:"index;size(30);null" json:"last_login_ip"`
	IsOnline      int       `orm:"index;default(0);description(0 在线 1 下线)" json:"is_online"`
	OnlineId      string    `orm:"unique;size(56);index;null;description(单用户唯一在线id)" json:"online_id"`
	Type          int       `orm:"index;default(0);description(0登录,1解绑,2下线,3心跳)"`
	LoginType     int       `orm:"index;default(0);description(0普通登录,1微信登录,2QQ登录,3微博登录)" json:"login_type"`
	OpenId        string    `orm:"size(30);null" json:"open_id"`
	AccessToken   string    `orm:"size(50);null" json:"access_token"`
	CreateTime    time.Time `orm:"auto_now_add;type(datetime);index" json:"create_time"`
}

type OnlineMember struct {
	Member     string      `json:"member"`
	Mac        string      `json:"mac"`
	ClientList []OnlineObj `json:"client_list"`
	Addtime    int         `json:"addtime"`
	ClientTime int         `json:"clienttime"`
	Ip         string      `json:"ip"`
	Email      string      `json:"email"`
	Name       string      `json:"name"`
	NickName   string      `json:"nick_name"`
	MemberId   int         `json:"member_id"`
	ManagerId  int         `json:"manager_id"`
	ProjectId  int         `json:"project_id"`
	Client     string      `json:"client"`
	Uid        string      `json:"uid"`
}

type OnlineObj struct {
	Mac        string `json:"mac"`
	Addtime    int    `json:"addtime"`
	ClientTime int    `json:"clienttime"`
	Ip         string `json:"ip"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	NickName   string `json:"nick_name"`
	MemberId   int    `json:"member_id"`
	ManagerId  int    `json:"manager_id"`
	ProjectId  int    `json:"project_id"`
	Client     string `json:"client"`
	Uid        string `json:"uid"`
}

func AddLog(ip string, times int, mac string, logType int, token string, member Member) bool {
	timeEnd := int(time.Now().Unix())
	timeStart := timeEnd - (30 * 60)
	_, ac := common.GetCacheAC()
	m := MemberLogin{
		ManagerId:     member.ManagerId,
		ProjectId:     member.ProjectId,
		Member:        member.Name,
		MemberId:      member.ID,
		Mac:           mac,
		LastLoginTime: times,
		LastLoginIp:   ip,
		Type:          logType,
		LoginType:     0,
		AccessToken:   token,
	}
	logs.Error("写入日志", m)
	if logType == 3 {
		m.OnlineId = strconv.Itoa(member.ID)
		m.IsOnline = 0
	} else {
		m.OnlineId = common.GetToken()
	}
	if logType == 2 {
		var heartMacList []string
		var h strings.Builder
		h.WriteString(heartPrefix)
		h.WriteString(strconv.Itoa(member.ProjectId))
		h.WriteString("-")
		h.WriteString(strconv.Itoa(member.ID))
		heartJson := common.Strval(ac.Get(h.String()))
		logs.Info("在线列表：", heartJson)
		_ = json.Unmarshal([]byte(common.Strval(heartJson)), &heartMacList)
		index := 0
		for _, i := range heartMacList {
			heartIsOnline := common.Strval(ac.Get("h-" + i))
			if heartIsOnline == "" {
				index++
			}
		}
		if index == len(heartMacList) {

			o := orm.NewOrm()
			qs := o.QueryTable("MemberLogin")
			qs = qs.Filter("Type", 3)
			qs = qs.Filter("IsOnline", 0)
			qs = qs.Filter("LastLoginTime__gte", timeStart)
			qs = qs.Filter("LastLoginTime__lte", timeEnd)
			qs = qs.Filter("OnlineId", member.ID)
			row, _ := qs.Update(orm.Params{"IsOnline": 1})
			if row > 0 {
				logs.Error("清理指定会员在线成功")
			}
		}
	}
	o := orm.NewOrm()
	if logType == 3 {
		id, err := o.Insert(&m)
		if err != nil {
			qs := o.QueryTable("MemberLogin")
			row, err := qs.Filter("OnlineId", m.OnlineId).Update(orm.Params{
				"ManagerId":     member.ManagerId,
				"ProjectId":     member.ProjectId,
				"Member":        member.Name,
				"MemberId":      member.ID,
				"Mac":           mac,
				"LastLoginTime": times,
				"LastLoginIp":   ip,
				"Type":          logType,
				"LoginType":     0,
				"AccessToken":   token,
				"IsOnline":      0,
			})
			if err != nil {
				return false
			}
			if row > 0 {
				return true
			}
			return false
		}
		if id > 0 {
			return true
		}
		return false

	} else {
		id, err := o.Insert(&m)
		if err != nil {
			return false
		}
		if id > 0 {
			return true
		}
	}
	return false
}

func (m *MemberLogin) Add() bool {
	o := orm.NewOrm()
	id, err := o.Insert(m)
	if err != nil {
		logs.Error("写入登录记录失败", err)
		return false
	}
	if id > 0 {
		return true
	}
	return false
}

func (m *MemberLogin) FetchLog(bind int, actionType int, id int) (int64, error) {
	o := orm.NewOrm()
	if bind == 0 {
		year, month, day := time.Now().Date()
		t := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
		startDate, _ := time.Parse("2006-01-02 15:04:05", t.Format("2006-01-02 15:04:05"))
		startTimestamp := startDate.Unix()
		t = time.Date(year, month, day, 23, 59, 59, 0, time.Local)
		endDate, _ := time.Parse("2006-01-02 15:04:05", t.Format("2006-01-02 15:04:05"))
		endTimestamp := endDate.Unix()
		qs := o.QueryTable("MemberLogin")
		count, err := qs.Filter("LastLoginTime__gte", endTimestamp).Filter("LastLoginTime__lte", startTimestamp).Filter("Type", actionType).Filter("MemberId", id).Count()
		return count, err
	} else {

		year, month, _ := time.Now().Date()
		t := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
		monthStartDate, _ := time.Parse("2006-01-02 15:04:05", t.Format("2006-01-02 15:04:05"))
		monthStartTimestamp := monthStartDate.Unix()
		year, month, day := time.Now().AddDate(0, 1, -1).Date()
		t = time.Date(year, month, day, 23, 59, 59, 0, time.Local)
		monthEndDate, _ := time.Parse("2006-01-02 15:04:05", t.Format("2006-01-02 15:04:05"))
		monthEndTimestamp := monthEndDate.Unix()
		qs := o.QueryTable("MemberLogin")
		count, err := qs.Filter("LastLoginTime__gte", monthEndTimestamp).Filter("LastLoginTime__lte", monthStartTimestamp).Filter("Type", actionType).Filter("MemberId", id).Count()
		return count, err
	}
}

func (m *MemberLogin) FetchOnline(projectId int, pageSize int64, page int64, member string) (status bool, pager Pager) {

	var data []MemberLogin
	var count int64
	var err error
	timeEnd := int(time.Now().Unix())
	timeStart := timeEnd - (30 * 60)

	o := orm.NewOrm()
	qs := o.QueryTable("MemberLogin")
	if member != "" {
		qs = qs.Filter("Member", member)
	}
	qs = qs.Filter("Type", 3)
	qs = qs.Filter("IsOnline", 0)
	qs = qs.Filter("LastLoginTime__gte", timeStart)
	qs = qs.Filter("LastLoginTime__lte", timeEnd)
	if projectId == -1 {
		count, err = qs.Count()
	} else {
		count, err = qs.Filter("ProjectId", projectId).Count()
	}
	totalPage, offset, currentPage := PageCount(count, pageSize, page)
	qs = o.QueryTable("MemberLogin")
	if member != "" {
		qs = qs.Filter("Member", member)
	}
	qs = qs.Filter("Type", 3)
	qs = qs.Filter("IsOnline", 0)
	qs = qs.Filter("LastLoginTime__gte", timeStart)
	qs = qs.Filter("LastLoginTime__lte", timeEnd)
	if projectId == -1 {
		_, err = qs.Limit(pageSize, offset).All(&data)
	} else {
		_, err = qs.Filter("ProjectId", projectId).Limit(pageSize, offset).All(&data)
	}
	if err != nil {
		return false, Pager{
			Count:       count,
			CurrentPage: currentPage,
			Data:        data,
			PageSize:    pageSize,
			TotalPages:  totalPage,
		}
	}
	// 初始化缓存
	_, ac := common.GetCacheAC()
	var onlineList []OnlineMember
	for _, item := range data {
		var heartMacList []string
		var h strings.Builder
		h.WriteString(heartPrefix)
		h.WriteString(strconv.Itoa(item.ProjectId))
		h.WriteString("-")
		h.WriteString(strconv.Itoa(item.MemberId))
		heartJson := common.Strval(ac.Get(h.String()))
		logs.Info("在线列表：", heartJson)
		_ = json.Unmarshal([]byte(common.Strval(heartJson)), &heartMacList)
		var clientList []OnlineObj
		for _, i := range heartMacList {
			logs.Error("在线token:", i)
			var s strings.Builder
			s.WriteString("h-")
			s.WriteString(i)
			heartDataStr := common.Strval(ac.Get(s.String()))
			logs.Info("在线数据:", heartDataStr)
			if heartDataStr != "" {
				oldHeartData := api.OnlineData{}
				_ = json.Unmarshal([]byte(common.Strval(heartDataStr)), &oldHeartData)
				heartData := OnlineObj{
					Mac:        oldHeartData.Mac,
					Addtime:    oldHeartData.Addtime,
					ClientTime: oldHeartData.ClientTime,
					Ip:         oldHeartData.Ip,
					Email:      oldHeartData.Email,
					Name:       oldHeartData.Name,
					NickName:   oldHeartData.NickName,
					MemberId:   oldHeartData.MemberId,
					ManagerId:  oldHeartData.ManagerId,
					ProjectId:  oldHeartData.ProjectId,
					Client:     oldHeartData.Client,
					Uid:        common.GetToken(),
				}
				clientList = append(clientList, heartData)
			}
		}
		onlineItem := OnlineMember{
			ManagerId:  item.ManagerId,
			ProjectId:  item.ProjectId,
			Member:     item.Member,
			Name:       item.Member,
			MemberId:   item.MemberId,
			Mac:        item.Mac,
			ClientTime: item.LastLoginTime,
			Ip:         item.LastLoginIp,
			ClientList: clientList,
			Uid:        common.GetToken(),
		}
		onlineList = append(onlineList, onlineItem)
	}

	return true, Pager{
		Count:       count,
		CurrentPage: currentPage,
		Data:        onlineList,
		PageSize:    pageSize,
		TotalPages:  totalPage,
	}
}

func (m *MemberLogin) MemberLogout(projectId int, id int, client string) {
	timeEnd := int(time.Now().Unix())
	timeStart := timeEnd - (30 * 60)
	_, ac := common.GetCacheAC()
	if client == "" {
		var heartMacList []string
		var h strings.Builder
		h.WriteString(heartPrefix)
		h.WriteString(strconv.Itoa(projectId))
		h.WriteString("-")
		h.WriteString(strconv.Itoa(id))
		heartJson := common.Strval(ac.Get(h.String()))
		logs.Info("在线列表：", heartJson)
		_ = json.Unmarshal([]byte(common.Strval(heartJson)), &heartMacList)
		for _, i := range heartMacList {
			var s strings.Builder
			s.WriteString("h-o")
			s.WriteString(i)
			_ = ac.Put(s.String(), "强制下线", 1*60*60*time.Second)
			_ = ac.Delete("h-" + i)
		}
	} else {
		var s strings.Builder
		s.WriteString("h-o")
		s.WriteString(client)
		_ = ac.Put(s.String(), "强制下线", 1*60*60*time.Second)
		_ = ac.Delete("h-" + client)
	}
	var heartMacList []string
	var h strings.Builder
	h.WriteString(heartPrefix)
	h.WriteString(strconv.Itoa(projectId))
	h.WriteString("-")
	h.WriteString(strconv.Itoa(id))
	heartJson := common.Strval(ac.Get(h.String()))
	_ = json.Unmarshal([]byte(common.Strval(heartJson)), &heartMacList)
	index := 0
	for _, i := range heartMacList {
		heartIsOnline := common.Strval(ac.Get("h-" + i))
		if heartIsOnline == "" {
			index++
		}
	}
	if index == len(heartMacList) {

		o := orm.NewOrm()
		qs := o.QueryTable("MemberLogin")
		qs = qs.Filter("Type", 3)
		qs = qs.Filter("IsOnline", 0)
		qs = qs.Filter("LastLoginTime__gte", timeStart)
		qs = qs.Filter("LastLoginTime__lte", timeEnd)
		qs = qs.Filter("OnlineId", id)
		row, _ := qs.Update(orm.Params{"IsOnline": 1})
		if row > 0 {
			logs.Error("清理指定会员在线成功")
		}
	}
}

func (p *MemberLogin) GetRangeCount(dateRange []common.DateObj, projectId int, managerIdArr []int) []int64 {
	var countArray []int64
	o := orm.NewOrm()
	qs := o.QueryTable(p)
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
