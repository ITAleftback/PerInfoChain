package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB
var err error

func init() {
	mysqladmin := beego.AppConfig.String("mysqladmin")
	mysqlPwd := beego.AppConfig.String("mysqlPwd")
	mysqlDB := beego.AppConfig.String("mysqlDB")
	DB, err =
		gorm.Open("mysql", mysqladmin+":"+mysqlPwd+"@/"+mysqlDB+"?charset=utf8"+
			"&parseTime=True&loc=Local")
	if err != nil {
		logs.Error(err)
		logs.Error("连接MySql数据库失败")
	} else {
		logs.Info("连接MySql数据库成功")
	}
	DB.AutoMigrate(
		User{}, FabricUserInfo{}, EducationInfo{},
		EducationScoreInfo{}, HousingInfo{}, CriminalCase{},
		BankInfo{}, BankDetails{}, User{},
	)
}
