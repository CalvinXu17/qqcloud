package admin

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	"strconv"
	"strings"
	"time"
	"verification/controllers/common"
	"verification/models"
	"verification/validation/admin"
)

type ProjectController struct {
	BaseController
}

// GetProjectList @Title 获取项目列表
// @router /getProjectList [post]
func (p *ProjectController) GetProjectList() {
	p.IsDeveloper()
	limit, _ := p.GetInt64("limit", 10)
	page, _ := p.GetInt64("page", 1)
	var project *models.Project
	status, list := project.GetProjectList(limit, page)
	if status == true {
		p.Success(0, list, "获取成功")
	}
	p.Error(400, "拉取列表失败")
}

// ProjectList @Title 获取全部项目列表
// @router /projectList [post]
func (p *ProjectController) ProjectList() {
	var status bool
	var list models.Pager
	var project *models.Project
	if p.ManagerInfo.Pid > 0 {
		m := models.ManagerCards{}
		idList := m.GetAgentProjectId(p.ManagerId)
		if len(idList) == 0 {
			p.Success(0, list, "获取成功")
		}
		status, list = project.GetAgentProjectList(idList)
	} else {
		status, list = project.ProjectList()
	}

	if status == true {
		p.Success(0, list, "获取成功")
	}
	p.Error(400, "拉取列表失败")
}

// CreateProject @Title 创建新项目
// @router /createProject [post]
func (p *ProjectController) CreateProject() {
	p.IsDeveloper()
	name := p.GetString("name", "")
	if name == "" {
		p.Error(400, "请填写项目名称")
	}
	projectType, _ := p.GetInt("type", 0)
	statusType, _ := p.GetInt("status_type", 0)
	encrypt, _ := p.GetInt("encrypt", 0)
	notice := p.GetString("notice", "")
	api := p.GetString("api", "")
	sign, _ := p.GetInt("sign", 0)
	managerId := common.GetManagerId(p.GetSession("token"))
	var project *models.Project
	id := project.Add(name, projectType, statusType, encrypt, notice, api, managerId, sign)
	if id == 0 {
		p.Error(400, "创建失败")
	}
	p.Success(0, id, "创建成功")
}

// UpdateProject @Title 创建新项目
// @router /updateProject [post]
func (p *ProjectController) UpdateProject() {
	p.IsDeveloper()
	id, _ := p.GetInt("id", 0)
	name := p.GetString("name")
	if name == "" {
		p.Error(400, "请填写项目名称")
	}
	projectType, _ := p.GetInt("type", 0)
	statusType, _ := p.GetInt("status_type", 0)
	encrypt, _ := p.GetInt("encrypt", 0)
	notice := p.GetString("notice", "")
	api := p.GetString("api", "")
	updateRsa, _ := p.GetInt("update_rsa", 0)
	updateKey, _ := p.GetInt("update_key", 0)
	updateAppKey, _ := p.GetInt("update_app_key", 0)
	updateSecretKey, _ := p.GetInt("update_secret_key", 0)
	sign, _ := p.GetInt("sign", 0)
	var project *models.Project
	status := project.Update(id, name, projectType, statusType, encrypt, notice, api, updateRsa, updateKey, updateAppKey, updateSecretKey, sign)
	if status == false {
		p.Error(400, "更新失败")
	}
	p.Success(0, id, "更新成功")
}

// DeleteProject @Title 删除项目
// @router /deleteProject [post]
func (p *ProjectController) DeleteProject() {
	p.IsDeveloper()
	id, _ := p.GetInt("id", 0)
	var project *models.Project
	status := project.Delete(id)
	if status == false {
		p.Error(400, "删除失败")
	}
	p.Success(0, id, "删除成功")
}

// bindProjectLogin @Title 项目绑定登录规则
// @router /bindProjectLogin [post]
func (p *ProjectController) BindProjectLogin() {
	p.IsDeveloper()
	id, _ := p.GetInt("id", 0)
	projectLoginId, _ := p.GetInt("project_login_id", 0)
	var project *models.Project
	status, msg := project.Bind(id, projectLoginId)
	if status == false {
		p.Error(400, msg)
	}
	p.Success(0, id, msg)
}

// CreateVersion @Title 创建版本
// @router /createVersion [post]
func (p *ProjectController) CreateVersion() {
	p.IsDeveloper()
	projectId, _ := p.GetInt("project_id", 0)
	version, _ := p.GetFloat("version", 1.00)
	mustUpdate, _ := p.GetInt("is_must_update", 0)
	active, _ := p.GetInt("is_active", 0)
	notice := p.GetString("notice", "")
	wgt := p.GetString("wgt_url", "")
	var projectVersion *models.ProjectVersion
	managerId := common.GetManagerId(p.GetSession("token"))
	id, msg := projectVersion.Add(projectId, version, mustUpdate, active, notice, wgt, managerId)
	if id > 0 {
		p.Success(0, id, msg)
	}
	p.Error(400, msg)
}

// GetVersionList @Title 获取版本号列表
// @router /getVersionList [post]
func (p *ProjectController) GetVersionList() {
	p.IsDeveloper()
	projectId, _ := p.GetInt("project_id", 0)
	limit, _ := p.GetInt64("limit", 10)
	page, _ := p.GetInt64("page", 1)
	var projectVersion *models.ProjectVersion
	status, list := projectVersion.GetVersionList(projectId, limit, page)
	if status == true {
		p.Success(0, list, "获取成功")
	}
	p.Error(400, "拉取列表失败")
}

