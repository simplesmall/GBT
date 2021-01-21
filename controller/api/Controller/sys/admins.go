package Sys

import (
	"GBT/common/Utils/Hash"
	"GBT/controller/api/Controller/common"
	"GBT/model/Model/Sys"
	models "GBT/model/Model/common"
	"github.com/gin-gonic/gin"
)

type Admins struct{}

// 分页数据
func (Admins) List(c *gin.Context) {
	page := common.GetPageIndex(c)
	limit := common.GetPageLimit(c)
	sort := common.GetPageSort(c)
	key := common.GetPageKey(c)
	status := common.GetQueryToUint(c, "status")
	var whereOrder []models.PageWhereOrder
	order := "ID DESC"
	if len(sort) >= 2 {
		orderType := sort[0:1]
		order = sort[1:len(sort)]
		if orderType == "+" {
			order += " ASC"
		} else {
			order += " DESC"
		}
	}
	whereOrder = append(whereOrder, models.PageWhereOrder{Order: order})
	if key != "" {
		v := "%" + key + "%"
		var arr []interface{}
		arr = append(arr, v)
		arr = append(arr, v)
		whereOrder = append(whereOrder, models.PageWhereOrder{Where: "user_name like ? or real_name like ?", Value: arr})
	}
	if status > 0 {
		var arr []interface{}
		arr = append(arr, status)
		whereOrder = append(whereOrder, models.PageWhereOrder{Where: "status = ?", Value: arr})
	}
	var total uint64
	list:= []Sys.Admins{}
	err := models.GetPage(&Sys.Admins{}, &Sys.Admins{}, &list, page, limit, &total, whereOrder...)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	common.ResSuccessPage(c, total, &list)
}

// 详情
func (Admins) Detail(c *gin.Context) {
	id := common.GetQueryToUint64(c, "id")
	var model Sys.Admins
	where := Sys.Admins{}
	where.ID = id
	_, err := models.First(&where, &model)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	model.Password = ""
	common.ResSuccess(c, &model)
}

// 更新
func (Admins) Update(c *gin.Context) {
	model := Sys.Admins{}
	err := c.Bind(&model)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	where := Sys.Admins{}
	where.ID = model.ID
	modelOld := Sys.Admins{}
	_, err = models.First(&where, &modelOld)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	model.UserName = modelOld.UserName
	model.Password = modelOld.Password
	err = models.Save(&model)
	if err != nil {
		common.ResFail(c, "操作失败")
		return
	}
	common.ResSuccessMsg(c)
}

//新增
func (Admins) Create(c *gin.Context) {
	model := Sys.Admins{}
	err := c.Bind(&model)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	model.Password = Hash.Md5String(common.MD5_PREFIX + model.Password)
	err = models.Create(&model)
	if err != nil {
		common.ResFail(c, "操作失败")
		return
	}
	common.ResSuccess(c, gin.H{"id": model.ID})
}

// 删除数据
func (Admins) Delete(c *gin.Context) {
	var ids []uint64
	err := c.Bind(&ids)
	if err != nil || len(ids) == 0 {
		common.ResErrSrv(c, err)
		return
	}
	admin:=Sys.Admins{}
  err = admin.Delete(ids)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	common.ResSuccessMsg(c)
}

// 获取用户下的角色ID列表
func (Admins) AdminsRoleIDList(c *gin.Context) {
	adminsid := common.GetQueryToUint64(c, "adminsid")
	roleList := []uint64{}
	where := Sys.AdminsRole{AdminsID: adminsid}
	err := models.PluckList(&Sys.AdminsRole{}, &where, &roleList, "role_id")
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	common.ResSuccess(c, &roleList)
}

// 分配用户角色权限
func (Admins) SetRole(c *gin.Context) {
	adminsid := common.GetQueryToUint64(c, "adminsid")
	var roleids []uint64
	err := c.Bind(&roleids)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	ar := Sys.AdminsRole{}
	err = ar.SetRole(adminsid, roleids)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	go common.CsbinAddRoleForUser(adminsid)
	common.ResSuccessMsg(c)
}


