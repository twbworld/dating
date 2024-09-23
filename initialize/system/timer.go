package system

import (
	"github.com/robfig/cron/v3"
	"github.com/twbworld/dating/global"
	"github.com/twbworld/dating/service"
	"github.com/twbworld/dating/task"
)

var c *cron.Cron

// startCronJob 启动一个新的定时任务
func startCronJob(schedule string, task func() error, name string) error {
	_, err := c.AddFunc(schedule, func() {
		defer func() {
			text := "任务完成"
			if p := recover(); p != nil {
				text = "任务出错[gqnoj]: " + p.(string)
			}
			service.Service.UserServiceGroup.TgService.TgSend(name + text)
		}()
		if err := task(); err != nil {
			panic(err)
		}
	})
	return err
}

func timerStart() error {
	c = cron.New([]cron.Option{
		cron.WithLocation(global.Tz),
		// cron.WithSeconds(), //精确到秒
	}...)

	if err := startCronJob("0 3 * * *", task.Clear, "清除过期数据"); err != nil {
		return err
	}

	c.Start() //已含协程
	global.Log.Infoln("定时器启动成功")
	return nil
}

func timerStop() error {
	if c == nil {
		global.Log.Warnln("定时器未启动")
		return nil
	}
	c.Stop()
	global.Log.Infoln("定时器停止成功")
	return nil
}
