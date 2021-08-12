/**
 * @Author: Anpw
 * @Description:
 * @File:  HousingInfoController
 * @Version: 1.0.0
 * @Date: 2021/7/23 21:15
 */

package controllers

import (
	"PerInfoChain/models"
	"PerInfoChain/pkg/app"
	"PerInfoChain/pkg/convert"
	"PerInfoChain/pkg/errcode"
	"PerInfoChain/pkg/upload"
	"PerInfoChain/pkg/utils"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	"time"
)

type HousingInfoController struct {
	BaseController
}
type AddHousingInfoRequest struct {
	UserID    string  `alias:"个人id" valid:"Required;"`
	Address   string  `alias:"地址" valid:"Required;"`
	HouseSize uint    `alias:"房屋大小" valid:"Required;"`
	Price     float32 `alias:"购买价格" valid:"Required;"`
	HousingNo string  `alias:"房产证号" valid:"Required;"`
	PayTime   string  `alias:"购买时间" valid:"Required;"`
	UsingYear uint    `alias:"使用年限" valid:"Required;"`
}

type UpdateHousingInfoRequest struct {
	Address   string  `alias:"地址" valid:""`
	HouseSize uint    `alias:"房屋大小" valid:""`
	Price     float32 `alias:"购买价格" valid:""`
	HousingNo string  `alias:"房产证号" valid:""`
	PayTime   string  `alias:"购买时间" valid:""`
	UsingYear uint    `alias:"使用年薪" valid:""`
}

