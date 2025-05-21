package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "AddManagerCards",
            Router: `/addManagerCards`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "BatchKeys",
            Router: `/batchKeys`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "BatchMember",
            Router: `/batchMember`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "BindProjectLogin",
            Router: `/bindProjectLogin`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "CardList",
            Router: `/cardList`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "CreateCard",
            Router: `/createCard`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "CreateKeys",
            Router: `/createKeys`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "CreateLoginRule",
            Router: `/createLoginRule`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "CreateProject",
            Router: `/createProject`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "CreateRoleUser",
            Router: `/createRoleUser`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "CreateVersion",
            Router: `/createVersion`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "DeleteCard",
            Router: `/deleteCard`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "DeleteKeys",
            Router: `/deleteKeys`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "DeleteLoginRule",
            Router: `/deleteLoginRule`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "DeleteManagerCards",
            Router: `/deleteManagerCards`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "DeleteMember",
            Router: `/deleteMember`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "DeleteProject",
            Router: `/deleteProject`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "DeleteProjectVersion",
            Router: `/deleteProjectVersion`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "DeleteRole",
            Router: `/deleteRole`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "GetAgent",
            Router: `/getAgent`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "GetCardList",
            Router: `/getCardList`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "GetKeysList",
            Router: `/getKeysList`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "GetLoginRuleList",
            Router: `/getLoginRuleList`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "GetManagerCards",
            Router: `/getManagerCards`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "GetManagerList",
            Router: `/getManagerList`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "GetMemberList",
            Router: `/getMemberList`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "GetOnlineList",
            Router: `/getOnlineList`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "GetOrder",
            Router: `/getOrder`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "GetProjectList",
            Router: `/getProjectList`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "GetRole",
            Router: `/getRole`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "GetRoleUser",
            Router: `/getRoleUser`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "GetUserRole",
            Router: `/getUserRole`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "GetVersionList",
            Router: `/getVersionList`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "LockKeys",
            Router: `/lockKey`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "LockMember",
            Router: `/lockMember`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "LoginRuleList",
            Router: `/loginRuleList`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "ManagerAdd",
            Router: `/managerAdd`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "ManagerDelete",
            Router: `/managerDelete`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "ManagerIdList",
            Router: `/managerIdList`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "ManagerUpdate",
            Router: `/managerUpdate`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "MemberLogout",
            Router: `/memberLogout`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "ProjectList",
            Router: `/projectList`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "QueryOrderRmb",
            Router: `/queryOrderRmb`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "RolePower",
            Router: `/rolePower`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "UnbindMember",
            Router: `/unbindMember`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "UpdateCache",
            Router: `/updateCache`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "UpdateCard",
            Router: `/updateCard`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "UpdateLoginRule",
            Router: `/updateLoginRule`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "UpdateManagerCards",
            Router: `/updateManagerCards`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "UpdateMember",
            Router: `/updateMember`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "UpdateProject",
            Router: `/updateProject`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "UpdateProjectVersion",
            Router: `/updateProjectVersion`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:ProjectController"],
        beego.ControllerComments{
            Method: "UpdateRoleUser",
            Router: `/updateRoleUser`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:UserController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:UserController"],
        beego.ControllerComments{
            Method: "GetInfo",
            Router: `/getInfo`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:UserController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:UserController"],
        beego.ControllerComments{
            Method: "GetKeysEcharts",
            Router: `/getKeysEcharts`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:UserController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:UserController"],
        beego.ControllerComments{
            Method: "GetMemberEcharts",
            Router: `/getMemberEcharts`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:UserController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:UserController"],
        beego.ControllerComments{
            Method: "GetSysNotice",
            Router: `/getSysNotice`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:UserController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:UserController"],
        beego.ControllerComments{
            Method: "Info",
            Router: `/info`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:UserController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:UserController"],
        beego.ControllerComments{
            Method: "Login",
            Router: `/login`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:UserController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:UserController"],
        beego.ControllerComments{
            Method: "Logout",
            Router: `/logout`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/admin:UserController"] = append(beego.GlobalControllerRouter["verification/controllers/admin:UserController"],
        beego.ControllerComments{
            Method: "Update",
            Router: `/update`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers/api:IndexController"] = append(beego.GlobalControllerRouter["verification/controllers/api:IndexController"],
        beego.ControllerComments{
            Method: "Index",
            Router: `/`,
            AllowHTTPMethods: []string{"post","get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["verification/controllers:MainController"] = append(beego.GlobalControllerRouter["verification/controllers:MainController"],
        beego.ControllerComments{
            Method: "Index",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
