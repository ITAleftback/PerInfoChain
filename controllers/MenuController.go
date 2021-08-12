/**
 * @Author: Anpw
 * @Description:
 * @File:  MenuController
 * @Version: 1.0.0
 * @Date: 2021/7/25 22:19
 */

package controllers

import (
	"PerInfoChain/models"
	"github.com/astaxie/beego"
)

type MenuController struct {
	beego.Controller
}

func (c *MenuController) Menu() {
	var menu []models.Menu
	models.DB.Find(&menu)
	c.Data["json"] = menu
	c.ServeJSON()
}
