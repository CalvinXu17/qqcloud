package models

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
	"verification/controllers/common"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type ProjectVersion struct {
	ID           int       `orm:"column(id)" json:"id"`
	ManagerId    int       `orm:"index;default(0)" json:"manager_id"`
	ProjectId    int       `orm:"index;default(0);description(关联项目ID)" json:"project_id"`
	Version      float64   `orm:"digits(12);decimals(2);default(1.00);description(版本号)" json:"version" valid:"MaxSize(5)"`
	IsMustUpdate int       `orm:"index;default(0);description(0启用强制更新,1关闭强制更新)" json:"is_must_update" valid:"MaxSize(1)"`
	IsActive     int       `orm:"index;default(0);description(0启用,1关闭)" json:"is_active" valid:"MaxSize(1)"`
	Notice       string    `orm:"type(text);description(更新公告);size(9000);null" json:"notice" valid:"MaxSize(9000)"`
	WgtUrl       string    `orm:"null;description(热更新文件下载地址)" json:"wgt_url" valid:"MaxSize(255)"`
	CreateTime   time.Time `orm:"auto_now_add;type(datetime);index" json:"create_time"`
}

func (p *ProjectVersion) Add(projectId int, version float64, mustUpdate int, active int, noitce string, wgt string, manageId int) (id int64, msg string) {
	o := orm.NewOrm()
	project := Project{}
	project.ID = projectId
	err := o.Read(&project)
	if err != nil {
		return 0, "项目不存在"
	}
	v := ProjectVersion{
		ManagerId:    manageId,
		ProjectId:    projectId,
		Version:      version,
		IsMustUpdate: mustUpdate,
		IsActive:     active,
		Notice:       noitce,
		WgtUrl:       wgt,
	}

	id, err = o.Insert(&v)
	if err != nil {
		return 0, "创建失败"
	}
	_, ac := common.GetCacheAC()
	data, _ := json.Marshal(&v)
	var s strings.Builder
	s.WriteString(strconv.Itoa(projectId))
	s.WriteString("-")
	s.WriteString(common.GetVersionString(v.Version))
	_ = ac.Put(s.String(), string(data), 365*24*60*60*time.Second)
	var u strings.Builder
	u.WriteString("id-")
	u.WriteString(strconv.Itoa(v.ProjectId))
	u.WriteString("-up-version")
	_ = ac.Delete(u.String())
	return id, "创建成功"
}

func (p *ProjectVersion) GetVersionList(projectId int, pageSize int64, page int64) (status bool, pager Pager) {
	var data []ProjectVersion
	o := orm.NewOrm()
	var count int64
	var err error
	if projectId == 0 {
		count, err = o.QueryTable(&p).Count()
	} else {
		count, err = o.QueryTable(&p).Filter("ProjectId", projectId).Count()
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

func (c *ProjectVersion) Delete(id int) bool {
	p := ProjectVersion{
		ID: id,
	}
	o := orm.NewOrm()
	err := o.Read(&p)
	if err != nil {
		logs.Error("查询版本号失败")
		return false
	}
	rows, err := o.Delete(&p)
	if rows > 0 {
		_, ac := common.GetCacheAC()
		var s strings.Builder
		s.WriteString(strconv.Itoa(p.ProjectId))
		s.WriteString("-")
		s.WriteString(common.GetVersionString(p.Version))
		_ = ac.Delete(s.String())
		var u strings.Builder
		u.WriteString("id-")
		u.WriteString(strconv.Itoa(p.ProjectId))
		u.WriteString("-up-version")
		_ = ac.Delete(u.String())
		return true
	}
	return false
}

func (c *ProjectVersion) Update(id int, Version float64, IsMustUpdate int, IsActive int, Notice string, WgtUrl string) bool {
	_, ac := common.GetCacheAC()
	o := orm.NewOrm()
	p := ProjectVersion{ID: id}
	err := o.Read(&p)
	if err != nil {
		return false
	}
	if Notice != "" {
		p.Notice = Notice
	}
	if WgtUrl != "" {
		p.WgtUrl = WgtUrl
	}
	if Version > 0 {
		var s strings.Builder
		s.WriteString(strconv.Itoa(p.ProjectId))
		s.WriteString("-")
		s.WriteString(common.GetVersionString(p.Version))
		_ = ac.Delete(s.String())
		p.Version = Version
	}
	p.IsMustUpdate = IsMustUpdate
	p.IsActive = IsActive
	row, err := o.Update(&p)
	if row > 0 {
		data, _ := json.Marshal(&p)
		var s strings.Builder
		s.WriteString(strconv.Itoa(p.ProjectId))
		s.WriteString("-")
		s.WriteString(common.GetVersionString(p.Version))
		_ = ac.Put(s.String(), string(data), 365*24*60*60*time.Second)
		var u strings.Builder
		u.WriteString("id-")
		u.WriteString(strconv.Itoa(p.ProjectId))
		u.WriteString("-up-version")
		_ = ac.Delete(u.String())

		return true
	}
	return false
}

func (c *ProjectVersion) GetUpVersion() (bool, ProjectVersion) {
	p := ProjectVersion{}
	o := orm.NewOrm()
	qs := o.QueryTable("ProjectVersion")
	err := qs.OrderBy("-Version").One(&p)
	if err != nil {
		return false, p
	}
	_, ac := common.GetCacheAC()

	data, _ := json.Marshal(&p)
	var u strings.Builder
	u.WriteString("id-")
	u.WriteString(strconv.Itoa(p.ProjectId))
	u.WriteString("-up-version")
	versionUpString := common.Strval(ac.Get(u.String()))
	if versionUpString == "" {
		_ = ac.Put(u.String(), string(data), 365*24*60*60*time.Second)
	}
	return true, p
}

func (p *ProjectVersion) AllVersion() []ProjectVersion {
	data := make([]ProjectVersion, 0)
	o := orm.NewOrm()
	_, err := o.QueryTable(&p).All(&data)
	if err != nil {
		return data
	}
	return data
}
