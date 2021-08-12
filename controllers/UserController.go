/**
 * @Author: Anpw
 * @Description:
 * @File:  AuthController
 * @Version: 1.0.0
 * @Date: 2021/7/13 14:18
 */

package controllers

import (
	"PerInfoChain/models"
	"PerInfoChain/pkg/app"
	"PerInfoChain/pkg/errcode"
	"PerInfoChain/pkg/utils"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	"time"
)

type UserController struct {
	BaseController
}

// RegisterRequest tags中 alias 表示验证不通过时候的提示名称，valid表示验证的格式 分号分割不同的验证内容
type RegisterRequest struct {
	Username string `alias:"用户名" valid:"MinSize(3)"`
	Password string `alias:"密码" valid:"MinSize(6)"`
	Mobile   string `alias:"电话号码" valid:"Mobile"`
	Role     uint   `alias:"角色" valid:"Required;"`
	Status   uint   `form:"default=0"`
}

type LoginRequest struct {
	Mobile   string `alias:"电话号码" valid:"Required;"`
	Password string `alias:"密码" valid:"Required;"`
}

func (c *UserController) Register() {
	user := RegisterRequest{}

	if err := c.ParseForm(&user); err != nil {
		logs.Error("ParseForm err: %v", err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}

	valid, errs := app.BindAndValid(&user)
	if !valid {
		logs.Error("app.BindAndValid err: %v", errs)
		c.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	//判断用户是否已经存在
	var users []models.User
	if err := models.DB.Where("phone=?", user.Mobile).Find(&users).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if len(users) > 0 {
		c.ToErrorResponse(errcode.ErrorCreateUserFail.WithDetails(errors.New("手机号码已被注册").Error()))
		return
	}
	hashedPassword, err := utils.GeneratePassHash(user.Password)
	if err != nil {
		logs.Error("加密错误 err:", err)
		c.ToErrorResponse(errcode.ServerError)
	}
	u := models.User{
		LoginUserID: utils.GenerateUID(),
		Username:    user.Username,
		Password:    hashedPassword,
		Phone:       user.Mobile,
		LastLogin:   time.Now(),
		Role:        user.Role,
	}
	if err = models.DB.Create(&u).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "注册成功",
	})
	return
}

func (c *UserController) Login() {

	userLogin := LoginRequest{}

	if err := c.ParseForm(&userLogin); err != nil {
		logs.Error("ParseForm err: %v", err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}

	valid, errs := app.BindAndValid(&userLogin)
	if !valid {
		logs.Error("app.BindAndValid err: %v", errs)
		c.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	var user models.User
	if err := models.DB.Where("phone=?", userLogin.Mobile).First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}

	//判断是否存在
	if user.LoginUserID == "" {
		c.ToErrorResponse(errcode.ErrorLoginUserFail.WithDetails(errors.New("手机或密码不正确").Error()))
		return
	}

	flag := utils.CompareHash(user.Password, userLogin.Password)
	if !flag {
		c.ToErrorResponse(errcode.ErrorLoginUserFail.WithDetails(errors.New("手机或密码不正确").Error()))
		return
	} else {
		if user.Status == 0 {
			c.ToErrorResponse(errcode.ErrorLoginUserFail.WithDetails(errors.New("你的身份未确认，请等待管理员核实").Error()))
		}

		token, err := app.GenerateToken(user)
		if err != nil {
			logs.Error("发放token错误:", err)
			c.ToErrorResponse(errcode.ServerError)
			return
		}
		if err = models.DB.Model(&user).Update("last_login", time.Now()).Error; err != nil && err != gorm.ErrRecordNotFound {
			logs.Error(err)
			c.ToErrorResponse(errcode.ServerError)
			return
		}
		c.ToResponse(map[string]interface{}{
			"msg":   "登录成功",
			"token": token,
			"role":  user.Role,
		})
		return
	}
}

// UserList 获取所有用户
func (c *UserController) UserList() {
	page, _ := c.GetInt("page")
	pageSize, _ := c.GetInt("page_size")
	var count int
	if err := models.DB.Table(models.User{}.TableName()).Count(&count).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}

	pager := Pager{Page: app.GetPage(page), PageSize: app.GetPageSize(pageSize), TotalRows: count}
	pageOffset := app.GetPageOffset(pager.Page, pager.PageSize)
	if pageOffset >= 0 && pageSize > 0 {
		models.DB = models.DB.Offset(pageOffset).Limit(pageSize)
	}

	var users []*models.User
	if err := models.DB.Find(&users).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponseList(users, &pager)
	return
}

// UserStatusList 获取所有status为0的用户
func (c *UserController) UserStatusList() {
	page, _ := c.GetInt("page")
	pageSize, _ := c.GetInt("page_size")
	var count int
	if err := models.DB.Table(models.User{}.TableName()).Where("status=?", 0).Count(&count).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}

	pager := Pager{Page: app.GetPage(page), PageSize: app.GetPageSize(pageSize), TotalRows: count}
	pageOffset := app.GetPageOffset(pager.Page, pager.PageSize)
	if pageOffset >= 0 && pageSize > 0 {
		models.DB = models.DB.Offset(pageOffset).Limit(pageSize)
	}

	var users []*models.User
	if err := models.DB.Where("status=?", 0).Find(&users).Error; err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponseList(users, &pager)
	return
}

func (c *UserController) DeleteUser() {
	loginUserID := c.Ctx.Input.Param(":login_user_id")
	var user models.User
	if err := models.DB.Where("login_user_id=?", loginUserID).Delete(&user).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "删除成功",
	})
}

func (c *UserController) UpdateUserStatus() {
	loginUserID := c.Ctx.Input.Param(":login_user_id")
	var user models.User
	err := models.DB.Where("login_user_id=?", loginUserID).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	if user.LoginUserID == "" {
		c.ToErrorResponse(errcode.ErrorCreateUserFail.WithDetails(errors.New("用户不存在").Error()))
		return
	}

	values := map[string]interface{}{
		"status": 1,
	}

	if err := models.DB.Model(&user).Where("login_user_id=?", loginUserID).Updates(values).Error; err != nil {
		logs.Error(err)
		c.ToErrorResponse(errcode.ServerError)
		return
	}
	c.ToResponse(map[string]interface{}{
		"msg": "审核通过",
	})
	return
}
