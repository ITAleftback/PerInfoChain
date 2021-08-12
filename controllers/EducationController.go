/**
 * @Author: Anpw
 * @Description:
 * @File:  EducationController
 * @Version: 1.0.0
 * @Date: 2021/7/22 12:22
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

type EducationInfoController struct {
	BaseController
}

type AddEducationInfoRequest struct {
	UserID     string `alias:"个人id" valid:"Required;"`
	SchoolName string `alias:"学校名称" valid:"Required"`
	StartYear  uint   `alias:"起始年份" valid:"Required;"`
	EndYear    uint   `alias:"终止年份" valid:"Required;"`
	Major      string `alias:"专业" valid:"Required"`
}

type UpdateEducationInfoRequest struct {
	SchoolName string `alias:"学校名称"`
	StartYear  uint   `alias:"起始年份"`
	EndYear    uint   `alias:"终止年份"`
	Major      string `alias:"专业"`
}

// AddEducationUserInfo 添加教育用户基本信息
func (c *EducationInfoController) AddEducationUserInfo() {
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

	educationInfo := AddEducationInfoRequest{}
	if err := c.ParseForm(&educationInfo); err != nil {
		logs.Error("ParseForm err: %v", err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}

	valid, errs := app.BindAndValid(&educationInfo)
	if !valid {
		logs.Error("app.BindAndValid err: %v", errs)
		c.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	var fabricUserInfo models.FabricUserInfo
	if err := models.DB.Where("user_id=?", educationInfo.UserID).First(&fabricUserInfo).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if fabricUserInfo.UserID == "" {
		c.ToErrorResponse(errcode.ErrorAddEducationUserInfoFail.WithDetails(errors.New("用户不存在").Error()))
		return
	}
	fileInfo, err := upload.Upload(upload.FileType(fileType), file, fileHeader, "educationinfo")
	if err != nil {
		logs.Error("UploadFile err:", err)
		c.ToErrorResponse(errcode.ERROR_UPLOAD_FILE_FAIL.WithDetails(err.Error()))
		return
	}
	u := models.EducationInfo{
		EducationID:                 utils.GenerateUID(),
		UserID:                      educationInfo.UserID,
		LoginUserID:                 fabricUserInfo.LoginUserID,
		SchoolName:                  educationInfo.SchoolName,
		StartYear:                   educationInfo.StartYear,
		EndYear:                     educationInfo.EndYear,
		Major:                       educationInfo.Major,
		EducationInfoFilePath:       fileInfo.Dst,
		EducationInfoFileHash:       fileInfo.Hash,
		EducationInfoFileUploadTime: time.Now(),
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

// EducationUserInfoList 获取用户列表信息 以及 查询
func (c *EducationInfoController) EducationUserInfoList() {

	page, _ := c.GetInt("page")
	pageSize, _ := c.GetInt("page_size")

	var count int
	if err := models.DB.Table(models.EducationInfo{}.TableName()).Count(&count).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}

	pager := Pager{Page: app.GetPage(page), PageSize: app.GetPageSize(pageSize), TotalRows: count}
	pageOffset := app.GetPageOffset(pager.Page, pager.PageSize)
	if pageOffset >= 0 && pageSize > 0 {
		models.DB = models.DB.Offset(pageOffset).Limit(pageSize)
	}

	var educationInfos []*models.EducationInfo
	if err := models.DB.Find(&educationInfos).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponseList(educationInfos, &pager)
	return
}

// GetEducationUserInfo 查询
func (c *EducationInfoController) GetEducationUserInfo() {
	var userInfo models.FabricUserInfo
	var educationInfo []*models.EducationInfo
	IDCard := c.Ctx.Input.Param(":id_card")
	err := models.DB.Where("id_card=?", IDCard).First(&userInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	err = models.DB.Where("user_id=?", userInfo.UserID).First(&educationInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(educationInfo)
	return
}

func (c *EducationInfoController) UpdateEducationUserInfo() {
	educationID := c.Ctx.Input.Param(":education_id")
	educationInfo := UpdateEducationInfoRequest{}
	if err := c.ParseForm(&educationInfo); err != nil {
		logs.Error("ParseForm err: %v", err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	valid, errs := app.BindAndValid(&educationInfo)
	if !valid {
		logs.Error("app.BindAndValid err: %v", errs)
		c.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	values := map[string]interface{}{}
	if educationInfo.SchoolName != "" {
		values["school_name"] = educationInfo.SchoolName
	}
	if educationInfo.StartYear != 0 {
		values["start_year"] = educationInfo.StartYear
	}
	if educationInfo.EndYear != 0 {
		values["end_year"] = educationInfo.EndYear
	}
	if educationInfo.Major != "" {
		values["major"] = educationInfo.Major
	}

	var educationUserInfo models.EducationInfo
	if err := models.DB.Model(&educationUserInfo).Where("education_id=?", educationID).Updates(values).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "更新成功",
	})
	return
}

func (c *EducationInfoController) DeleteEducationUserInfo() {
	educationID := c.Ctx.Input.Param(":education_id")
	var educationUserInfo models.EducationInfo
	models.DB.Where("education_id=?", educationID)
	if err := models.DB.Where("education_id=?", educationID).First(&educationUserInfo).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	//将其对应文件删除
	err := upload.Delete(educationUserInfo.EducationInfoFilePath)
	if err != nil {
		c.ToErrorResponse(errcode.ERROR_DELETE_FILE_FAIL.WithDetails(err.Error()))
		return
	}
	if err = models.DB.Delete(&educationUserInfo).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "删除成功",
	})
}

func (c *EducationInfoController) UpdateFile() {
	userID := c.Ctx.Input.Param(":user_id")
	var educationInfo models.EducationInfo
	err := models.DB.Where("user_id=?", userID).First(&educationInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if educationInfo.UserID == "" {
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
	fileInfo, err := upload.Update(upload.FileType(fileType), file, fileHeader, educationInfo.EducationInfoFilePath)
	if err != nil {
		logs.Error("UploadFile err:", err)
		c.ToErrorResponse(errcode.ERROR_UPDATE_FILE_FAIL.WithDetails(err.Error()))
		return
	}
	values := map[string]interface{}{
		"education_info_file_hash":        fileInfo.Hash,
		"education_info_file_path":        fileInfo.Dst,
		"education_info_file_upload_time": time.Now(),
	}

	if err := models.DB.Model(&educationInfo).Where("user_id=?", userID).Updates(values).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "更新文件成功",
	})
}

func (c *EducationInfoController) Download() {
	userID := c.Ctx.Input.Param(":user_id")
	var educationInfo models.EducationInfo
	err := models.DB.Where("user_id=?", userID).First(&educationInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if educationInfo.UserID == "" {
		c.ToErrorResponse(errcode.ErrorCreateUserFail.WithDetails(errors.New("用户不存在").Error()))
		return
	}
	c.Ctx.Output.Download(educationInfo.EducationInfoFilePath)
	return
}
