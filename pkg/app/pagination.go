/**
 * @Author: Anpw
 * @Description:
 * @File:  pagination
 * @Version: 1.0.0
 * @Date: 2021/7/20 4:39
 */

package app

import (
	"PerInfoChain/pkg/convert"
	"github.com/astaxie/beego"
)

func GetPage(page int) int {
	if page <= 0 {
		return 1
	}
	return page
}

func GetPageSize(pageSize int) int {
	if pageSize <= 0 {
		defaultPageSize := beego.AppConfig.String("DefaultPageSize")
		pageSize = convert.StrTo(defaultPageSize).MustInt()
		return pageSize
	}
	maxPageSize := beego.AppConfig.String("MaxPageSize")
	maxSize := convert.StrTo(maxPageSize).MustInt()
	if pageSize > maxSize {
		return maxSize
	}
	return pageSize
}

// GetPageOffset 跳过result条数据 查询后pageSize条数据
func GetPageOffset(page, pageSize int) int {
	result := 0
	if page > 0 {
		result = (page - 1) * pageSize
	}
	return result
}


