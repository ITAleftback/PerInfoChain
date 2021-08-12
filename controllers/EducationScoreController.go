/**
 * @Author: Anpw
 * @Description:
 * @File:  EducationScoreController
 * @Version: 1.0.0
 * @Date: 2021/7/23 19:41
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

type EducationScoreController struct {
	BaseController
}

type AddScoreInfoRequest struct {
	EducationID string  `alias:"教育信息id" valid:"Required;"`
	TestTime    uint    `alias:"考试年份" valid:"Required;"`
	ClassName   string  `alias:"课程名称" valid:"Required;"`
	Score       float32 `alias:"分数" valid:"Required;"`
}

type UpdateScoreInfoRequest struct {
	TestTime   uint    `alias:"考试年份"`
	ClassName  string  `alias:"课程名称"`
	Score      float32 `alias:"分数"`
}

// AddEducationScoreInfo 添加教育用户基本信息
func (c *EducationScoreController) AddEducationScoreInfo() {
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

	scoreInfo := AddScoreInfoRequest{}
	if err := c.ParseForm(&scoreInfo); err != nil {
		logs.Error("ParseForm err: %v", err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}

	valid, errs := app.BindAndValid(&scoreInfo)
	if !valid {
		logs.Error("app.BindAndValid err: %v", errs)
		c.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//判断教育id是否存在是否已经存在
	var educationInfos models.EducationInfo
	if err := models.DB.Where("education_id=?", scoreInfo.EducationID).First(&educationInfos).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if educationInfos.EducationID == "" {
		c.ToErrorResponse(errcode.ErrorAddEducationScoreInfoFail.WithDetails(errors.New("教育信息不存在").Error()))
		return
	}
	fileInfo, err := upload.Upload(upload.FileType(fileType), file, fileHeader, "scoreinfo")
	if err != nil {
		logs.Error("UploadFile err:", err)
		c.ToErrorResponse(errcode.ERROR_UPLOAD_FILE_FAIL.WithDetails(err.Error()))
		return
	}
	u := models.EducationScoreInfo{
		ScoreID:                 utils.GenerateUID(),
		EducationID:             scoreInfo.EducationID,
		LoginUserID:             educationInfos.LoginUserID,
		TestTime:                scoreInfo.TestTime,
		ClassName:               scoreInfo.ClassName,
		Score:                   scoreInfo.Score,
		UploadTime:              time.Now(),
		ScoreInfoFilePath:       fileInfo.Dst,
		ScoreInfoFileHash:       fileInfo.Hash,
		ScoreInfoFileUploadTime: time.Now(),
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

// EducationScoreInfoList 获取用户列表信息 以及 查询
func (c *EducationScoreController) EducationScoreInfoList() {

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

	var scoreInfos []*models.EducationScoreInfo
	if err := models.DB.Find(&scoreInfos).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponseList(scoreInfos, &pager)
	return
}

func (c *EducationScoreController) UpdateEducationScoreInfo() {
	scoreID := c.Ctx.Input.Param(":score_id")
	scoreInfo := UpdateScoreInfoRequest{}
	if err := c.ParseForm(&scoreInfo); err != nil {
		logs.Error("ParseForm err: %v", err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	valid, errs := app.BindAndValid(&scoreInfo)
	if !valid {
		logs.Error("app.BindAndValid err: %v", errs)
		c.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	values := map[string]interface{}{}
	if scoreInfo.TestTime != 0 {
		values["test_time"] = scoreInfo.TestTime
	}
	if scoreInfo.ClassName != "" {
		values["class_name"] = scoreInfo.ClassName
	}
	if scoreInfo.Score != 0 {
		values["score"] = scoreInfo.Score
	}
	values["upload_time"] = time.Now()

	var educationScoreInfo models.EducationScoreInfo
	if err := models.DB.Model(&educationScoreInfo).Where("score_id=?", scoreID).Updates(values).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "更新成功",
	})
	return
}

func (c *EducationScoreController) DeleteEducationScoreInfo() {
	scoreID := c.Ctx.Input.Param(":score_id")
	var educationScoreInfo models.EducationScoreInfo
	err := models.DB.Where("score_id=?", scoreID).First(&educationScoreInfo).Error
	if err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	err = upload.Delete(educationScoreInfo.ScoreInfoFilePath)
	if err != nil{
		c.ToErrorResponse(errcode.ERROR_DELETE_FILE_FAIL.WithDetails(err.Error()))
		return
	}
	if err = models.DB.Delete(&educationScoreInfo).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "删除成功",
	})
	return
}

func (c *EducationScoreController) UpdateFile() {
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
	var scoreInfo models.EducationScoreInfo
	err = models.DB.Where("education_id=?", educationInfo.EducationID).First(&scoreInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if scoreInfo.EducationID == "" {
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
	fileInfo, err := upload.Update(upload.FileType(fileType), file, fileHeader, scoreInfo.ScoreInfoFilePath)
	if err != nil {
		logs.Error("UploadFile err:", err)
		c.ToErrorResponse(errcode.ERROR_UPDATE_FILE_FAIL.WithDetails(err.Error()))
		return
	}
	values := map[string]interface{}{
		"score_info_file_hash" : fileInfo.Hash,
		"score_info_file_path" : fileInfo.Dst,
		"score_info_file_upload_time" : time.Now(),
	}

	if err = models.DB.Model(&scoreInfo).Where("education_id=?", educationInfo.EducationID).Updates(values).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "更新文件成功",
	})
}

func (c *EducationScoreController) Download() {
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
	var scoreInfo models.EducationScoreInfo
	err = models.DB.Where("education_id=?", educationInfo.EducationID).First(&scoreInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if scoreInfo.EducationID == "" {
		c.ToErrorResponse(errcode.ErrorCreateUserFail.WithDetails(errors.New("用户不存在").Error()))
		return
	}

	c.Ctx.Output.Download(scoreInfo.ScoreInfoFilePath)
	return
}