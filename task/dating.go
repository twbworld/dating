package task

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/twbworld/dating/dao"
)

const cleanDaysAgo = 7 //清除?天前的数据

func Clean() error {
	cutoffTime, ids := time.Now().AddDate(0, 0, -cleanDaysAgo), []uint{}
	if err := dao.App.DatingDb.GetCleanDating(&ids, cutoffTime.Unix()); err != nil || len(ids) < 1 {
		return err
	}

	return dao.Tx(func(tx *sqlx.Tx) error {
		return dao.App.DatingDb.CloseDating(ids, tx)
	})

}
