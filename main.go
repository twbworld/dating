package main

import (
	"fmt"

	"github.com/twbworld/dating/global"
	"github.com/twbworld/dating/initialize"
	initGlobal "github.com/twbworld/dating/initialize/global"
	"github.com/twbworld/dating/initialize/system"
	"github.com/twbworld/dating/task"
)

func main() {
	initGlobal.New().Start()
	initialize.InitializeLogger()
	if err := system.DbStart(); err != nil {
		global.Log.Fatalf("连接数据库失败[fbvk89]: %v", err)
	}
	defer system.DbClose()

	// service.Service.UserServiceGroup.DatingService.Match(4)

	defer func() {
		if p := recover(); p != nil {
			global.Log.Println(p)
		}
	}()

	switch initGlobal.Act {
	case "":
		initialize.Start()
	case "clear":
		task.Clear()
	default:
		fmt.Println("参数可选: clear|expiry")
	}

}
