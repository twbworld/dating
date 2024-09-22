package main

import (
	"github.com/twbworld/dating/global"
	"github.com/twbworld/dating/initialize"
	initGlobal "github.com/twbworld/dating/initialize/global"
	"github.com/twbworld/dating/initialize/system"
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
	initialize.Start()
}
