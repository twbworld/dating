package global

import (
	"fmt"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram"
	"github.com/twbworld/dating/global"
)

func (*GlobalInit) initMiniProgram() error {
	var err error
	global.MiniProgramApp, err = miniProgram.NewMiniProgram(&miniProgram.UserConfig{
		AppID:     global.Config.Weixin.XcxAppid,
		Secret:    global.Config.Weixin.XcxSecret,
		Debug:     global.Config.Debug,
		HttpDebug: global.Config.Debug,
		Log: miniProgram.Log{
			File:  `./log/wechat.log`,
			Error: `./log/wechat_error.log`,
		},
	})

	if err != nil {
		return fmt.Errorf("小程序配置错误[opjkgh]: %w", err)
	}
	return nil
}
