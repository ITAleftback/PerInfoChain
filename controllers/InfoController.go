/**
 * @Author: Anpw
 * @Description:
 * @File:  InfoController
 * @Version: 1.0.0
 * @Date: 2021/7/19 1:19
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

type InfoController struct {
	BaseController
}

// AddUserInfoRequest FabricUserInfo 添加了一些唯一索引
type AddUserInfoRequest struct {
	LoginUserID   string `alias:"个人id" valid:"Required;"`
	IDCard        string `alias:"身份证" valid:"Required;Length(18);"`
	Name          string `alias:"姓名" valid:"Required;"`
	Sex           uint   `alias:"性别" valid:"Required;"`
	Nation        string `alias:"民族" valid:"Required;"`
	Native        string `alias:"籍贯" valid:"Required;"`
	Birthday      string `alias:"生日" valid:"Required;"`
	Phone         string `alias:"手机" valid:"Mobile;"`
	Email         string `alias:"邮箱" valid:"Email;"`
	PoliticalLook string `alias:"政治面貌" valid:"Required;"`
	HomeAddress   string `alias:"家庭住址" valid:"Required;"`
}

type UpdateUserInfoRequest struct {
	IDCard        string `alias:"身份证" valid:"Length(18);"`
	Name          string `alias:"姓名"`
	Sex           uint   `alias:"性别"`
	Nation        string `alias:"民族"`
	Native        string `alias:"籍贯"`
	Birthday      string `alias:"生日"`
	Phone         string `alias:"手机" valid:"Mobile;"`
	Email         string `alias:"邮箱" valid:"Email;"`
	PoliticalLook string `alias:"政治面貌"`
	HomeAddress   string `alias:"家庭住址"`
}

// AddUserInfo 添加用户基本信息
func (c *InfoController) AddUserInfo() {
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

	userInfo := AddUserInfoRequest{}
	if err := c.ParseForm(&userInfo); err != nil {
		logs.Error("ParseForm err: %v", err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}

	valid, errs := app.BindAndValid(&userInfo)
	if !valid {
		logs.Error("app.BindAndValid err: %v", errs)
		c.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	//判断用户是否已经存在
	var users []models.User
	if err := models.DB.Where("login_user_id=?", userInfo.LoginUserID).Find(&users).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if len(users) == 0 {
		c.ToErrorResponse(errcode.ErrorCreateUserFail.WithDetails(errors.New("用户不存在").Error()))
		return
	}
	//判断用户是否已经存在
	var userInfos []models.FabricUserInfo

	if err = models.DB.Where("id_card=? OR phone=?", userInfo.IDCard, userInfo.Phone).Find(&userInfos).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if len(userInfos) > 0 {
		c.ToErrorResponse(errcode.ErrorAddUserInfoFail.WithDetails(errors.New("基本信息已添加").Error()))
		return
	}
	/**
	 * TODO-Anpw: 2021/7/29 7:50 此处有个小bug 就是文件上传后添加数据库失败会造成不一，想过用事务但是文件上传不属于数据库范畴，想办法化为原子操作
	 * Description:
	 */
	fileInfo, err := upload.Upload(upload.FileType(fileType), file, fileHeader, "userinfo")
	if err != nil {
		logs.Error("UploadFile err:", err)
		c.ToErrorResponse(errcode.ERROR_UPLOAD_FILE_FAIL.WithDetails(err.Error()))
		return
	}
	u := models.FabricUserInfo{
		UserID:                 utils.GenerateUID(),
		LoginUserID:            userInfo.LoginUserID,
		IDCard:                 userInfo.IDCard,
		Name:                   userInfo.Name,
		Sex:                    userInfo.Sex,
		Nation:                 userInfo.Nation,
		Native:                 userInfo.Native,
		Birthday:               userInfo.Birthday,
		Phone:                  userInfo.Phone,
		Email:                  userInfo.Email,
		PoliticalLook:          userInfo.PoliticalLook,
		HomeAddress:            userInfo.HomeAddress,
		UserInfoFilePath:       fileInfo.Dst,
		UserInfoFileHash:       fileInfo.Hash,
		UserInfoFileUploadTime: time.Now(),
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

// UserInfoList 获取用户列表信息
func (c *InfoController) UserInfoList() {
	page, _ := c.GetInt("page")
	pageSize, _ := c.GetInt("page_size")

	var count int
	if err := models.DB.Table(models.FabricUserInfo{}.TableName()).Count(&count).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}

	pager := Pager{Page: app.GetPage(page), PageSize: app.GetPageSize(pageSize), TotalRows: count}
	pageOffset := app.GetPageOffset(pager.Page, pager.PageSize)
	if pageOffset >= 0 && pageSize > 0 {
		models.DB = models.DB.Offset(pageOffset).Limit(pageSize)
	}

	var userInfos []*models.FabricUserInfo
	if err := models.DB.Find(&userInfos).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponseList(userInfos, &pager)
	return
}