// DeleteProjectVersion @Title 删除项目版本号
// @router /deleteProjectVersion [post]
func (p *ProjectController) DeleteProjectVersion() {
	p.IsDeveloper()
	id, _ := p.GetInt("id", 0)
	var version *models.ProjectVersion
	status := version.Delete(id)
	if status == false {
		p.Error(400, "删除失败")
	}
	p.Success(0, id, "删除成功")
}

// UpdateProjectVersion @Title 更新版本号
// @router /updateProjectVersion [post]
func (p *ProjectController) UpdateProjectVersion() {
	p.IsDeveloper()
	ID, _ := p.GetInt("id", 0)
	Version, _ := p.GetFloat("version", 1)
	IsMustUpdate, _ := p.GetInt("is_must_update", 0)
	IsActive, _ := p.GetInt("is_active", 0)
	Notice := p.GetString("notice", "")
	WgtUrl := p.GetString("wgt_url", "")
	var v *models.ProjectVersion
	status := v.Update(ID, Version, IsMustUpdate, IsActive, Notice, WgtUrl)
	if status == false {
		p.Error(400, "更新失败")
	}
	p.Success(0, ID, "更新成功")
}

// CreateCard @Title 创建激活码类型
// @router /createCard [post]
func (p *ProjectController) CreateCard() {
	p.IsDeveloper()
	ProjectId, _ := p.GetInt("project_id", 0)
	Title := p.GetString("title", "")
	if Title == "" {
		p.Error(400, "请填写类型名称")
	}
	Price, _ := p.GetFloat("price", 1)
	KeyPrefix := p.GetString("key_prefix", "")
	LevelId, _ := p.GetInt("level_id", 0)
	Days, _ := p.GetFloat("days", 1)
	Points, _ := p.GetInt("points", 0)
	KeyExtAttr := p.GetString("key_ext_attr", "")
	Tag := p.GetString("tag", "")
	IsLock, _ := p.GetInt("is_lock", 0)
	if Tag == "免费用户" {
		p.Error(400, "标签：【免费用户】不可用")
	}
	managerId := common.GetManagerId(p.GetSession("token"))
	valid := validation.Validation{}
	y := admin.CreateCards{
		Title:      Title,
		Price:      Price,
		KeyPrefix:  KeyPrefix,
		LevelId:    LevelId,
		Days:       Days,
		Points:     Points,
		KeyExtAttr: KeyExtAttr,
		Tag:        Tag,
		IsLock:     IsLock,
	}
	b, _ := valid.Valid(&y)
	if !b {
		for _, err := range valid.Errors {
			p.Error(400, err.Message)
		}
	}
	var c *models.Cards
	id := c.Add(ProjectId, Title, Price, KeyPrefix, LevelId, Days, Points, KeyExtAttr, Tag, IsLock, managerId)
	if id > 0 {
		p.Success(0, id, "创建成功")
	}
	p.Error(400, "创建失败")
}

// UpdateCard @Title 创建激活码类型
// @router /updateCard [post]
func (p *ProjectController) UpdateCard() {
	p.IsDeveloper()
	ID, _ := p.GetInt("id", 0)
	Price, _ := p.GetFloat("price", 1)
	Title := p.GetString("title", "")
	if Title == "" {
		p.Error(400, "请填写类型名称")
	}
	KeyPrefix := p.GetString("key_prefix", "")
	LevelId, _ := p.GetInt("level_id", 0)
	Days, _ := p.GetFloat("days", 1)
	Points, _ := p.GetInt("points", 0)
	KeyExtAttr := p.GetString("key_ext_attr", "")
	Tag := p.GetString("tag", "")
	IsLock, _ := p.GetInt("is_lock", 0)

	valid := validation.Validation{}
	y := admin.CreateCards{
		Title:      Title,
		Price:      Price,
		KeyPrefix:  KeyPrefix,
		LevelId:    LevelId,
		Days:       Days,
		Points:     Points,
		KeyExtAttr: KeyExtAttr,
		Tag:        Tag,
		IsLock:     IsLock,
	}
	b, _ := valid.Valid(&y)
	if !b {
		for _, err := range valid.Errors {
			p.Error(400, err.Message)
		}
	}
	var c *models.Cards
	id := c.Update(ID, Title, Price, KeyPrefix, LevelId, Days, Points, KeyExtAttr, Tag, IsLock)
	if id == true {
		p.Success(0, id, "修改成功")
	}
	p.Error(400, "修改失败")
}

// GetCardList @Title 获取激活码类型列表
// @router /getCardList [post]
func (p *ProjectController) GetCardList() {
	p.IsDeveloper()
	projectId, _ := p.GetInt("project_id", 0)
	limit, _ := p.GetInt64("limit", 10)
	page, _ := p.GetInt64("page", 1)
	var c *models.Cards
	status, list := c.GetCardList(projectId, limit, page)
	if status == true {
		p.Success(0, list, "获取成功")
	}
	p.Error(400, "拉取列表失败")
}

