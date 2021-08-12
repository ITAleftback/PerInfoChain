/**
 * @Author: Anpw
 * @Description:
 * @File:  BankInfoController
 * @Version: 1.0.0
 * @Date: 2021/7/24 18:18
 */

package controllers

import (
	"PerInfoChain/models"
	"PerInfoChain/pkg/app"
	"PerInfoChain/pkg/convert"
	"PerInfoChain/pkg/errcode"
	"PerInfoChain/pkg/upload"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	"time"
)

type BankInfoController struct {
	BaseController
}

type AddBankInfoRequest struct {
	UserID     string  `alias:"个人id" valid:"Required"`
	BankName   string  `alias:"银行名称" valid:"Required"`
	Deposit    float64 `alias:"存款" valid:"Required;"`
	BankCardNo string  `alias:"银行卡号" valid:"Required;"`
}

type UpdateBankInfoRequest struct {
	BankName   string  `alias:"银行名称"`
	Deposit    float64 `alias:"存款"`
}

// AddBankInfo 添加银行用户基本信息
func (c *BankInfoController) AddBankInfo() {
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

	bankInfo := AddBankInfoRequest{}
	if err := c.ParseForm(&bankInfo); err != nil {
		logs.Error("ParseForm err: %v", err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}

	valid, errs := app.BindAndValid(&bankInfo)
	if !valid {
		logs.Error("app.BindAndValid err: %v", errs)
		c.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//判断用户是否已经存在
	var fabricUserInfo models.FabricUserInfo
	if err := models.DB.Where("user_id=?", bankInfo.UserID).First(&fabricUserInfo).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if fabricUserInfo.UserID == "" {
		c.ToErrorResponse(errcode.ErrorAddBankInfoFail.WithDetails(errors.New("用户不存在").Error()))
		return
	}
	var bankInfos []models.BankInfo
	if err := models.DB.Where("bank_card_no=?", bankInfo.BankCardNo).First(&bankInfos).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if len(bankInfos) > 0 {
		c.ToErrorResponse(errcode.ErrorAddBankInfoFail.WithDetails(errors.New("银行卡已添加").Error()))
		return
	}
	fileInfo, err := upload.Upload(upload.FileType(fileType), file, fileHeader, "bankinfo")
	if err != nil {
		logs.Error("UploadFile err:", err)
		c.ToErrorResponse(errcode.ERROR_UPLOAD_FILE_FAIL.WithDetails(err.Error()))
		return
	}
	u := models.BankInfo{
		UserID:                 bankInfo.UserID,
		LoginUserID:            fabricUserInfo.LoginUserID,
		BankCardNo:             bankInfo.BankCardNo,
		BankName:               bankInfo.BankName,
		Deposit:                bankInfo.Deposit,
		BankInfoFilePath:       fileInfo.Dst,
		BankInfoFileHash:       fileInfo.Hash,
		BankInfoFileUploadTime: time.Now(),
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

// BankInfoList 获取用户列表信息 以及 查询
func (c *BankInfoController) BankInfoList() {

	page, _ := c.GetInt("page")
	pageSize, _ := c.GetInt("page_size")

	var count int
	if err := models.DB.Table(models.BankInfo{}.TableName()).Count(&count).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}

	pager := Pager{Page: app.GetPage(page), PageSize: app.GetPageSize(pageSize), TotalRows: count}
	pageOffset := app.GetPageOffset(pager.Page, pager.PageSize)
	if pageOffset >= 0 && pageSize > 0 {
		models.DB = models.DB.Offset(pageOffset).Limit(pageSize)
	}

	var bankInfos []*models.BankInfo
	if err := models.DB.Find(&bankInfos).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponseList(bankInfos, &pager)
	return
}

// GetBankInfo 查询
func (c *BankInfoController) GetBankInfo() {
	var userInfo models.FabricUserInfo
	var bankInfo []*models.BankInfo
	IDCard := c.Ctx.Input.Param(":id_card")

	err := models.DB.Where("id_card=?", IDCard).First(&userInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	err = models.DB.Where("user_id=?", userInfo.UserID).First(&bankInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(bankInfo)
	return
}

func (c *BankInfoController) UpdateBankInfo() {
	bankCardNo := c.Ctx.Input.Param(":bank_card_no")
	bankInfoRequest := UpdateBankInfoRequest{}
	if err := c.ParseForm(&bankInfoRequest); err != nil {
		logs.Error("ParseForm err: %v", err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	valid, errs := app.BindAndValid(&bankInfoRequest)
	if !valid {
		logs.Error("app.BindAndValid err: %v", errs)
		c.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	values := map[string]interface{}{}
	if bankInfoRequest.BankName != "" {
		values["bank_name"] = bankInfoRequest.BankName
	}
	if bankInfoRequest.Deposit != 0 {
		values["deposit"] = bankInfoRequest.Deposit
	}
	var bankInfo models.BankInfo
	if err := models.DB.Model(&bankInfo).Where("bank_card_no=?", bankCardNo).Updates(values).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "更新成功",
	})
	return
}

func (c *BankInfoController) DeleteBankInfo() {
	bankCardNo := c.Ctx.Input.Param(":bank_card_no")
	var bankInfo models.BankInfo
	err := models.DB.Where("bank_card_no=?", bankCardNo).First(&bankInfo).Error
	if err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	err = upload.Delete(bankInfo.BankInfoFilePath)
	if err != nil{
		c.ToErrorResponse(errcode.ERROR_DELETE_FILE_FAIL.WithDetails(err.Error()))
		return
	}
	if err = models.DB.Delete(&bankInfo).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "删除成功",
	})
}


func (c *BankInfoController) UpdateFile() {
	userID := c.Ctx.Input.Param(":user_id")
	var bankInfo models.BankInfo
	err := models.DB.Where("user_id=?", userID).First(&bankInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if bankInfo.UserID == "" {
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
	fileInfo, err := upload.Update(upload.FileType(fileType), file, fileHeader, bankInfo.BankInfoFilePath)
	if err != nil {
		logs.Error("UploadFile err:", err)
		c.ToErrorResponse(errcode.ERROR_UPDATE_FILE_FAIL.WithDetails(err.Error()))
		return
	}
	values := map[string]interface{}{
		"bank_info_file_hash" : fileInfo.Hash,
		"bank_info_file_path" : fileInfo.Dst,
		"bank_info_file_upload_time" : time.Now(),
	}

	if err := models.DB.Model(&bankInfo).Where("user_id=?", userID).Updates(values).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "更新文件成功",
	})
}

func (c *BankInfoController) Download() {
	userID := c.Ctx.Input.Param(":user_id")
	var bankInfo models.BankInfo
	err := models.DB.Where("user_id=?", userID).First(&bankInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if bankInfo.UserID == "" {
		c.ToErrorResponse(errcode.ErrorCreateUserFail.WithDetails(errors.New("用户不存在").Error()))
		return
	}
	c.Ctx.Output.Download(bankInfo.BankInfoFilePath)
	return
}