// GetUserInfo 查询
func (c *InfoController) GetUserInfo() {
	var fabricUserInfos []*models.FabricUserInfo
	IDCard := c.Ctx.Input.Param(":id_card")

	err := models.DB.Where("id_card = ?", IDCard).First(&fabricUserInfos).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(fabricUserInfos)
	return
}

func (c *InfoController) UpdateUserInfo() {

	userID := c.Ctx.Input.Param(":user_id")
	userInfo := UpdateUserInfoRequest{}
	if err := c.ParseForm(&userInfo); err != nil {
		logs.Error("ParseForm err: %v", err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	valid, errs := app.BindAndValid(&userInfo)
	if !valid {
		logs.Error("app.BindAndValid err: %v", errs)
		c.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	var user models.FabricUserInfo
	err := models.DB.Where("user_id=?", userID).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	//判断是不是自己
	if user.Phone != userInfo.Phone || user.IDCard != userInfo.IDCard {
		//判断用户是否已经存在
		var users []models.FabricUserInfo
		if err = models.DB.Where("id_card=? OR phone=?", userInfo.IDCard, userInfo.Phone).Find(&users).Error; err != nil && err != gorm.ErrRecordNotFound {
			logs.Error(err)
			c.ToErrorResponse(errcode.ServerError)
			return
		}
		if len(users) > 0 {
			c.ToErrorResponse(errcode.ErrorUpdateUserInfoFail.WithDetails(errors.New("电话号码或身份证已被使用").Error()))
			return
		}
	}

	values := map[string]interface{}{}
	if userInfo.IDCard != "" {
		values["id_card"] = userInfo.IDCard
	}
	if userInfo.Sex != 0 {
		values["sex"] = userInfo.Sex
	}
	if userInfo.Birthday != "" {
		values["birthday"] = userInfo.Birthday
	}
	if userInfo.Name != "" {
		values["name"] = userInfo.Name
	}
	if userInfo.Nation != "" {
		values["nation"] = userInfo.Nation
	}
	if userInfo.Native != "" {
		values["native"] = userInfo.Native
	}
	if userInfo.PoliticalLook != "" {
		values["political_look"] = userInfo.PoliticalLook
	}
	if userInfo.HomeAddress != "" {
		values["home_address"] = userInfo.HomeAddress
	}
	if userInfo.Phone != "" {
		values["phone"] = userInfo.Phone
	}
	if userInfo.Email != "" {
		values["email"] = userInfo.Email
	}


	var fabricUserInfo models.FabricUserInfo
	if err := models.DB.Model(&fabricUserInfo).Where("user_id=?", userID).Updates(values).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "更新成功",
	})
	return
}

func (c *InfoController) DeleteUserInfo() {
	userID := c.Ctx.Input.Param(":user_id")
	var fabricUserInfo models.FabricUserInfo
	err := models.DB.Where("user_id=?", userID).First(&fabricUserInfo).Error
	if err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	err = upload.Delete(fabricUserInfo.UserInfoFilePath)
	if err != nil{
		c.ToErrorResponse(errcode.ERROR_DELETE_FILE_FAIL.WithDetails(err.Error()))
		return
	}
	if err = models.DB.Delete(&fabricUserInfo).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "删除成功",
	})
	return
}


func (c *InfoController) UpdateFile() {
	userID := c.Ctx.Input.Param(":user_id")
	var fabricUserInfo models.FabricUserInfo
	err := models.DB.Where("user_id=?", userID).First(&fabricUserInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if fabricUserInfo.UserID == "" {
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
	fileInfo, err := upload.Update(upload.FileType(fileType), file, fileHeader, fabricUserInfo.UserInfoFilePath)
	if err != nil {
		logs.Error("UploadFile err:", err)
		c.ToErrorResponse(errcode.ERROR_UPDATE_FILE_FAIL.WithDetails(err.Error()))
		return
	}
	values := map[string]interface{}{
		"user_info_file_hash" : fileInfo.Hash,
		"user_info_file_path" : fileInfo.Dst,
		"user_info_file_upload_time" : time.Now(),
	}

	if err := models.DB.Model(&fabricUserInfo).Where("user_id=?", userID).Updates(values).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "更新文件成功",
	})
}

func (c *InfoController) Download() {
	userID := c.Ctx.Input.Param(":user_id")
	var fabricUserInfo models.FabricUserInfo
	err := models.DB.Where("user_id=?", userID).First(&fabricUserInfo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if fabricUserInfo.UserID == "" {
		c.ToErrorResponse(errcode.ErrorCreateUserFail.WithDetails(errors.New("用户不存在").Error()))
		return
	}
	c.Ctx.Output.Download(fabricUserInfo.UserInfoFilePath)
	return
}