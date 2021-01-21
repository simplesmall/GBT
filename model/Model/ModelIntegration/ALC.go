package ModelIntegration

import (
	"GBT/Config"
	"GBT/model/Model/RBAC"
	"fmt"
)

//type RBACIntegration struct {
//	Role *RBAC.Role
//	Menu *RBAC.Menu
//	Admins *RBAC.Admins
//	RoleMenu *RBAC.RoleMenu
//	AdminsRole *RBAC.AdminsRole
//}
//
//var (
//	Role *RBAC.Role
//	Menu *RBAC.Menu
//	Admins *RBAC.Admins
//	RoleMenu *RBAC.RoleMenu
//	AdminsRole *RBAC.AdminsRole
//)

func AutoMigrateRBAC()  {
	err:=Config.DB.AutoMigrate(&RBAC.Menu{},&RBAC.Admins{},&RBAC.RoleMenu{},&RBAC.Role{},&RBAC.AdminsRole{}).Error
	if err != nil {
		fmt.Println("AutoMigrate is wrong... : ",err)
		return
	}
}
