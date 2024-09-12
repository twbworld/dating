package system

import (
	"github.com/robfig/cron/v3"
	"github.com/twbworld/dating/global"
	"github.com/twbworld/dating/task"
)

var c *cron.Cron

func timerStart() error {
	var option []cron.Option
	// option = append(option, cron.WithSeconds()) //精确到秒
	c = cron.New(option...)

	_, err := c.AddFunc("0 3 * * *", func() {
		if err := task.Clean(); err != nil {
			global.Log.Errorf("任务出错[osjd]: %s", err)
		}
	})
	if err != nil {
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
