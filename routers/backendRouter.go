package routers

import (
	"PerInfoChain/common"
	"PerInfoChain/controllers"
	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/api/v1",
		beego.NSRouter("/register", &controllers.UserController{}, "post:Register"),
		beego.NSRouter("/login", &controllers.UserController{}, "post:Login"),

		beego.NSNamespace("/user",
			beego.NSBefore(common.ChenkToken),
			//基本信息模块
			beego.NSRouter("/adduserinfo", &controllers.InfoController{}, "post:AddUserInfo"),
			beego.NSRouter("/list", &controllers.InfoController{}, "get:UserInfoList"),
			beego.NSRouter("/:id_card", &controllers.InfoController{}, "get:GetUserInfo"),
			beego.NSRouter("/:user_id", &controllers.InfoController{}, "put:UpdateUserInfo"),
			beego.NSRouter("/:user_id", &controllers.InfoController{}, "delete:DeleteUserInfo"),
			beego.NSRouter("/file/:user_id", &controllers.InfoController{}, "put:UpdateFile"),
			beego.NSRouter("/file/:user_id", &controllers.InfoController{}, "get:Download"),
		),
		beego.NSNamespace("/education",
			beego.NSBefore(common.ChenkToken),
			//教育模块
			beego.NSRouter("/addeducationinfo", &controllers.EducationInfoController{}, "post:AddEducationUserInfo"),
			beego.NSRouter("/list", &controllers.EducationInfoController{}, "get:EducationUserInfoList"),
			beego.NSRouter("/:id_card", &controllers.EducationInfoController{}, "get:GetEducationUserInfo"),
			beego.NSRouter("/:education_id", &controllers.EducationInfoController{}, "put:UpdateEducationUserInfo"),
			beego.NSRouter("/:education_id", &controllers.EducationInfoController{}, "delete:DeleteEducationUserInfo"),
			beego.NSRouter("/file/:user_id", &controllers.EducationInfoController{}, "put:UpdateFile"),
			beego.NSRouter("/file/:user_id", &controllers.EducationInfoController{}, "get:Download"),

			//教育成绩模块
			beego.NSRouter("/addscoreinfo", &controllers.EducationScoreController{}, "post:AddEducationScoreInfo"),
			beego.NSRouter("/score/list", &controllers.EducationScoreController{}, "get:EducationScoreInfoList"),
			beego.NSRouter("/score/:score_id", &controllers.EducationScoreController{}, "put:UpdateEducationScoreInfo"),
			beego.NSRouter("/score/:score_id", &controllers.EducationScoreController{}, "delete:DeleteEducationScoreInfo"),
			beego.NSRouter("/score/file/:user_id", &controllers.EducationScoreController{}, "put:UpdateFile"),
			beego.NSRouter("/score/file/:user_id", &controllers.EducationScoreController{}, "get:Download"),

		),
		beego.NSNamespace("/case",
			beego.NSBefore(common.ChenkToken),
			//公安
			beego.NSRouter("/addcriminalcase", &controllers.CriminalCaseController{}, "post:AddCriminalCaseInfo"),
			beego.NSRouter("/list", &controllers.CriminalCaseController{}, "get:CriminalCaseList"),
			beego.NSRouter("/:id_card", &controllers.CriminalCaseController{}, "get:GetCriminalInfo"),
			beego.NSRouter("/:case_id", &controllers.CriminalCaseController{}, "put:UpdateCriminalCaseInfo"),
			beego.NSRouter("/:case_id", &controllers.CriminalCaseController{}, "delete:DeleteCriminalCaseInfo"),
			beego.NSRouter("/file/:user_id", &controllers.CriminalCaseController{}, "put:UpdateFile"),
			beego.NSRouter("/file/:user_id", &controllers.CriminalCaseController{}, "get:Download"),

		),
		beego.NSNamespace("/house",
			beego.NSBefore(common.ChenkToken),
			//住房
			beego.NSRouter("/addhouseinfo", &controllers.HousingInfoController{}, "post:AddHousingInfo"),
			beego.NSRouter("/list", &controllers.HousingInfoController{}, "get:HousingInfoList"),
			beego.NSRouter("/:id_card", &controllers.HousingInfoController{}, "get:GetHousingInfo"),
			beego.NSRouter("/:housing_id", &controllers.HousingInfoController{}, "put:UpdateHousingInfo"),
			beego.NSRouter("/:housing_id", &controllers.HousingInfoController{}, "delete:DeleteHousingInfo"),
			beego.NSRouter("/file/update/:user_id", &controllers.HousingInfoController{}, "put:UpdateFile"),
			beego.NSRouter("/file/download/:user_id", &controllers.HousingInfoController{}, "get:Download"),

		),
		beego.NSNamespace("/bank",
			beego.NSBefore(common.ChenkToken),
			//银行
			beego.NSRouter("/addbankinfo", &controllers.BankInfoController{}, "post:AddBankInfo"),
			beego.NSRouter("/list", &controllers.BankInfoController{}, "get:BankInfoList"),
			beego.NSRouter("/:id_card", &controllers.BankInfoController{}, "get:GetBankInfo"),
			beego.NSRouter("/:bank_card_no", &controllers.BankInfoController{}, "put:UpdateBankInfo"),
			beego.NSRouter("/:bank_card_no", &controllers.BankInfoController{}, "delete:DeleteBankInfo"),
			beego.NSRouter("/file/:user_id", &controllers.BankInfoController{}, "put:UpdateFile"),
			beego.NSRouter("/file/:user_id", &controllers.BankInfoController{}, "get:Download"),
			//银行明细
			beego.NSRouter("/addbankdetails", &controllers.BankDetailsController{}, "post:AddBankDetails"),
			beego.NSRouter("/details/list", &controllers.BankDetailsController{}, "get:BankDetailsList"),
			beego.NSRouter("/details/:bank_month_id", &controllers.BankDetailsController{}, "put:UpdateBankDetails"),
			beego.NSRouter("/details/:bank_month_id", &controllers.BankDetailsController{}, "delete:DeleteBankDetails"),
			beego.NSRouter("/details/file/:user_id", &controllers.BankDetailsController{}, "put:UpdateFile"),
			beego.NSRouter("/details/file/:user_id", &controllers.BankDetailsController{}, "get:Download"),
		),
		beego.NSNamespace("/admin",
			beego.NSBefore(common.ChenkToken),
			beego.NSRouter("/user/:login_user_id", &controllers.UserController{}, "put:UpdateUserStatus"),
			beego.NSRouter("/user/list", &controllers.UserController{}, "get:UserList"),
			beego.NSRouter("/user/status/list", &controllers.UserController{}, "get:UserStatusList"),
			beego.NSRouter("/user/:login_user_id", &controllers.UserController{}, "delete:DeleteUser"),
		),
	)
	beego.AddNamespace(ns)
}
