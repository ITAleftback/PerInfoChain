package main

import (
	_ "PerInfoChain/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
	"log"
)

func init() {
	//跨域
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))

	err := setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}
}
func main() {
	beego.Run()
}

func setupLogger() error {
	logSavePath := beego.AppConfig.String("LogSavePath")
	logFileName := beego.AppConfig.String("LogFileName")
	logFileExt := beego.AppConfig.String("LogFileExt")
	logPath := logSavePath + "/" + logFileName + logFileExt
	logs.SetLogger(logs.AdapterMultiFile, `{"filename": "`+logPath+`", "maxdays": 2, "maxlines": 1000000}`)
	return nil
}
