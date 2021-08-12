/**
 * @Author: Anpw
 * @Description:
 * @File:  BankDetailsController
 * @Version: 1.0.0
 * @Date: 2021/7/25 18:53
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
	"math"
	"time"
)

type BankDetailsController struct {
	BaseController
}

type AddBankDetailsRequest struct {
	BankCardNo    string  `alias:"银行卡号" valid:"Required;"`
	BeforeDeposit float64 `alias:"变更前存款" valid:"Required;"`
	AfterDeposit  float64 `alias:"变更后存款" valid:"Required;"`
	OccurTime     string  `alias:"发生时间" valid:"Required;"`
}

type UpdateBankDetailsRequest struct {
	BankCardNo    string
	BeforeDeposit float64
	AfterDeposit  float64
	OccurTime     string
}

// AddBankDetails 添加教育用户基本信息
func (c *BankDetailsController) AddBankDetails() {
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

	bankDetails := AddBankDetailsRequest{}
	if err := c.ParseForm(&bankDetails); err != nil {
		logs.Error("ParseForm err: %v", err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}

	valid, errs := app.BindAndValid(&bankDetails)
	if !valid {
		logs.Error("app.BindAndValid err: %v", errs)
		c.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	fileInfo, err := upload.Upload(upload.FileType(fileType), file, fileHeader, "bankdetails")
	if err != nil {
		logs.Error("UploadFile err:", err)
		c.ToErrorResponse(errcode.ERROR_UPLOAD_FILE_FAIL.WithDetails(err.Error()))
		return
	}
	//判断银行卡号是否已经存在
	var bankInfo models.BankInfo
	if err := models.DB.Where("bank_card_no=?", bankDetails.BankCardNo).First(&bankInfo).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if bankInfo.BankCardNo == "" {
		c.ToErrorResponse(errcode.ErrorAddBankDetailsFail.WithDetails(errors.New("银行卡号不存在").Error()))
		return
	}

	u := models.BankDetails{
		BankMonthID:               utils.GenerateUID(),
		LoginUserID:               bankInfo.LoginUserID,
		BankCardNo:                bankDetails.BankCardNo,
		BeforeDeposit:             bankDetails.BeforeDeposit,
		AfterDeposit:              bankDetails.AfterDeposit,
		Difference:                math.Abs(bankDetails.AfterDeposit - bankDetails.BeforeDeposit),
		OccurTime:                 bankDetails.OccurTime,
		BankDetailsFilePath:       fileInfo.Dst,
		BankDetailsFileHash:       fileInfo.Hash,
		BankDetailsFileUploadTime: time.Now(),
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

// BankDetailsList 获取列表信息 以及 查询
func (c *BankDetailsController) BankDetailsList() {

	page, _ := c.GetInt("page")
	pageSize, _ := c.GetInt("page_size")

	var count int
	if err := models.DB.Table(models.EducationScoreInfo{}.TableName()).Count(&count).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}

	pager := Pager{Page: app.GetPage(page), PageSize: app.GetPageSize(pageSize), TotalRows: count}
	pageOffset := app.GetPageOffset(pager.Page, pager.PageSize)
	if pageOffset >= 0 && pageSize > 0 {
		models.DB = models.DB.Offset(pageOffset).Limit(pageSize)
	}

	var bankDetails []*models.BankDetails
	if err := models.DB.Find(&bankDetails).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponseList(bankDetails, &pager)
	return
}

func (c *BankDetailsController) UpdateBankDetails() {
	bankMonthID := c.Ctx.Input.Param(":bank_month_id")
	bankDetailsRequest := UpdateBankDetailsRequest{}
	if err := c.ParseForm(&bankDetailsRequest); err != nil {
		logs.Error("ParseForm err: %v", err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	valid, errs := app.BindAndValid(&bankDetailsRequest)
	if !valid {
		logs.Error("app.BindAndValid err: %v", errs)
		c.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	values := map[string]interface{}{}
	if bankDetailsRequest.BankCardNo != "" {
		values["bank_card_no"] = bankDetailsRequest.BankCardNo
	}
	if bankDetailsRequest.BeforeDeposit != 0 {
		values["before_deposit"] = bankDetailsRequest.BeforeDeposit
	}
	if bankDetailsRequest.AfterDeposit != 0 {
		values["after_deposit"] = bankDetailsRequest.AfterDeposit
	}
	values["difference"] = math.Abs(bankDetailsRequest.AfterDeposit - bankDetailsRequest.BeforeDeposit)

	var bankDetails models.BankDetails
	if err := models.DB.Model(&bankDetails).Where("bank_month_id=?", bankMonthID).Updates(values).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "更新成功",
	})
	return
}

func (c *BankDetailsController) DeleteBankDetails() {
	bankMonthID := c.Ctx.Input.Param(":bank_month_id")
	var bankDetails models.BankDetails
	err := models.DB.Where("bank_month_id=?", bankMonthID).First(&bankDetails).Error
	if err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	err = upload.Delete(bankDetails.BankDetailsFilePath)
	if err != nil {
		c.ToErrorResponse(errcode.ERROR_DELETE_FILE_FAIL.WithDetails(err.Error()))
		return
	}
	if err = models.DB.Delete(&bankDetails).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "删除成功",
	})
}

func (c *BankDetailsController) UpdateFile() {
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

	var bankDetails models.BankDetails
	err = models.DB.Where("bank_card_no=?", bankInfo.BankCardNo).First(&bankDetails).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if bankDetails.BankMonthID == "" {
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
	fileInfo, err := upload.Update(upload.FileType(fileType), file, fileHeader, bankDetails.BankDetailsFilePath)
	if err != nil {
		logs.Error("UploadFile err:", err)
		c.ToErrorResponse(errcode.ERROR_UPDATE_FILE_FAIL.WithDetails(err.Error()))
		return
	}
	values := map[string]interface{}{
		"bank_details_file_hash":        fileInfo.Hash,
		"bank_details_file_path":        fileInfo.Dst,
		"bank_details_file_upload_time": time.Now(),
	}

	if err := models.DB.Model(&bankDetails).Where("bank_card_no=?", bankInfo.BankCardNo).Updates(values).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "更新文件成功",
	})
}

func (c *BankDetailsController) Download() {
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
	var bankDetails models.BankDetails
	err = models.DB.Where("bank_card_no=?", bankInfo.BankCardNo).First(&bankDetails).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if bankDetails.BankMonthID == "" {
		c.ToErrorResponse(errcode.ErrorCreateUserFail.WithDetails(errors.New("用户不存在").Error()))
		return
	}

	c.Ctx.Output.Download(bankDetails.BankDetailsFilePath)
	return
}