// CardList @Title 获取激活码类型列表
// @router /cardList [post]
func (p *ProjectController) CardList() {
	//p.IsDeveloper()
	var status bool
	var list models.Pager
	var c *models.Cards
	if p.ManagerInfo.Pid > 0 {
		m := models.ManagerCards{}
		statusManagerCards, pager := m.GetCardList(p.ManagerInfo.ID)
		if statusManagerCards == false {
			p.Error(400, "拉取列表失败")
		}
		var cardsIdList []int
		for _, i := range pager {
			cardsIdList = append(cardsIdList, i.CardsId)
		}
		if len(cardsIdList) == 0 {
			p.Success(0, list, "获取成功")
		}
		logs.Error("代理权限列表：", cardsIdList)
		status, list = c.GetAgentCardList(cardsIdList)
	} else {
		status, list = c.CardList()
	}
	if status == true {
		p.Success(0, list, "获取成功")
	}
	p.Error(400, "拉取列表失败")

}

// DeleteCard @Title 删除激活码类型
// @router /deleteCard [post]
func (p *ProjectController) DeleteCard() {
	p.IsDeveloper()
	id, _ := p.GetInt("id", 0)
	var c *models.Cards
	status := c.Delete(id)
	if status == false {
		p.Error(400, "删除失败")
	}
	p.Success(0, id, "删除成功")
}

// CreateLoginRule @Title 创建登录规则
// @router /createLoginRule [post]
func (p *ProjectController) CreateLoginRule() {
	p.IsDeveloper()
	Title := p.GetString("title", "")
	if Title == "" {
		p.Error(400, "请填写规则名称")
	}
	Mode, _ := p.GetInt("mode", 0)
	RegMode, _ := p.GetInt("reg_mode", 0)
	EmailReg, _ := p.GetInt("email_reg", 0)
	UnbindMode, _ := p.GetInt("unbind_mode", 0)
	UnbindWeakenMode, _ := p.GetInt("unbind_weaken_mode", 0)
	UnbindWeaken, _ := p.GetFloat("unbind_weaken", 0)
	UnbindWeakenPoints, _ := p.GetInt("unbind_weaken_points", 0)
	UnbindTimes, _ := p.GetInt("unbind_times", 0)
	UnbindDate, _ := p.GetInt("UnbindDate", 0)
	UnbindBefore, _ := p.GetInt("unbind_before", 0)
	NumberMode, _ := p.GetInt("number_mode", 0)
	NumberMore, _ := p.GetInt("number_more", 0)
	NumberWeaken, _ := p.GetInt("number_weaken", 0)
	NumberWeakenTime, _ := p.GetInt("number_weaken_time", 0)
	PcMore, _ := p.GetInt("pc_more", 0)
	PcCodeMore, _ := p.GetInt("pc_code_more", 0)
	managerId := common.GetManagerId(p.GetSession("token"))
	m := models.ProjectLogin{
		Title:              Title,
		Mode:               Mode,
		RegMode:            RegMode,
		EmailReg:           EmailReg,
		UnbindMode:         UnbindMode,
		UnbindWeakenMode:   UnbindWeakenMode,
		UnbindWeaken:       UnbindWeaken,
		UnbindWeakenPoints: UnbindWeakenPoints,
		UnbindTimes:        UnbindTimes,
		UnbindDate:         UnbindDate,
		UnbindBefore:       UnbindBefore,
		NumberMode:         NumberMode,
		NumberMore:         NumberMore,
		NumberWeaken:       NumberWeaken,
		NumberWeakenTime:   NumberWeakenTime,
		PcMore:             PcMore,
		PcCodeMore:         PcCodeMore,
		ManagerId:          managerId,
	}
	id := m.Add()
	if id > 0 {
		p.Success(0, id, "创建成功")
	}
	p.Error(400, "创建失败")
}

// GetLoginRuleList @Title 获取登录规则列表
// @router /getLoginRuleList [post]
func (p *ProjectController) GetLoginRuleList() {
	p.IsDeveloper()
	limit, _ := p.GetInt64("limit", 10)
	page, _ := p.GetInt64("page", 1)
	var c *models.ProjectLogin
	status, list := c.GetLoginRuleList(limit, page)
	if status == true {
		p.Success(0, list, "获取成功")
	}
	p.Error(400, "拉取列表失败")
}

// LoginRuleList @Title 获取登录规则列表
// @router /loginRuleList [post]
func (p *ProjectController) LoginRuleList() {
	p.IsDeveloper()
	var c *models.ProjectLogin
	status, list := c.LoginRuleList()
	if status == true {
		p.Success(0, list, "获取成功")
	}
	p.Error(400, "拉取列表失败")
}

