package dao

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/twbworld/dating/model/db"
)

type FeedbackDb struct{}

func (u *FeedbackDb) AddFeedback(f *db.Feedback, tx *sqlx.Tx) (uint, error) {
	if tx == nil {
		return 0, errors.New(`请使用事务[iodhja]`)
	}

	currentTime := time.Now().Unix()
	sql, args := utils.getInsertSql(f, map[string]interface{}{
		"user_id":     f.UserId,
		"desc":        f.Desc,
		"file_id":     f.FileId,
		"add_time":    currentTime,
		"update_time": currentTime,
	})

	res, err := tx.Exec(sql, args...)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	return uint(id), err
}
