package global

import (
	"time"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/twbworld/dating/model/config"
)

// 全局变量
// 业务逻辑禁止修改
var (
	Config         *config.Config = new(config.Config) //指针类型, 给与其内存空间
	Log            *logrus.Logger
	Tz             *time.Location
	MiniProgramApp *miniProgram.MiniProgram
	Bot            *tg.BotAPI
)