// UpdateLoginRule @Title 修改登录规则
// @router /updateLoginRule [post]
func (p *ProjectController) UpdateLoginRule() {
	p.IsDeveloper()
	ID, _ := p.GetInt("id", 0)
	Title := p.GetString("title", "")
	if Title == "" {
		p.Error(400, "请填写规则名称")
	}
	Mode, _ := p.GetInt("mode", 0)
	RegMode, _ := p.GetInt("reg_mode", 0)
	EmailReg, _ := p.GetInt("email_reg", 0)
	UnbindMode, _ := p.GetInt("unbind_mode", 0)
	UnbindWeakenMode, _ := p.GetInt("unbind_weaken_mode", 0)
	UnbindWeaken, _ := p.GetFloat("unbind_weaken", 0)
	UnbindWeakenPoints, _ := p.GetInt("unbind_weaken_points", 0)
	UnbindTimes, _ := p.GetInt("unbind_times", 0)
	UnbindDate, _ := p.GetInt("UnbindDate", 0)
	UnbindBefore, _ := p.GetInt("unbind_before", 0)
	NumberMode, _ := p.GetInt("number_mode", 0)
	NumberMore, _ := p.GetInt("number_more", 0)
	NumberWeaken, _ := p.GetInt("number_weaken", 0)
	NumberWeakenTime, _ := p.GetInt("number_weaken_time", 0)
	PcMore, _ := p.GetInt("pc_more", 0)
	PcCodeMore, _ := p.GetInt("pc_code_more", 0)
	m := models.ProjectLogin{ID: ID}
	status := m.Update(models.ProjectLogin{
		Title:              Title,
		Mode:               Mode,
		RegMode:            RegMode,
		EmailReg:           EmailReg,
		UnbindMode:         UnbindMode,
		UnbindWeakenMode:   UnbindWeakenMode,
		UnbindWeaken:       UnbindWeaken,
		UnbindWeakenPoints: UnbindWeakenPoints,
		UnbindTimes:        UnbindTimes,
		UnbindDate:         UnbindDate,
		UnbindBefore:       UnbindBefore,
		NumberMode:         NumberMode,
		NumberMore:         NumberMore,
		NumberWeaken:       NumberWeaken,
		NumberWeakenTime:   NumberWeakenTime,
		PcMore:             PcMore,
		PcCodeMore:         PcCodeMore,
	})
	if status > 0 {
		p.Success(0, status, "更新成功")
	}
	p.Error(400, "更新失败")
}

// DeleteLoginRule @Title 删除登录规则
// @router /deleteLoginRule [post]
func (p *ProjectController) DeleteLoginRule() {
	p.IsDeveloper()
	id, _ := p.GetInt("id", 0)
	m := models.ProjectLogin{ID: id}
	row := m.Delete()
	if row > 0 {
		p.Success(0, row, "删除成功")
	}
	p.Error(400, "删除失败")
}

// QueryOrderRmb @Title 预创建激活码金额查询
// @router /queryOrderRmb [post]
func (p *ProjectController) QueryOrderRmb() {
	cardsId, _ := p.GetInt("cards_id", 0)
	count, _ := p.GetInt64("count", 1)
	length, _ := p.GetInt("length", 8)
	createType, _ := p.GetInt("create_type", 0)
	projectId, _ := p.GetInt("project_id", -1)
	tag := p.GetString("tag", "")
	var agentPrice float64
	agentPrice = 0
	if tag == "免费用户" {
		p.Error(400, "标签：【免费用户】不可用")
	}
	if p.ManagerInfo.Pid > 0 {
		m := models.ManagerCards{}
		status, pager := m.GetCardList(p.ManagerInfo.ID)
		if status == false {
			p.Error(400, "您暂无可用的授权激活码类型")
		}
		hasPermission := false
		for _, i := range pager {
			if i.CardsId == cardsId {
				hasPermission = true
				agentPrice = i.Price
			}
		}
		if hasPermission == false {
			p.Error(400, "暂无此激活码类型提卡权限")
		}
	}
	valid := validation.Validation{}
	y := admin.CreateKeys{
		CardsId:    cardsId,
		Count:      count,
		Length:     length,
		CreateType: createType,
		ProjectId:  projectId,
		Tag:        tag,
	}
	b, _ := valid.Valid(&y)
	if !b {
		for _, err := range valid.Errors {
			p.Error(400, err.Message)
		}
	}
	managerId := common.GetManagerId(p.GetSession("token"))
	var c *models.Keys
	type res struct {
		Status bool    `json:"status"`
		Msg    string  `json:"msg"`
		Cost   float64 `json:"cost"`
	}
	status, cost, msg := c.QueryCreate(cardsId, count, managerId, agentPrice)
	result := res{
		Status: status,
		Msg:    msg,
		Cost:   cost,
	}
	p.Success(0, result, msg)
}

// CreateKeys @Title 创建激活码
// @router /createKeys [post]
func (p *ProjectController) CreateKeys() {
	cardsId, _ := p.GetInt("cards_id", 0)
	count, _ := p.GetInt64("count", 1)
	length, _ := p.GetInt("length", 8)
	createType, _ := p.GetInt("create_type", 0)
	projectId, _ := p.GetInt("project_id", -1)
	tag := p.GetString("tag", "")
	var agentPrice float64
	agentPrice = 0
	if tag == "免费用户" {
		p.Error(400, "标签：【免费用户】不可用")
	}
	if p.ManagerInfo.Pid > 0 {
		m := models.ManagerCards{}
		status, pager := m.GetCardList(p.ManagerInfo.ID)
		if status == false {
			p.Error(400, "您暂无可用的授权激活码类型")
		}
		hasPermission := false
		for _, i := range pager {
			if i.CardsId == cardsId {
				hasPermission = true
				agentPrice = i.Price
			}
		}
		if hasPermission == false {
			p.Error(400, "暂无此激活码类型提卡权限")
		}
	}
	valid := validation.Validation{}
	y := admin.CreateKeys{
		CardsId:    cardsId,
		Count:      count,
		Length:     length,
		CreateType: createType,
		ProjectId:  projectId,
		Tag:        tag,
	}
	b, _ := valid.Valid(&y)
	if !b {
		for _, err := range valid.Errors {
			p.Error(400, err.Message)
		}
	}
	managerId := common.GetManagerId(p.GetSession("token"))
	var c *models.Keys
	status, keys, msg := c.Create(cardsId, count, length, createType, managerId, projectId, tag, agentPrice)
	if status == true {
		p.Success(0, keys, msg)
	}
	p.Error(400, msg)
}

