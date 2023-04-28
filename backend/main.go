package main

import (
	"LibManSys/app"
	"LibManSys/conf"
	"LibManSys/web"
)

func main() {
	conf.Init()

	app.LMS = app.NewLibraryManagementSystemImpl(conf.GetMysqlLoginConfig())
	app.LMS.Init()
	defer app.LMS.Free()

	web.InitWebFramework()
	web.StartServer()
}