// AddHousingInfo 添加住房基本信息
func (c *HousingInfoController) AddHousingInfo() {
	//获取文件
	file, fileHeader, err := c.GetFile("file")
	fileType := convert.StrTo(c.GetString("type")).MustInt()
	if err != nil {
		c.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}
	if fileHeader == nil || fileType <= 0 {
		c.ToErrorResponse(errcode.InvalidParams)
		return
	}

	housingInfo := AddHousingInfoRequest{}
	if err := c.ParseForm(&housingInfo); err != nil {
		logs.Error("ParseForm err: %v", err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	valid, errs := app.BindAndValid(&housingInfo)
	if !valid {
		logs.Error("app.BindAndValid err: %v", errs)
		c.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//判断用户是否存在
	var fabricUserInfo models.FabricUserInfo
	if err := models.DB.Where("user_id=?", housingInfo.UserID).First(&fabricUserInfo).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if fabricUserInfo.UserID == "" {
		c.ToErrorResponse(errcode.ErrorAddHouseInfoFail.WithDetails(errors.New("用户不存在").Error()))
		return
	}
	fileInfo, err := upload.Upload(upload.FileType(fileType), file, fileHeader, "houseinfo")
	if err != nil {
		logs.Error("UploadFile err:", err)
		c.ToErrorResponse(errcode.ERROR_UPLOAD_FILE_FAIL.WithDetails(err.Error()))
		return
	}
	u := models.HousingInfo{
		HousingID:               utils.GenerateUID(),
		UserID:                  housingInfo.UserID,
		LoginUserID:             fabricUserInfo.LoginUserID,
		Address:                 housingInfo.Address,
		HouseSize:               housingInfo.HouseSize,
		Price:                   housingInfo.Price,
		HousingNo:               housingInfo.HousingNo,
		PayTime:                 housingInfo.PayTime,
		UsingYear:               housingInfo.UsingYear,
		HouseInfoFilePath:       fileInfo.Dst,
		HouseInfoFileHash:       fileInfo.Hash,
		HouseInfoFileUploadTime: time.Now(),
	}
	if err := models.DB.Create(&u).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}

	c.ToResponse(map[string]interface{}{
		"msg": "添加成功",
	})
	return
}

// HousingInfoList 获取列表信息 以及 查询
func (c *HousingInfoController) HousingInfoList() {
	page, _ := c.GetInt("page")
	pageSize, _ := c.GetInt("page_size")

	var count int
	if err := models.DB.Table(models.CriminalCase{}.TableName()).Count(&count).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}

	pager := Pager{Page: app.GetPage(page), PageSize: app.GetPageSize(pageSize), TotalRows: count}
	pageOffset := app.GetPageOffset(pager.Page, pager.PageSize)
	if pageOffset >= 0 && pageSize > 0 {
		models.DB = models.DB.Offset(pageOffset).Limit(pageSize)
	}

	var housingInfos []*models.HousingInfo
	if err := models.DB.Find(&housingInfos).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponseList(housingInfos, &pager)
	return
}

// GetHousingInfo 查询
func (c *HousingInfoController) GetHousingInfo() {
	var userInfo models.FabricUserInfo
	var housingInfo []*models.HousingInfo
	IDCard := c.Ctx.Input.Param(":id_card")
	err := models.DB.Where("id_card=?", IDCard).First(&userInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	err = models.DB.Where("user_id=?", userInfo.UserID).First(&housingInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(housingInfo)
	return
}

func (c *HousingInfoController) UpdateHousingInfo() {
	housingID := c.Ctx.Input.Param(":housing_id")
	housingInfo := UpdateHousingInfoRequest{}
	if err := c.ParseForm(&housingInfo); err != nil {
		logs.Error("ParseForm err: %v", err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	valid, errs := app.BindAndValid(&housingInfo)
	if !valid {
		logs.Error("app.BindAndValid err: %v", errs)
		c.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	values := map[string]interface{}{}
	if housingInfo.HousingNo != "" {
		values["housing_no"] = housingInfo.HousingNo
	}
	if housingInfo.HouseSize != 0 {
		values["house_size"] = housingInfo.HouseSize
	}
	if housingInfo.Address != "" {
		values["address"] = housingInfo.Address
	}
	if housingInfo.UsingYear != 0 {
		values["using_year"] = housingInfo.UsingYear
	}
	if housingInfo.PayTime != "" {
		values["pay_time"] = housingInfo.PayTime
	}
	if housingInfo.Price != 0 {
		values["price"] = housingInfo.Price
	}

	var housingInfos models.HousingInfo
	if err := models.DB.Model(&housingInfos).Where("housing_id=?", housingID).Updates(values).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "更新成功",
	})
	return
}

func (c *HousingInfoController) DeleteHousingInfo() {
	HousingID := c.Ctx.Input.Param(":housing_id")
	var housingInfo models.HousingInfo
	err := models.DB.Where("housing_id=?", HousingID).First(&housingInfo).Error
	if err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	err = upload.Delete(housingInfo.HouseInfoFilePath)
	if err != nil{
		c.ToErrorResponse(errcode.ERROR_DELETE_FILE_FAIL.WithDetails(err.Error()))
		return
	}
	if err = models.DB.Delete(&housingInfo).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "删除成功",
	})
}


func (c *HousingInfoController) UpdateFile() {
	userID := c.Ctx.Input.Param(":user_id")
	var houseInfo models.HousingInfo
	err := models.DB.Where("user_id=?", userID).First(&houseInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if houseInfo.UserID == "" {
		c.ToErrorResponse(errcode.ERROR_UPDATE_FILE_FAIL.WithDetails(errors.New("用户不存在").Error()))
		return
	}

	//获取文件
	file, fileHeader, err := c.GetFile("file")
	fileType := convert.StrTo(c.GetString("type")).MustInt()
	if err != nil {
		c.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}
	if fileHeader == nil || fileType <= 0 {
		c.ToErrorResponse(errcode.InvalidParams)
		return
	}
	fileInfo, err := upload.Update(upload.FileType(fileType), file, fileHeader, houseInfo.HouseInfoFilePath)
	if err != nil {
		logs.Error("UploadFile err:", err)
		c.ToErrorResponse(errcode.ERROR_UPDATE_FILE_FAIL.WithDetails(err.Error()))
		return
	}
	values := map[string]interface{}{
		"house_info_file_hash" : fileInfo.Hash,
		"house_info_file_path" : fileInfo.Dst,
		"house_info_file_upload_time" : time.Now(),
	}

	if err := models.DB.Model(&houseInfo).Where("user_id=?", userID).Updates(values).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "更新文件成功",
	})
}

func (c *HousingInfoController) Download() {
	userID := c.Ctx.Input.Param(":user_id")
	var houseInfo models.HousingInfo
	err := models.DB.Where("user_id=?", userID).First(&houseInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if houseInfo.UserID == "" {
		c.ToErrorResponse(errcode.ErrorCreateUserFail.WithDetails(errors.New("用户不存在").Error()))
		return
	}
	c.Ctx.Output.Download(houseInfo.HouseInfoFilePath)
	return
}