// GetKeysList @Title 获取激活码列表
// @router /getKeysList [post]
func (p *ProjectController) GetKeysList() {
	projectId, _ := p.GetInt("project_id", 0)
	limit, _ := p.GetInt64("limit", 10)
	page, _ := p.GetInt64("page", 1)
	longKeys := p.GetString("long_keys", "")
	cardsId, _ := p.GetInt("cards_id", 0)
	isActive, _ := p.GetInt("is_active", -1)
	isLock, _ := p.GetInt("is_lock", -1)
	member := p.GetString("member", "")
	orderId, _ := p.GetInt("order_id", 0)
	var k *models.Keys
	status, list := k.GetKeysList(p.ManagerIdArr, projectId, limit, page, longKeys, cardsId, isActive, isLock, member, orderId)
	if status == true {
		p.Success(0, list, "获取成功")
	}
	p.Error(400, "拉取列表失败")
}

// BatchKeys @Title 批量操作激活码
// @router /batchKeys [post]
func (p *ProjectController) BatchKeys() {
	var idList []int
	var rows int64
	projectId, _ := p.GetInt("project_id", 0)
	longKeys := p.GetString("long_keys", "")
	cardsId, _ := p.GetInt("cards_id", 0)
	isActive, _ := p.GetInt("is_active", -1)
	isLock, _ := p.GetInt("is_lock", -1)
	member := p.GetString("member", "")
	orderId, _ := p.GetInt("order_id", 0)
	id := p.GetString("id", "")
	opType := p.GetString("type", "")
	if id != "" {
		err := json.Unmarshal([]byte(id), &idList)
		if err != nil {
			p.Error(400, "提供的ID参数格式错误")
		}
	}
	var k *models.Keys
	switch opType {
	case "lock_1":
		rows = k.LockSelect(1, idList, p.ManagerIdArr)
		p.Success(0, nil, "成功锁定"+strconv.FormatInt(rows, 10)+"条数据")
		break
	case "unlock_1":
		rows = k.LockSelect(0, idList, p.ManagerIdArr)
		p.Success(0, nil, "成功解锁"+strconv.FormatInt(rows, 10)+"条数据")
		break
	case "delete_1":
		rows = k.DeleteSelect(idList, p.ManagerIdArr)
		p.Success(0, nil, "成功删除"+strconv.FormatInt(rows, 10)+"条数据")
		break
	case "lock_2":
		rows = k.LockQuery(1, p.ManagerIdArr, projectId, longKeys, cardsId, isActive, isLock, member, orderId)
		p.Success(0, nil, "成功锁定"+strconv.FormatInt(rows, 10)+"条数据")
		break
	case "unlock_2":
		rows = k.LockQuery(0, p.ManagerIdArr, projectId, longKeys, cardsId, isActive, isLock, member, orderId)
		p.Success(0, nil, "成功解锁"+strconv.FormatInt(rows, 10)+"条数据")
		break
	case "delete_2":
		logs.Error("代理列表", p.ManagerIdArr)
		rows = k.DeleteQuery(p.ManagerIdArr, projectId, longKeys, cardsId, isActive, isLock, member, orderId)
		p.Success(0, nil, "成功删除"+strconv.FormatInt(rows, 10)+"条数据")
		break
	default:
		p.Error(400, "未知操作类型")
	}

}

// LockKeys @Title 锁定解锁激活码
// @router /lockKey [post]
func (p *ProjectController) LockKeys() {
	id, _ := p.GetInt("id", 0)
	var c *models.Keys
	status, msg := c.Lock(id)
	if status == true {
		p.Success(0, nil, msg)
	}
	p.Error(400, msg)
}

// DeleteKeys @Title 删除激活码
// @router /deleteKeys [post]
func (p *ProjectController) DeleteKeys() {
	id, _ := p.GetInt("id", 0)
	var c *models.Keys
	status := c.Delete(id)
	if status == true {
		p.Success(0, nil, "删除成功")
	}
	p.Error(400, "删除失败")
}

// GetMemberList @Title 获取会员列表
// @router /getMemberList [post]
func (p *ProjectController) GetMemberList() {
	projectId, _ := p.GetInt("project_id", 0)
	limit, _ := p.GetInt64("limit", 10)
	page, _ := p.GetInt64("page", 1)
	mac := p.GetString("mac", "")
	member := p.GetString("member", "")
	isLock, _ := p.GetInt("is_lock", -1)
	var c *models.Member
	status, list := c.GetMemberList(p.ManagerIdArr, projectId, limit, page, mac, member, isLock)
	if status == true {
		p.Success(0, list, "获取成功")
	}
	p.Error(400, "拉取列表失败")
}

