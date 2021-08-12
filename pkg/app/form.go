/**
 * @Author: Anpw
 * @Description: 验证器封装
 * @File:  form
 * @Version: 1.0.0
 * @Date: 2021/7/14 2:29
 */

package app

import (
	"github.com/astaxie/beego/validation"
	"reflect"
	"strings"
)

type ValidError struct {
	Name string
	Key     string
	Message string
}

type ValidErrors []*ValidError

func (v *ValidError) Error() string {
	return v.Message
}

func (v ValidErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}

func (v ValidErrors) Errors() []string {
	var errs []string
	for _, err := range v {
		errs = append(errs, err.Error())
	}
	return errs
}

// BindAndValid 验证器设置
func BindAndValid(v interface{}) (bool, ValidErrors) {
	valid := validation.Validation{}
	var errs ValidErrors
	// 自定义消息提示
	var Message = map[string]string{
		"Required": "不能为空",
		"MinSize":  "最短长度为 %d",
		"Length":   "长度必须为 %d",
		"Numeric":  "必须是有效的数字",
		"Email":    "无效的电子邮件地址",
		"Mobile":   "无效的手机号码",
	}
	validation.SetDefaultMessage(Message)

	// 验证是否struct tag是否正确
	validResult, _ := valid.Valid(v)

	if !validResult {
		// 验证没通过
		st := reflect.TypeOf(v).Elem()
		for _, err := range valid.Errors {
			//获取验证的字段名和提示信息的别名
			filed, _ := st.FieldByName(err.Field)
			var alias = filed.Tag.Get("alias")
			errs = append(errs, &ValidError{
				Key:     err.Key,
				Message: alias+err.Message,
			})
		}
		return false, errs
	}
	return true, nil
}
