package global

import (
	"fmt"
	"time"

	"github.com/twbworld/dating/global"
)

func (*GlobalInit) initTz() error {
	Location, err := time.LoadLocation(global.Config.Tz)
	if err != nil {
		return fmt.Errorf("时区配置失败[siortuj]: %w", err)
	}
	global.Tz = Location
	return nil
}
