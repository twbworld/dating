package task

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/twbworld/dating/dao"
	"github.com/twbworld/dating/global"
)

const cleanDaysAgo = 7 //清除?天前的数据

func Clear() error {
	cutoffTime, ids := time.Now().AddDate(0, 0, -cleanDaysAgo), []uint{}
	if err := dao.App.DatingDb.GetCleanDating(&ids, cutoffTime.Unix()); err != nil {
		global.Log.Errorln("清除失败[sdfsjn]", err)
		return err
	}
	if len(ids) < 1 {
		global.Log.Infoln("不需清理")
		return nil
	}
	err := dao.Tx(func(tx *sqlx.Tx) error {
		return dao.App.DatingDb.CloseDating(ids, tx)
	})
	if err != nil {
		global.Log.Errorln("清除失败[s8jn]", ids)
		return err
	}

	global.Log.Infoln("成功清除过期数据", ids)
	return nil
}
