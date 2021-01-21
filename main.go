package main

import (
	"GBT/Config"
	"GBT/model/Model/ModelIntegration"
	"GBT/model/Model/RBAC"
	"GBT/routers"
	"fmt"
)


func AutoMigrate() {
	err:=Config.DB.AutoMigrate(&RBAC.Menu{},&RBAC.Admins{},&RBAC.RoleMenu{},&RBAC.Role{},&RBAC.AdminsRole{}).Error
	if err != nil {
		fmt.Println("AutoMigrate is wrong... : ",err)
		return
	}
}

func main() {
	defer Config.CloseDB()
	routers.InitServer()
	//AutoMigrate()
	ModelIntegration.AutoMigrateRBAC()
}
