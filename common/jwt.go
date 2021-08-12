/**
 * @Author: Anpw
 * @Description:
 * @File:  jwt
 * @Version: 1.0.0
 * @Date: 2021/7/18 19:27
 */

package common

import (
	"PerInfoChain/models"
	"PerInfoChain/pkg/app"
	"PerInfoChain/pkg/errcode"
	"encoding/json"
	"github.com/astaxie/beego/context"
	"github.com/dgrijalva/jwt-go"
)

var ChenkToken = func(ctx *context.Context) {
	var (
		token string
		ecode = errcode.Success
	)
	token = ctx.Input.Header("token")
	if token == "" {
		ecode = errcode.InvalidParams
	} else {
		claims, err := app.ParseToken(token)
		if err != nil {
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				ecode = errcode.UnauthorizedTokenTimeout
			default:
				ecode = errcode.UnauthorizedTokenError
			}
		} else {
			var user models.User
			models.DB.Where("login_user_id=?", claims.LoginUserID).First(&user)
			if user.LoginUserID == "" {
				ecode = errcode.UnauthorizedAuthNotExist
			}
		}

	}

	if ecode != errcode.Success {
		JsonData := ReturnData{
			Code: ecode.Code,
			Msg:  ecode.Msg,
		}
		//为了架构  自行封装的json
		ret, _ := json.Marshal(JsonData)
		ctx.Output.Header("Content-Type", "application/json")
		ctx.ResponseWriter.Write(ret)
		return
	}
}

type ReturnData struct {
	Code int
	Msg  interface{}
}
