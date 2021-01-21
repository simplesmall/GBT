package common

import (
	"GBT/common/EnforcerInstance"
	"GBT/common/Utils/Convert"
	"GBT/model/Model/Sys"
	"GBT/model/Model/common"
	"github.com/casbin/casbin"
)

const (
	PrefixUserID = "u"
	PrefixRoleID = "r"
)

var Enforcer *casbin.Enforcer

// 角色-URL导入
func InitCsbinEnforcer() (err error) {
	var enforcer *casbin.Enforcer
	enforcer, err = EnforcerInstance.EnforcerInstance()

	if err != nil {
		return
	}
	var roles []Sys.Role
	err = common.Find(&Sys.Role{}, &roles)
	if err != nil {
		return
	}
	if len(roles) == 0 {
		Enforcer = enforcer
		return
	}
	for _, role := range roles {
		setRolePermission(enforcer, role.ID)
	}
	Enforcer = enforcer
	return
}

// 删除角色
func CsbinDeleteRole(roleids []uint64) {
	if Enforcer == nil {
		return
	}
	for _, rid := range roleids {
		Enforcer.DeletePermissionsForUser(PrefixRoleID + Convert.ToString(rid))
		Enforcer.DeleteRole(PrefixRoleID + Convert.ToString(rid))
	}
}

// 设置角色权限
func CsbinSetRolePermission(roleid uint64) {
	if Enforcer == nil {
		return
	}
	Enforcer.DeletePermissionsForUser(PrefixRoleID + Convert.ToString(roleid))
	setRolePermission(Enforcer, roleid)
}

// 设置角色权限
func setRolePermission(enforcer *casbin.Enforcer, roleid uint64) {
	var rolemenus []Sys.RoleMenu
	err := common.Find(&Sys.RoleMenu{RoleID: roleid}, &rolemenus)
	if err != nil {
		return
	}
	for _, rolemenu := range rolemenus {
		menu := Sys.Menu{}
		where := Sys.Menu{}
		where.ID = rolemenu.MenuID
		_, err = common.First(&where, &menu)
		if err != nil {
			return
		}
		if menu.MenuType == 3 {
			enforcer.AddPermissionForUser(PrefixRoleID+Convert.ToString(roleid), "/api"+menu.URL, "GET|POST")
		}
	}
}

// 检查用户是否有权限
func CsbinCheckPermission(userID, url, methodtype string) (bool, error) {
	return Enforcer.EnforceSafe(PrefixUserID+userID, url, methodtype)
}

// 用户角色处理
func CsbinAddRoleForUser(userid uint64)(err error){
	if Enforcer == nil {
		return
	}
	uid:=PrefixUserID+Convert.ToString(userid)
	Enforcer.DeleteRolesForUser(uid)
	var adminsroles []Sys.AdminsRole
	err = common.Find(&Sys.AdminsRole{AdminsID: userid}, &adminsroles)
	if err != nil {
		return
	}
	for _, adminsrole := range adminsroles {
		Enforcer.AddRoleForUser(uid, PrefixRoleID+Convert.ToString(adminsrole.RoleID))
	}
	return
}
