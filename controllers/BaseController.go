/**
 * @Author: Anpw
 * @Description:
 * @File:  BaseController
 * @Version: 1.0.0
 * @Date: 2021/7/14 2:11
 */

package controllers

import (
	"PerInfoChain/pkg/errcode"
	"github.com/astaxie/beego"
	"net/http"
)

type BaseController struct {
	beego.Controller
}

// Pager page 指第几页 pageSize指每页请求多少条数据 totalRows指总计数据
type Pager struct {
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	TotalRows int `json:"total_rows"`
}

// ToResponse 返回json数据
func (c *BaseController) ToResponse(data interface{}) {
	response := map[string]interface{}{"code": http.StatusOK, "data": data}
	c.Data["json"] = response
	c.ServeJSON()
}

func (c *BaseController) ToResponseList(list interface{}, pager *Pager) {
	data := map[string]interface{}{
		"list": list,
		"pager": Pager{
			Page:      pager.Page,
			PageSize:  pager.PageSize,
			TotalRows: pager.TotalRows,
		},
	}
	c.Data["json"] = map[string]interface{}{"code": http.StatusOK, "data": data}
	c.ServeJSON()
}

func (c *BaseController) ToErrorResponse(err *errcode.Error) {
	response := map[string]interface{}{"code": err.Code, "msg": err.Msg}
	details := err.Details
	if len(details) > 0 {
		response["details"] = details
	}
	c.Data["json"] = response
	c.ServeJSON()
}
