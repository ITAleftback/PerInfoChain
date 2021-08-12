/**
 * @Author: Anpw
 * @Description:
 * @File:  CriminalCaseController
 * @Version: 1.0.0
 * @Date: 2021/7/23 20:37
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

type CriminalCaseController struct {
	BaseController
}

type AddCriminalCaseInfoRequest struct {
	UserID string `alias:"个人id" valid:"Required;"`
	Time   string `alias:"发生时间" valid:"Required;"`
	Place  string `alias:"发生地点" valid:"Required;"`
	Case   string `alias:"案件经过" valid:"Required;"`
	Degree uint `alias:"案件严重程度" valid:"Required;"`
	Punish string `alias:"所受处罚" valid:"Required;"`
}

type UpdateCriminalCaseInfoRequest struct {
	Time   string
	Place  string
	Case   string
	Degree uint
	Punish string
}

// AddCriminalCaseInfo 添加刑事案件基本信息
func (c *CriminalCaseController) AddCriminalCaseInfo() {
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

	criminalCaseInfo := AddCriminalCaseInfoRequest{}
	if err := c.ParseForm(&criminalCaseInfo); err != nil {
		logs.Error("ParseForm err: %v", err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}

	valid, errs := app.BindAndValid(&criminalCaseInfo)
	if !valid {
		logs.Error("app.BindAndValid err: %v", errs)
		c.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//判断用户是否存在
	var fabricUserInfo models.FabricUserInfo
	if err := models.DB.Where("user_id=?", criminalCaseInfo.UserID).First(&fabricUserInfo).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if fabricUserInfo.UserID == "" {
		c.ToErrorResponse(errcode.ErrorAddCriminalCaseFail.WithDetails(errors.New("用户不存在").Error()))
		return
	}
	fileInfo, err := upload.Upload(upload.FileType(fileType), file, fileHeader, "caseinfo")
	if err != nil {
		logs.Error("UploadFile err:", err)
		c.ToErrorResponse(errcode.ERROR_UPLOAD_FILE_FAIL.WithDetails(err.Error()))
		return
	}

	u := models.CriminalCase{
		CaseID:                 utils.GenerateUID(),
		LoginUserID:            fabricUserInfo.LoginUserID,
		UserID:                 criminalCaseInfo.UserID,
		Time:                   criminalCaseInfo.Time,
		Place:                  criminalCaseInfo.Place,
		Case:                   criminalCaseInfo.Case,
		Degree:                 criminalCaseInfo.Degree,
		Punish:                 criminalCaseInfo.Punish,
		CaseInfoFilePath:       fileInfo.Dst,
		CaseInfoFileHash:       fileInfo.Hash,
		CaseInfoFileUploadTime: time.Now(),
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

// CriminalCaseList 获取列表信息 以及 查询
func (c *CriminalCaseController) CriminalCaseList() {
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

	var criminalInfos []*models.CriminalCase
	if err := models.DB.Find(&criminalInfos).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponseList(criminalInfos, &pager)
	return
}

// GetCriminalInfo 查询
func (c *CriminalCaseController) GetCriminalInfo() {
	var userInfo models.FabricUserInfo
	var criminalCaseInfo []*models.CriminalCase
	IDCard := c.Ctx.Input.Param(":id_card")
	err := models.DB.Where("id_card=?", IDCard).First(&userInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	err = models.DB.Where("user_id=?", userInfo.UserID).First(&criminalCaseInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(criminalCaseInfo)
	return
}

func (c *CriminalCaseController) UpdateCriminalCaseInfo() {
	caseID := c.Ctx.Input.Param(":case_id")
	criminalCaseInfo := UpdateCriminalCaseInfoRequest{}
	if err := c.ParseForm(&criminalCaseInfo); err != nil {
		logs.Error("ParseForm err: %v", err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	valid, errs := app.BindAndValid(&criminalCaseInfo)
	if !valid {
		logs.Error("app.BindAndValid err: %v", errs)
		c.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	values := map[string]interface{}{}
	if criminalCaseInfo.Case != "" {
		values["case"] = criminalCaseInfo.Case
	}
	if criminalCaseInfo.Punish != "" {
		values["punish"] = criminalCaseInfo.Punish
	}
	if criminalCaseInfo.Time != "" {
		values["time"] = criminalCaseInfo.Time
	}
	if criminalCaseInfo.Degree != 0 {
		values["degree"] = criminalCaseInfo.Degree
	}
	if criminalCaseInfo.Place != "" {
		values["place"] = criminalCaseInfo.Place
	}

	var criminalCase models.CriminalCase
	if err := models.DB.Model(&criminalCase).Where("case_id=?", caseID).Updates(values).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "更新成功",
	})
	return
}

func (c *CriminalCaseController) DeleteCriminalCaseInfo() {
	caseID := c.Ctx.Input.Param(":case_id")
	var criminalCase models.CriminalCase
	err := models.DB.Where("case_id=?", caseID).First(&criminalCase).Error
	if err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	err = upload.Delete(criminalCase.CaseInfoFilePath)
	if err != nil{
		c.ToErrorResponse(errcode.ERROR_DELETE_FILE_FAIL.WithDetails(err.Error()))
		return
	}
	if err = models.DB.Delete(&criminalCase).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "删除成功",
	})
}


func (c *CriminalCaseController) UpdateFile() {
	userID := c.Ctx.Input.Param(":user_id")
	var caseInfo models.CriminalCase
	err := models.DB.Where("user_id=?", userID).First(&caseInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if caseInfo.UserID == "" {
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
	fileInfo, err := upload.Update(upload.FileType(fileType), file, fileHeader, caseInfo.CaseInfoFilePath)
	if err != nil {
		logs.Error("UploadFile err:", err)
		c.ToErrorResponse(errcode.ERROR_UPDATE_FILE_FAIL.WithDetails(err.Error()))
		return
	}
	values := map[string]interface{}{
		"case_info_file_hash" : fileInfo.Hash,
		"case_info_file_path" : fileInfo.Dst,
		"case_info_file_upload_time" : time.Now(),
	}

	if err := models.DB.Model(&caseInfo).Where("user_id=?", userID).Updates(values).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "更新文件成功",
	})
}

func (c *CriminalCaseController) Download() {
	userID := c.Ctx.Input.Param(":user_id")
	var criminalCase models.CriminalCase
	err := models.DB.Where("user_id=?", userID).First(&criminalCase).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if criminalCase.UserID == "" {
		c.ToErrorResponse(errcode.ErrorCreateUserFail.WithDetails(errors.New("用户不存在").Error()))
		return
	}
	c.Ctx.Output.Download(criminalCase.CaseInfoFilePath)
	return
}