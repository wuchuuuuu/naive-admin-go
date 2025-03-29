package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"naive-admin-go/db"
	"naive-admin-go/inout"
	"naive-admin-go/model"
	"naive-admin-go/utils"
	"strconv"
)

var Permissions = &permissions{}

type permissions struct {
}

func (permissions) GetBts(c *gin.Context) {
	id := c.Param("id")
	// 先根据id找到parentId，再根据parentId找到所有子节点，再根据子节点找到所有按钮,返回一个数组bts
	var partId int
	db.Dao.Model(model.Permission{}).Where("id = ?", id).Select("parentId").Find(&partId)
	var bts []model.Permission
	db.Dao.Model(model.Permission{}).Where("parentId = ?", partId).Where("type = 'BUTTON'").Find(&bts)
	utils.Resp.Succ(c, bts)
}

func (permissions) List(c *gin.Context) {
	var onePermissList = make([]model.Permission, 0)
	db.Dao.Model(model.Permission{}).Where("parentId is NULL").Order("`order` Asc").Find(&onePermissList)
	for i, perm := range onePermissList {
		var twoPerissList []model.Permission
		db.Dao.Model(model.Permission{}).Where("parentId = ?", perm.ID).Order("`order` Asc").Find(&twoPerissList)
		for i2, perm2 := range twoPerissList {
			var twoPerissList2 []model.Permission
			db.Dao.Model(model.Permission{}).Where("parentId = ?", perm2.ID).Order("`order` Asc").Find(&twoPerissList2)
			twoPerissList[i2].Children = twoPerissList2
		}
		onePermissList[i].Children = twoPerissList
	}

	utils.Resp.Succ(c, onePermissList)
}

func (permissions) ListPage(c *gin.Context) {
	var data = &inout.RoleListPageRes{}
	var name = c.DefaultQuery("name", "")
	var pageNoReq = c.DefaultQuery("pageNo", "1")
	var pageSizeReq = c.DefaultQuery("pageSize", "10")
	pageNo, _ := strconv.Atoi(pageNoReq)
	pageSize, _ := strconv.Atoi(pageSizeReq)
	orm := db.Dao.Model(model.Role{})
	if name != "" {
		orm = orm.Where("name like ?", "%"+name+"%")
	}
	orm.Count(&data.Total)

	orm.Offset((pageNo - 1) * pageSize).Limit(pageSize).Find(&data.PageData)
	for i, datum := range data.PageData {
		var perIdList []int64
		db.Dao.Model(model.RolePermissionsPermission{}).Where("roleId=?", datum.ID).Select("permissionId").Find(&perIdList)
		data.PageData[i].PermissionIds = perIdList
	}
	utils.Resp.Succ(c, data)
}
func (permissions) Add(c *gin.Context) {
	var params inout.AddPermissionReq
	err := c.Bind(&params)
	if err != nil {
		utils.Resp.Err(c, 20001, err.Error())
		return
	}

	err = db.Dao.Model(model.Permission{}).Create(&model.Permission{
		Name:      params.Name,
		Code:      params.Code,
		Type:      params.Type,
		ParentId:  params.ParentId, // insert value null
		Path:      params.Path,
		Icon:      params.Icon,
		Component: params.Component,
		Layout:    params.Layout,
		KeepAlive: params.KeepAlive,
		Show:      params.Show,
		Enable:    params.Enable,
		Order:     params.Order,
	}).Error
	if err != nil {
		utils.Resp.Err(c, 20001, err.Error())
		return
	}
	utils.Resp.Succ(c, "")
}
func (permissions) Delete(c *gin.Context) {
	id := c.Param("id")
	err := db.Dao.Transaction(func(tx *gorm.DB) error {
		tx.Where("id =?", id).Delete(&model.Permission{})
		tx.Where("permissionId =?", id).Delete(&model.RolePermissionsPermission{})
		return nil
	})
	if err != nil {
		utils.Resp.Err(c, 20001, err.Error())
		return
	}
	utils.Resp.Succ(c, "")
}
func (permissions) PatchPermission(c *gin.Context) {

	var params inout.PatchPermissionReq
	err := c.BindJSON(&params)
	if err != nil {
		utils.Resp.Err(c, 20001, err.Error())
		return
	}
	err = db.Dao.Model(model.Permission{}).Where("id=?", params.Id).Select(
		"Name", "Code", "Enable",
	).Updates(model.Permission{
		Name:      params.Name,
		Code:      params.Code,
		Type:      params.Type,
		Path:      params.Path,
		Icon:      params.Icon,
		Component: params.Component,
		Layout:    params.Layout,
		KeepAlive: params.KeepAlive,
		Method:    params.Method,
		Show:      params.Show,
		Enable:    params.Enable,
		Order:     params.Order,
	}).Error
	if err != nil {
		utils.Resp.Err(c, 20001, err.Error())
		return
	}
	utils.Resp.Succ(c, params)

}
func IsTrue(v bool) int {
	if v {
		return 1
	}
	return 0
}