// BatchMember @Title 批量操作会员
// @router /batchMember [post]
func (p *ProjectController) BatchMember() {
	var idList []int
	projectId, _ := p.GetInt("project_id", 0)
	mac := p.GetString("mac", "")
	member := p.GetString("member", "")
	isLock, _ := p.GetInt("is_lock", -1)
	id := p.GetString("id", "")
	opType := p.GetString("type", "")
	if id != "" {
		err := json.Unmarshal([]byte(id), &idList)
		if err != nil {
			p.Error(400, "提供的ID参数格式错误")
		}
	}
	var rows int64
	var m *models.Member
	switch opType {
	case "lock_1":
		rows = m.LockSelect(1, idList, p.ManagerIdArr)
		p.Success(0, nil, "成功锁定"+strconv.FormatInt(rows, 10)+"条数据")
		break
	case "unlock_1":
		rows = m.LockSelect(0, idList, p.ManagerIdArr)
		p.Success(0, nil, "成功解锁"+strconv.FormatInt(rows, 10)+"条数据")
		break
	case "unbind_1":
		rows = m.UnBindSelect(idList, p.ManagerIdArr)
		p.Success(0, nil, "成功解绑"+strconv.FormatInt(rows, 10)+"条数据")
		break
	case "delete_1":
		rows = m.DeleteSelect(idList, p.ManagerIdArr)
		p.Success(0, nil, "成功删除"+strconv.FormatInt(rows, 10)+"条数据")
		break
	case "lock_2":
		rows = m.LockQuery(1, p.ManagerIdArr, projectId, mac, member, isLock)
		p.Success(0, nil, "成功锁定"+strconv.FormatInt(rows, 10)+"条数据")
		break
	case "unlock_2":
		rows = m.LockQuery(0, p.ManagerIdArr, projectId, mac, member, isLock)
		p.Success(0, nil, "成功解锁"+strconv.FormatInt(rows, 10)+"条数据")
		break
	case "unbind_2":
		rows = m.UnBindQuery(p.ManagerIdArr, projectId, mac, member, isLock)
		p.Success(0, nil, "成功解绑"+strconv.FormatInt(rows, 10)+"条数据")
		break
	case "delete_2":
		m.DeleteQuery(p.ManagerIdArr, projectId, mac, member, isLock)
		p.Success(0, nil, "成功删除"+strconv.FormatInt(rows, 10)+"条数据")
		break
	default:
		p.Error(400, "错误的操作类型")
	}
}

// DeleteMember @Title 删除会员
// @router /deleteMember [post]
func (p *ProjectController) DeleteMember() {
	id, _ := p.GetInt("id", 0)
	m := models.Member{ID: id}
	status, msg := m.DeleteById()
	if status == true {
		p.Success(0, nil, msg)
	}
	p.Error(400, msg)
}

// LockMember @Title 锁定解锁会员
// @router /lockMember [post]
func (p *ProjectController) LockMember() {
	id, _ := p.GetInt("id", 0)
	m := models.Member{ID: id}
	status, msg := m.LockById()
	if status == true {
		p.Success(0, nil, msg)
	}
	p.Error(400, msg)
}

// UpdateMember @Title 更新会员信息
// @router /updateMember [post]
func (p *ProjectController) UpdateMember() {
	id, _ := p.GetInt("id", 0)
	m := models.Member{ID: id}
	name := p.GetString("name", "")
	nickName := p.GetString("nickname", "")
	password := p.GetString("password", "")
	safePassword := p.GetString("safe_password", "")
	days, _ := p.GetFloat("days", 0)
	points, _ := p.GetInt("points", 0)
	mac := p.GetString("mac", "")
	phone, _ := p.GetInt64("phone", 0)
	email := p.GetString("email", "")
	endtime, _ := p.GetInt("endtime", 0)
	member := models.Member{
		Name:         name,
		NickName:     nickName,
		Password:     password,
		SafePassword: safePassword,
		Days:         days,
		Points:       points,
		Mac:          mac,
		Phone:        phone,
		Email:        email,
		EndTime:      endtime,
	}
	status, msg := m.UpdateById(member)
	if status == true {
		p.Success(0, nil, msg)
	}
	p.Error(400, msg)
}

// UnbindMember @Title 会员解绑
// @router /unbindMember [post]
func (p *ProjectController) UnbindMember() {
	id, _ := p.GetInt("id", 0)
	m := models.Member{ID: id}
	status, msg := m.UnbindById()
	if status == true {
		p.Success(0, nil, msg)
	}
	p.Error(400, msg)
}

// GetOnlineList @Title 获取在线会员
// @router /getOnlineList [post]
func (p *ProjectController) GetOnlineList() {
	projectId, _ := p.GetInt("project_id", -1)
	limit, _ := p.GetInt64("limit", 10)
	page, _ := p.GetInt64("page", 1)
	member := p.GetString("member", "")
	var c *models.MemberLogin
	status, list := c.FetchOnline(projectId, limit, page, member)
	if status == true {
		p.Success(0, list, "获取成功")
	}
	p.Error(400, "拉取列表失败")
}

// MemberLogout @Title 会员下线
// @router /memberLogout [post]
func (p *ProjectController) MemberLogout() {
	projectId, _ := p.GetInt("project_id", 0)
	id, _ := p.GetInt("id", 0)
	client := p.GetString("client", "")
	var m *models.MemberLogin
	m.MemberLogout(projectId, id, client)
	p.Success(0, nil, "下线成功")
}

// GetManagerList @Title 获取所有代理
// @router /getManagerList [post]
func (p *ProjectController) GetManagerList() {
	p.IsDeveloper()
	managerId := common.GetManagerId(p.GetSession("token"))
	var m *models.Manager
	_, data := m.GetManagerList(managerId, p.ManagerIdArr)
	p.Success(0, data, "获取成功")
}

