package main

import (
	"log"

	"github.com/twbworld/dating/global"
	"github.com/twbworld/dating/initialize"
	initGlobal "github.com/twbworld/dating/initialize/global"
	"github.com/twbworld/dating/initialize/system"
	"github.com/twbworld/dating/task"
)

func main() {
	initGlobal.New().Start()
	initialize.InitializeLogger()
	sys := system.Start()
	defer sys.Stop()

	defer func() {
		if p := recover(); p != nil {
			global.Log.Println(p)
		}
	}()

	// service.Service.UserServiceGroup.DatingService.Match(4)
	switch initGlobal.Act {
	case "":
		initialize.Start()
	case "clear":
		task.Clear()
	default:
		log.Println("参数可选: clear")
	}

}
