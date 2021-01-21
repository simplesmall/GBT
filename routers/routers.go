package routers

import (
	"GBT/Config"
	MiddleJWT "GBT/common/jwt"
	"GBT/controller/api/Controller/Sys"
	"GBT/controller/api/Controller/common"
	"GBT/model/Model/ModelIntegration"
	"fmt"
	"github.com/casbin/casbin"
	"github.com/casbin/gorm-adapter"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var Enforcer *casbin.Enforcer

// 初始化casbin
func CasbinSetup() {

	// Initialize a gorm adapter with MySQL database.
	a := gormadapter.NewAdapter("mysql", "root:root@tcp(127.0.0.1:3306)/castest",true)

	// Casbin 文本形式创建
	/*m, err := model.NewModelFromString(`
		[request_definition]
		r = sub, obj, act

		[policy_definition]
		p = sub, obj, act

		[policy_effect]
		e = some(where (p.eft == allow))

		[matchers]
		m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
		`)
	if err != nil {
		log.Fatalf("error: model: %s", err)
	}
	*/
	e := casbin.NewEnforcer("./common/Casbin/rbac_models.conf", a)
	Enforcer = e
}

func InitServer() {
	server := gin.Default()
	// 初始化数据库连接
	Config.InitConnect()
	// 自动创建表
	ModelIntegration.AutoMigrateRBAC()
	// 初始化Casbin
	_ = common.InitCsbinEnforcer()

	//初始化casbin
	//server.Use(Authorize())

	// 配置swagger
	server.Use(MiddleJWT.Cors())
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 路由分组
	menu := Sys.Menu{}
	api := server.Group("api")
	{
		// orm练手测试
		api.GET("/login", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"data": "data",
			})
		})
		//增加policy
		api.POST("/api/v1/add", func(c *gin.Context) {
			fmt.Println("增加Policy")
			if ok := Enforcer.AddPolicy("admin", "/api/v1/world", "GET"); !ok {
				fmt.Println("Policy已经存在")
			} else {
				fmt.Println("增加成功")
			}
		})
		//删除policy
		api.DELETE("/api/v1/delete", func(c *gin.Context) {
			fmt.Println("删除Policy")
			if ok := Enforcer.RemovePolicy("admin", "/api/v1/world", "GET"); !ok {
				fmt.Println("Policy不存在")
			} else {
				fmt.Println("删除成功")
			}
		})
		//获取policy
		api.GET("/api/v1/get", func(c *gin.Context) {
			fmt.Println("查看policy")
			list := Enforcer.GetPolicy()
			for _, vlist := range list {
				for _, v := range vlist {
					fmt.Printf("value: %s, ", v)
				}
			}
		})
		// 测试Casbin
		api.GET("/menu/list", menu.List)
		api.GET("/menu/detail", menu.Detail)
		api.GET("/menu/allmenu", menu.AllMenu)
		api.GET("/menu/menubuttonlist", menu.MenuButtonList)
		api.POST("/menu/delete", menu.Delete)
		api.POST("/menu/update", menu.Update)
		api.POST("/menu/create", menu.Create)
		user := Sys.User{}
		api.GET("/user/info", user.Info)
		api.POST("/user/login", user.Login)
		api.POST("/user/logout", user.Logout)
		api.POST("/user/editpwd", user.EditPwd)
		admins := Sys.Admins{}
		api.GET("/admins/list", admins.List)
		api.GET("/admins/detail", admins.Detail)
		api.GET("/admins/adminsroleidlist", admins.AdminsRoleIDList)
		api.POST("/admins/delete", admins.Delete)
		api.POST("/admins/update", admins.Update)
		api.POST("/admins/create", admins.Create)
		api.POST("/admins/setrole", admins.SetRole)
		role := Sys.Role{}
		api.GET("/role/list", role.List)
		api.GET("/role/detail", role.Detail)
		api.GET("/role/rolemenuidlist", role.RoleMenuIDList)
		api.GET("/role/allrole", role.AllRole)
		api.POST("/role/delete", role.Delete)
		api.POST("/role/update", role.Update)
		api.POST("/role/create", role.Create)
		api.POST("/role/setrole", role.SetRole)
	}
	_ = server.Run(":8099")
}

func HandleTest()  {
	Enforcer.AddPermissionForUser("")
}
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		e := Enforcer

		//从DB加载策略
		_ = e.LoadPolicy()

		//获取请求的URI
		obj := c.Request.URL.RequestURI()
		//获取请求方法
		act := c.Request.Method
		//获取用户的角色 应该从db中读取
		sub := "admin"

		//判断策略中是否存在
		if ok := e.Enforce(sub, obj, act); ok {
			fmt.Println("恭喜您,权限验证通过")
			c.Next() // 进行下一步操作
		} else {
			fmt.Println("很遗憾,权限验证没有通过")
			c.Abort()
		}
	}
}