// ManagerAdd @Title 创建代理
// @router /managerAdd [post]
func (p *ProjectController) ManagerAdd() {
	p.IsDeveloper()
	money, _ := p.GetFloat("money", 0)
	powerId, _ := p.GetInt("power_id", 0)
	user := p.GetString("user", "")
	pwd := p.GetString("pwd", "")
	email := p.GetString("email", "")
	if user == "" {
		p.Error(400, "账号不能为空")
	}
	if pwd == "" {
		p.Error(400, "密码不能为空")
	}
	if email == "" {
		p.Error(400, "邮箱不能为空")
	}
	managerId := common.GetManagerId(p.GetSession("token"))
	i := models.Manager{ID: managerId}
	i, _ = i.GetInfoById(managerId)
	level := i.Level + 1
	m := models.Manager{
		User:    user,
		Pwd:     pwd,
		Money:   money,
		PowerId: powerId,
		Email:   email,
		Pid:     managerId,
		Level:   level,
	}
	id, msg := m.AddManager()
	if id > 0 {
		p.Success(0, id, "创建成功")
	}
	p.Error(400, msg)
}

// ManagerUpdate @Title 修改代理
// @router /managerUpdate [post]
func (p *ProjectController) ManagerUpdate() {
	p.IsDeveloper()
	id, _ := p.GetInt("id", 0)
	money, _ := p.GetFloat("money", 0)
	powerId, _ := p.GetInt("power_id", 0)
	user := p.GetString("user", "")
	pwd := p.GetString("pwd", "")
	email := p.GetString("email", "")
	if user == "" {
		p.Error(400, "账号不能为空")
	}
	if pwd == "" {
		p.Error(400, "密码不能为空")
	}
	if email == "" {
		p.Error(400, "邮箱不能为空")
	}
	m := models.Manager{ID: id}
	n := models.Manager{
		User:    user,
		Pwd:     pwd,
		Money:   money,
		PowerId: powerId,
		Email:   email,
	}
	row, msg := m.UpdateManager(n)
	if row > 0 {
		p.Success(0, row, "修改成功")
	}
	p.Error(400, msg)
}

// ManagerDelete @Title 删除代理
// @router /managerDelete [post]
func (p *ProjectController) ManagerDelete() {
	p.IsDeveloper()
	id, _ := p.GetInt("id", 0)
	index := 0
	for _, i := range p.ManagerIdArr {
		if i == id {
			index = 1
		}
	}
	if index == 0 {
		p.Error(400, "无法删除不属于自己的代理")
	}
	if id == p.ManagerId {
		p.Error(400, "无法删除自身")
	}
	m := models.Manager{ID: id}
	row := m.Delete()
	if row > 0 {
		p.Success(0, row, "删除成功")
	}
	p.Error(400, "删除失败")
}

// GetManagerCards @Title 获取代理权限
// @router /getManagerCards [post]
func (p *ProjectController) GetManagerCards() {
	id, _ := p.GetInt("id", 0)
	p.IsInMangaerList(id)
	m := models.ManagerCards{}
	status, list := m.CardList(id)
	if status == false {
		p.Error(400, "获取失败")
	}
	p.Success(0, list, "获取成功")
}

// AddManagerCards @Title 添加授权激活码
// @router /addManagerCards [post]
func (p *ProjectController) AddManagerCards() {
	managerId, _ := p.GetInt("manager_id", 0)
	logs.Error("managerId", managerId)
	p.IsInMangaerList(managerId)
	cardsId, _ := p.GetInt("cards_id", 0)
	price, _ := p.GetFloat("price", 0)
	c := models.Cards{ID: cardsId}
	status, cards := c.GetById()
	if status == false {
		p.Error(400, "激活码类型不存在")
	}
	m := models.ManagerCards{
		ProjectId: cards.ProjectId,
		ManagerId: managerId,
		CardsId:   cards.ID,
		Price:     price,
	}
	row, msg := m.Add()
	if row > 0 {
		p.Success(0, row, msg)
	}
	p.Error(400, msg)
}

// UpdateManagerCards @Title 添加授权激活码
// @router /updateManagerCards [post]
func (p *ProjectController) UpdateManagerCards() {
	id, _ := p.GetInt("id", 0)
	managerId, _ := p.GetInt("manager_id", 0)
	p.IsInMangaerList(managerId)
	price, _ := p.GetFloat("price", 0)
	m := models.ManagerCards{
		ID:        id,
		Price:     price,
		ManagerId: managerId,
	}
	status, msg := m.Update()
	if status == true {
		p.Success(0, status, msg)
	}
	p.Error(400, msg)
}

// DeleteManagerCards @Title 删除角色激活码权限
// @router /deleteManagerCards [post]
func (p *ProjectController) DeleteManagerCards() {
	p.IsDeveloper()
	id, _ := p.GetInt("id", 0)
	managerId, _ := p.GetInt("manager_id", 0)
	p.IsInMangaerList(managerId)
	r := models.ManagerCards{ManagerId: managerId, ID: id}
	row, msg := r.Delete()
	if row > 0 {
		p.Success(0, row, msg)
	}
	p.Error(400, msg)
}

// ManagerIdList @Title 获取代理ID
// @router /managerIdList [post]
func (p *ProjectController) ManagerIdList() {
	p.IsDeveloper()
	p.Success(0, p.ManagerIdArr, "获取成功")
}

// RolePower @Title 获取预设权限列表
// @router /rolePower [post]
func (p *ProjectController) RolePower() {
	p.IsDeveloper()
	p.Success(0, models.RoleList, "获取成功")
}

// GetRoleUser @Title 获取角色列表
// @router /getRoleUser [post]
func (p *ProjectController) GetRoleUser() {
	p.IsDeveloper()
	limit, _ := p.GetInt64("limit", 10)
	page, _ := p.GetInt64("page", 1)
	managerId := common.GetManagerId(p.GetSession("token"))
	var r *models.Role
	status, list := r.GetRoleList(limit, page, managerId)
	if status == true {
		p.Success(0, list, "获取成功")
	}
	p.Error(400, "拉取列表失败")
}

