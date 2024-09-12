package dao

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/twbworld/dating/model/db"
)

type UserDb struct{}

func (u *UserDb) GetUserById(user *db.User, id uint, tx ...*sqlx.Tx) error {
	sql := fmt.Sprintf("SELECT `id`, `avatar`, `nick_name`, `gender`, `phone_number`, `session_key`, `add_time` FROM `%s` WHERE `id` = ?", user.TableName())
	if len(tx) > 0 && tx[0] != nil {
		return tx[0].Get(user, sql, id)
	}
	return DB.Get(user, sql, id)
}

func (u *UserDb) GetUserByOpenId(user *db.User, openId string, tx ...*sqlx.Tx) error {
	sql := fmt.Sprintf("SELECT `id`, `avatar`, `nick_name`, `gender`, `phone_number`, `session_key`, `add_time` FROM `%s` WHERE `openid` = ?", user.TableName())
	if len(tx) > 0 && tx[0] != nil {
		return tx[0].Get(user, sql, openId)
	}
	return DB.Get(user, sql, openId)
}

func (u *UserDb) AddUser(user *db.User, tx *sqlx.Tx) (uint, error) {
	if tx == nil {
		return 0, errors.New(`请使用事务[iodhja]`)
	}
	currentTime := time.Now().Unix()
	sql, args := utils.getInsertSql(user, map[string]interface{}{
		"nick_name":   user.NickName,
		"avatar":      user.Avatar,
		"gender":      user.Gender,
		"openid":      user.OpenId,
		"unionid":     user.UnionId,
		"session_key": user.SessionKey,
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

func (u *UserDb) UpdateUser(user *db.User, tx *sqlx.Tx) (err error) {
	if tx == nil {
		return errors.New("请使用事务[ologhja]")
	}
	if user.Id < 1 {
		return errors.New("参数错误[gdfkmo90]")
	}

	var sql string
	if CanLock {
		sql = fmt.Sprintf("SELECT `id` FROM `%s` WHERE `id` = ? FOR UPDATE", db.User{}.TableName())
		if _, err = tx.Exec(sql, user.Id); err != nil {
			return
		}
	}

	sql, args := utils.getUpdateSql(db.User{}, user.Id, map[string]interface{}{
		"nick_name": user.NickName,
		"avatar":    user.Avatar,
	})
	_, err = tx.Exec(sql, args...)

	return
}

func (u *UserDb) UpdateSessionKey(userId uint, sessionKey string, tx *sqlx.Tx) (err error) {
	if tx == nil {
		return errors.New("请使用事务[ologhja]")
	}
	if userId < 1 {
		return errors.New("参数错误[olshja]")
	}

	var sql string
	if CanLock {
		sql = fmt.Sprintf("SELECT `id` FROM `%s` WHERE `id` = ? FOR UPDATE", db.User{}.TableName())
		if _, err = tx.Exec(sql, userId); err != nil {
			return
		}
	}

	sql, args := utils.getUpdateSql(db.User{}, userId, map[string]interface{}{
		"session_key": sessionKey,
	})
	_, err = tx.Exec(sql, args...)

	return
}
