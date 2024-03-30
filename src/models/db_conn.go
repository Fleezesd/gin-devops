package models

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/fleezesd/gin-devops/src/common"
	"github.com/fleezesd/gin-devops/src/config"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Db       *gorm.DB
	Enforcer *casbin.Enforcer
)

func InitDB(sc *config.ServerConfig) (err error) {
	Db, err = gorm.Open(mysql.Open(sc.Mysql.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}
	if err = Db.Use(otelgorm.NewPlugin()); err != nil {
		return err
	}
	return nil
}

func InitCasbin(sc *config.ServerConfig) (err error) {
	adapter, err := gormadapter.NewAdapterByDB(Db)
	if err != nil {
		sc.Logger.Error("Casbin适配器初始化失败", zap.Error(err))
		return err
	}
	Enforcer, err = casbin.NewEnforcer("./config/rbac_model.conf", adapter)
	if err != nil {
		sc.Logger.Error("Casbin创建Enforcer失败", zap.Error(err))
		return err
	}
	return nil
}

// MigrateTable 自动建表的逻辑
func MigrateTable() error {
	return Db.AutoMigrate(
		&User{},
		&Role{},
		&Menu{},
		&Api{},
	)
}

func MockUserRegister(sc *config.ServerConfig) {
	// 菜单
	menus := []*Menu{
		{
			Name:      "System",
			Title:     "系统管理",
			Icon:      "ion:settings-outline",
			Type:      "0",
			Show:      "1",
			OrderNo:   1,
			Component: "LAYOUT",
			Redirect:  "/system/account",
			Path:      "/system",
		},
		{
			Name:      "Permission",
			Title:     "权限管理",
			Icon:      "ion:layers-outline",
			Type:      "0",
			Show:      "1",
			OrderNo:   2,
			Component: "LAYOUT",
			Path:      "/permission",
			Redirect:  "/permission/front/page",
		},
		{
			Name:      "ServiceTree",
			Title:     "服务树与cmdb",
			Icon:      "ion:layers-outline",
			Type:      "0",
			Show:      "1",
			OrderNo:   3,
			Component: "LAYOUT",
			Path:      "/serviceTree",
			Redirect:  "/serviceTree/serviceTree/index",
		},
		{
			Name:      "WorkOrder",
			Title:     "工单系统",
			Icon:      "ion:git-compare-outline",
			Type:      "0",
			Show:      "1",
			OrderNo:   4,
			Component: "LAYOUT",
			Path:      "/workOrder",
			Redirect:  "/workOrder/process/index",
		},

		{
			Name:      "MenuManagement",
			Title:     "菜单管理",
			Icon:      "ant-design:account-book-filled",
			Type:      "1",
			Show:      "1",
			OrderNo:   1,
			Component: "/demo/system/menu/index",
			Pid:       1,
			Path:      "menu",
		},
		{
			Name:      "AccountManagement",
			Title:     "用户管理",
			Icon:      "ant-design:account-book-twotone",
			Type:      "1",
			Show:      "1",
			OrderNo:   2,
			Component: "/demo/system/account/index",
			Pid:       1,
			Path:      "account",
		},
		{
			Name:      "RoleManagement",
			Title:     "角色管理",
			Icon:      "ion:layers-outline",
			Type:      "1",
			Show:      "1",
			OrderNo:   3,
			Component: "/demo/system/role/index",
			Pid:       1,
			Path:      "role",
		},
		{
			Name:      "ChangePassword",
			Title:     "修改密码",
			Icon:      "ion:layers-outline",
			Type:      "1",
			Show:      "1",
			OrderNo:   4,
			Component: "/demo/system/password/index",
			Pid:       1,
			Path:      "changePassword",
		},

		{
			Name:      "ApiManagement",
			Title:     "api接口管理",
			Icon:      "ant-design:account-book-filled",
			Type:      "1",
			Show:      "1",
			OrderNo:   5,
			Component: "/demo/system/api/index",
			Pid:       1,
			Path:      "api",
		},

		{
			Name:      "PermissionFrontDemo",
			Title:     "前端权限管理",
			Icon:      "ion:layers-outline",
			Type:      "1",
			Show:      "1",
			OrderNo:   1,
			Component: "/demo/permission/front/index",
			Pid:       2,
			Path:      "front",
		},

		{
			Name:      "ServiceTreeIndex",
			Title:     "服务树",
			Icon:      "ant-design:account-book-filled",
			Type:      "1",
			Show:      "1",
			OrderNo:   1,
			Component: "/stree/stree/index",
			Pid:       3,
			Path:      "stree",
		},
		{
			Name:      "ServiceTreeIndexAsync",
			Title:     "服务树异步",
			Icon:      "ant-design:account-book-filled",
			Type:      "1",
			Show:      "1",
			OrderNo:   2,
			Component: "/stree/stree/indexAsync",
			Pid:       3,
			Path:      "streeAsync",
		},
		{
			Name:      "ProcessManagement",
			Title:     "流程管理",
			Icon:      "ant-design:account-book-filled",
			Type:      "1",
			Show:      "1",
			OrderNo:   1,
			Component: "/workorder/process/index",
			Pid:       4,
			Path:      "process",
		},
		{
			Name:      "FormDesignManagement",
			Title:     "表单设计管理",
			Icon:      "ant-design:account-book-filled",
			Type:      "1",
			Show:      "1",
			OrderNo:   3,
			Component: "/workorder/formDesign/index",
			Pid:       4,
			Path:      "formDesign",
		},
		{
			Name:      "WorkOrderTemplateManagement",
			Title:     "工单模板管理",
			Icon:      "ant-design:account-book-filled",
			Type:      "1",
			Show:      "1",
			OrderNo:   4,
			Component: "/workorder/template/index",
			Pid:       4,
			Path:      "template",
		},
		{
			Name:      "WorkOrderApply",
			Title:     "工单申请",
			Icon:      "ant-design:account-book-filled",
			Type:      "1",
			Show:      "1",
			OrderNo:   4,
			Component: "/workorder/apply/index",
			Pid:       4,
			Path:      "apply",
		},
		{
			Name:      "WorkOrderCreate",
			Title:     "工单创建",
			Icon:      "ant-design:account-book-filled",
			Type:      "1",
			Show:      "0",
			OrderNo:   4,
			Component: "/workorder/apply/create",
			Pid:       4,
			Path:      "create",
		},
		{
			Name:      "WorkOrderSearch",
			Title:     "工单查询",
			Icon:      "simple-icons:about-dot-me",
			Type:      "1",
			Show:      "1",
			OrderNo:   4,
			Component: "/workorder/apply/search",
			Pid:       4,
			Path:      "search",
		},
		{
			Name:      "WorkOrderDetail",
			Title:     "工单详情",
			Icon:      "ant-design:account-book-filled",
			Type:      "1",
			Show:      "0",
			OrderNo:   4,
			Component: "/workorder/detail/index",
			Pid:       4,
			Path:      "detail",
		},
	}
	// api
	apis := []*Api{
		{
			Path:   "/api/system/menu",
			Method: "GET",
			Title:  "系统管理-菜单相关",
			Type:   "0",
		},

		{
			Path:   "/api/*",
			Method: "ALL",
			Title:  "api的所有权限",
			Type:   "0",
		},
		{
			Path:   "/api/system/account",
			Method: "GET",
			Title:  "系统管理-用户相关",
			Type:   "0",
		},
		{
			Path:   "/api/system/getMenuList",
			Method: "GET",
			Pid:    1,
			Title:  "系统管理-根据用户获取菜单",
			Type:   "1",
		},
		{
			Path:   "/api/system/getMenuListAll",
			Method: "GET",
			Pid:    1,
			Title:  "系统管理-获取全量菜单",
			Type:   "1",
		},
		{
			Path:   "/api/system/updateMenu",
			Method: "POST",
			Pid:    1,
			Title:  "系统管理-更新菜单",
			Type:   "1",
		},
		{
			Path:   "/api/system/createMenu",
			Method: "POST",
			Pid:    1,
			Title:  "系统管理-创建菜单",
			Type:   "1",
		},
		{
			Path:   "/api/system/deleteMenu/:id",
			Method: "DELETE",
			Pid:    1,
			Title:  "系统管理-删除菜单",
			Type:   "1",
		},
		{
			Path:   "/api/getUserInfo",
			Method: "GET",
			Pid:    3,
			Title:  "获取用户信息",
			Type:   "1",
		},
		{
			Path:   "/api/getPermCode",
			Method: "GET",
			Pid:    3,
			Title:  "获取用户code",
			Type:   "1",
		},

		{
			Path:   "/api/system/getAccountList",
			Method: "GET",
			Pid:    3,
			Title:  "获取用户列表",
			Type:   "1",
		},

		{
			Path:   "/api/*",
			Method: "GET",
			Pid:    2,
			Title:  "所有api的GET权限",
			Type:   "1",
		},
		{
			Path:   "/api/*",
			Method: "POST",
			Pid:    2,
			Title:  "所有api的新增权限",
			Type:   "1",
		},
		{
			Path:   "/api/*",
			Method: "DELETE",
			Pid:    2,
			Title:  "所有api的DELETE权限",
			Type:   "1",
		},
	}
	// 用户
	user1 := User{
		Username: "admin",
		Password: "123456",
		RealName: "超管",
		Desc:     "",
		HomePath: "/system/account",
		Enable:   1,
		Roles: []*Role{
			{
				RoleName:  "超级管理员",
				RoleValue: "super",
				Menus:     menus,
				Apis:      apis,
			},
		},
	}

	user1.Password = common.BcryptHash(user1.Password)
	// 保存用户信息   关联role menu会自动migrate到各自表中
	err := Db.FirstOrCreate(&user1).Error
	if err != nil {
		sc.Logger.Error("模拟用户注册失败", zap.Error(err))
		return
	}

	err = Db.Create(apis).Error
	sc.Logger.Info("创建api结果", zap.Any("err", err))

	dbRole, _ := GetRoleByRoleValue("super")
	dbRole.Apis = apis
	err = dbRole.UpdateApis(apis)
	sc.Logger.Info("更新api结果", zap.Any("err", err))
	sc.Logger.Info("模拟用户注册成功!")

	// mock casbin的策略
	for _, api := range apis {
		_, err = Enforcer.AddPolicy("super", api.Path, api.Method)
		if err != nil {
			sc.Logger.Error("添加权限失败",
				zap.Error(err),
				zap.String("角色", "super"),
				zap.String("api.Path", api.Path),
				zap.String("api.Method", api.Method),
				zap.String("api.Title", api.Title),
			)
		}
	}
	sc.Logger.Info("模拟casbin策略添加成功!")
}
