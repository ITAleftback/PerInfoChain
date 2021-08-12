/**
 * @Author: Anpw
 * @Description:
 * @File:  apiRouter
 * @Version: 1.0.0
 * @Date: 2021/7/25 22:14
 */

package routers

import (
	"PerInfoChain/controllers"
	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/api/v1",
		beego.NSRouter("/", &controllers.MenuController{}),
		beego.NSRouter("/menu", &controllers.MenuController{}, "get:Menu"),
	)
	beego.AddNamespace(ns)
}