// CreateRoleUser @Title 创建角色
// @router /createRoleUser [post]
func (p *ProjectController) CreateRoleUser() {
	p.IsDeveloper()
	var nameList []string
	title := p.GetString("title", "")
	description := p.GetString("description", "")
	name := p.GetString("name", "")
	if name != "" {
		err := json.Unmarshal([]byte(name), &nameList)
		if err != nil {
			p.Error(400, "提供的权限参数格式错误")
		}
	}
	managerId := common.GetManagerId(p.GetSession("token"))
	r := models.Role{
		ManagerId:   managerId,
		Title:       title,
		Description: description,
	}
	id := r.Add(nameList)
	if id > 0 {
		p.Success(0, id, "创建成功")
	}
	p.Error(400, "创建失败")
}

// UpdateRoleUser @Title 修改角色
// @router /updateRoleUser [post]
func (p *ProjectController) UpdateRoleUser() {
	p.IsDeveloper()
	var nameList []string
	id, _ := p.GetInt("id", 0)
	title := p.GetString("title", "")
	description := p.GetString("description", "")
	name := p.GetString("name", "")
	if name != "" {
		err := json.Unmarshal([]byte(name), &nameList)
		if err != nil {
			p.Error(400, "提供的权限参数格式错误")
		}
	}
	r := models.Role{
		ID:          id,
		Title:       title,
		Description: description,
	}
	row := r.Update(nameList)
	if row > 0 {
		p.Success(0, row, "更新成功")
	}
	p.Error(400, "更新失败")
}

// GetUserRole @Title 获取角色权限
// @router /getUserRole [post]
func (p *ProjectController) GetUserRole() {
	p.IsDeveloper()
	roleId, _ := p.GetInt("role_id", 0)
	r := models.RoleItem{RoleId: roleId}
	list := r.GetUserRole()
	var result []string
	for _, i := range list {
		result = append(result, i.Name)
	}
	p.Success(0, result, "获取成功")
}

// GetRole @Title 获取所有角色
// @router /getRole [post]
func (p *ProjectController) GetRole() {
	p.IsDeveloper()
	r := models.Role{}
	managerId := common.GetManagerId(p.GetSession("token"))
	list := r.GetRoleAll(managerId)
	p.Success(0, list, "获取成功")
}

// DeleteRole @Title 删除角色权限
// @router /deleteRole [post]
func (p *ProjectController) DeleteRole() {
	p.IsDeveloper()
	roleId, _ := p.GetInt("id", 0)
	r := models.Role{ID: roleId}
	row := r.Delete()
	if row > 0 {
		p.Success(0, row, "删除成功")
	}
	p.Error(400, "删除失败")
}

// UpdateCache @Title 更新缓存
// @router /updateCache [post]
func (p *ProjectController) UpdateCache() {
	status, ac := common.GetCacheAC()
	if status == false {
		logs.Error("缓存初始化失败")
	}
	var managerIdList []int
	var upManagerIdList []int
	m := models.Manager{ID: p.ManagerId}
	m.GetManagerIdList(p.ManagerId, &managerIdList)
	for _, i := range managerIdList {
		_ = ac.Delete("manager-" + strconv.Itoa(i))
	}
	models.GetUpManagerList(p.ManagerId, &upManagerIdList)
	for _, i := range upManagerIdList {
		_ = ac.Delete("manager-" + strconv.Itoa(i))
	}
	if p.ManagerInfo.Pid == 0 {
		project := models.Project{}
		projectList := project.AllProject()
		for _, i := range projectList {
			data, _ := json.Marshal(&i)
			_ = ac.Put(common.GetStringMd5(i.AppKey), string(data), 365*24*60*60*time.Second)
			_ = ac.Put(i.AppKey, string(data), 365*24*60*60*time.Second)
		}
		projectLogin := models.ProjectLogin{}
		loginList := projectLogin.AllLogin()
		for _, i := range loginList {
			data, _ := json.Marshal(&i)
			_ = ac.Put("cache-project-login-"+strconv.Itoa(i.ID), string(data), 365*24*60*60*time.Second)
		}
		projectVersion := models.ProjectVersion{}
		versionList := projectVersion.AllVersion()
		for _, i := range versionList {
			var s strings.Builder
			s.WriteString(strconv.Itoa(i.ProjectId))
			s.WriteString("-")
			s.WriteString(common.GetVersionString(i.Version))
			data, _ := json.Marshal(&i)
			_ = ac.Put(s.String(), string(data), 365*24*60*60*time.Second)
		}
	}
	p.Success(0, 1, "缓存更新成功")
}

// GetOrder @Title 获取订单列表
// @router /getOrder [post]
func (p *ProjectController) GetOrder() {
	limit, _ := p.GetInt64("limit", 10)
	page, _ := p.GetInt64("page", 1)
	var o *models.Order
	status, list :=o.GetOrderList(limit, page, p.ManagerIdArr)
	if status == true {
		p.Success(0, list, "获取成功")
	}
	p.Error(400, "拉取列表失败")
}

// GetAgent @Title 获取全部代理列表
// @router /getAgent [post]
func (p *ProjectController) GetAgent() {
	var m *models.Manager
	list := m.GetAllAgent(p.ManagerIdArr)
	p.Success(0, list, "获取成功")
}