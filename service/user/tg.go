package user

import (
	"fmt"
	"strings"
	"sync"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/twbworld/dating/global"
)

type TgService struct{}

var lock sync.RWMutex

// 向Tg发送信息(请用协程执行)
func (t *TgService) TgSend(text string) (err error) {
	if global.Bot == nil {
		return fmt.Errorf("[ertioj98]出错")
	}

	if len(text) < 1 {
		return fmt.Errorf("[sioejn89]出错")
	}

	lock.RLock()
	defer lock.RUnlock()

	var str strings.Builder
	str.WriteString(`[`)
	str.WriteString(global.Config.ProjectName)
	str.WriteString(`]`)
	str.WriteString(text)

	msg := tg.NewMessage(global.Config.Telegram.Id, str.String())
	// msg.ParseMode = "MarkdownV2" //使用Markdown格式, 需要对特殊字符进行转义

	_, err = global.Bot.Send(msg)
	return
}
