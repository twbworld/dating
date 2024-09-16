package dao

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/twbworld/dating/model/common"
	"github.com/twbworld/dating/model/db"
)

type DatingDb struct{}

func (d *DatingDb) GetDatingAmount(amount *[]uint, userId uint, tx ...*sqlx.Tx) error {
	sql := fmt.Sprintf("SELECT DISTINCT `d`.`id` FROM `%s` AS ut JOIN `%s` AS d ON `ut`.`dating_id` = `d`.`id` AND `ut`.`user_id` = ? AND `ut`.`status` = 1 WHERE `d`.`status` = 1 ORDER BY `d`.`id`", db.UserTime{}.TableName(), db.Dating{}.TableName())
	if len(tx) > 0 && tx[0] != nil {
		return tx[0].Select(amount, sql, userId)
	}
	return DB.Select(amount, sql, userId)
}

func (d *DatingDb) GetDating(dating *db.Dating, datingId uint, tx ...*sqlx.Tx) error {
	sql := fmt.Sprintf("SELECT `id`, `status`, `result`, `create_user_id` FROM `%s` WHERE `id` = ?", dating.TableName())
	if len(tx) > 0 && tx[0] != nil {
		return tx[0].Get(dating, sql, datingId)
	}
	return DB.Get(dating, sql, datingId)
}

func (d *DatingDb) GetUserTime(userTime *db.UserTime, datingId, UserId uint, tx ...*sqlx.Tx) error {
	sql := fmt.Sprintf("SELECT `id`, `dating_id`, `user_id` FROM `%s` WHERE `status` = 1 AND `dating_id` = ? AND `user_id` = ?", userTime.TableName())
	if len(tx) > 0 && tx[0] != nil {
		return tx[0].Get(userTime, sql, datingId, UserId)
	}
	return DB.Get(userTime, sql, datingId, UserId)
}

func (d *DatingDb) GetUserTimeById(userTime *db.UserTime, UtId uint, tx ...*sqlx.Tx) error {
	sql := fmt.Sprintf("SELECT `id`, `dating_id`, `user_id` FROM `%s` WHERE `status` = 1 AND `id` = ?", userTime.TableName())
	if len(tx) > 0 && tx[0] != nil {
		return tx[0].Get(userTime, sql, UtId)
	}
	return DB.Get(userTime, sql, UtId)
}

func (d *DatingDb) GetUserTimeByDatingId(userTime *[]db.UserTime, datingId uint, tx ...*sqlx.Tx) error {
	sql := fmt.Sprintf("SELECT `id`, `info` FROM `%s` WHERE `dating_id`= ? AND `status` = 1 ORDER BY `id`", db.UserTime{}.TableName())
	if len(tx) > 0 && tx[0] != nil {
		return tx[0].Select(userTime, sql, datingId)
	}
	return DB.Select(userTime, sql, datingId)
}

func (d *DatingDb) GetDatingUsers(datingUsers *[]common.DatingUser, datingId uint, tx ...*sqlx.Tx) error {
	sql := fmt.Sprintf("SELECT `u`.`id`, `u`.`nick_name`, `f`.`path`, `ut`.`id` AS ut_id, `ut`.`info` FROM `%s` AS ut JOIN `%s` AS u ON `ut`.`dating_id`= ? AND `ut`.`user_id` = `u`.`id` AND `ut`.`status` = 1 JOIN `%s` AS f ON `f`.`id` = `u`.`avatar` ORDER BY `ut`.`id`", db.UserTime{}.TableName(), db.User{}.TableName(), db.File{}.TableName())
	if len(tx) > 0 && tx[0] != nil {
		return tx[0].Select(datingUsers, sql, datingId)
	}
	return DB.Select(datingUsers, sql, datingId)
}

func (d *DatingDb) GetDatingList(datingList *[]common.DatingList, userId, lastId, limit uint, tx ...*sqlx.Tx) error {
	var sql string
	if lastId < 1 {
		sql = fmt.Sprintf("SELECT `d`.`id`, `d`.`create_user_id`, `d`.`status`,`ut`.`add_time`,`ut`.`id` AS ut_id FROM `%s` AS ut JOIN `%s` AS d ON `ut`.`dating_id` = `d`.`id` AND `ut`.`user_id` = ? AND `ut`.`status` = 1 ORDER BY `ut`.`id` DESC LIMIT ?", db.UserTime{}.TableName(), db.Dating{}.TableName())
		if len(tx) > 0 && tx[0] != nil {
			return tx[0].Select(datingList, sql, userId, limit)
		}
		return DB.Select(datingList, sql, userId, limit)
	} else {
		sql = fmt.Sprintf("SELECT `d`.`id`, `d`.`create_user_id`, `d`.`status`,`ut`.`add_time`,`ut`.`id` AS ut_id FROM `%s` AS ut JOIN `%s` AS d ON `ut`.`id` < ? AND `ut`.`dating_id` = `d`.`id` AND `ut`.`user_id` = ? AND `ut`.`status` = 1 ORDER BY `ut`.`id` DESC LIMIT ?", db.UserTime{}.TableName(), db.Dating{}.TableName())
		if len(tx) > 0 && tx[0] != nil {
			return tx[0].Select(datingList, sql, lastId, userId, limit)
		}
		return DB.Select(datingList, sql, lastId, userId, limit)
	}
}

func (d *DatingDb) GetUserByDatingIds(datingUserAvatar *[]common.DatingUserAvatar, datingIds []uint, tx ...*sqlx.Tx) error {
	if len(datingIds) < 1 {
		return nil
	}
	sql := fmt.Sprintf("SELECT `ut`.`dating_id`, `f`.`path` FROM `%s` AS ut JOIN `%s` AS u ON `ut`.`dating_id` IN (?) AND `ut`.`user_id` = `u`.`id` AND `ut`.`status` = 1 JOIN `%s` AS f ON `f`.`id` = `u`.`avatar`", db.UserTime{}.TableName(), db.User{}.TableName(), db.File{}.TableName())

	query, args, err := sqlx.In(sql, datingIds)
	if err != nil {
		return err
	}
	if len(tx) > 0 && tx[0] != nil {
		return tx[0].Select(datingUserAvatar, DB.Rebind(query), args...)
	}
	return DB.Select(datingUserAvatar, DB.Rebind(query), args...)
}

func (d *DatingDb) GetCleanDating(ids *[]uint, unix int64, tx ...*sqlx.Tx) error {
	sql := fmt.Sprintf("SELECT `id` FROM `%s` WHERE `update_time` < ?", db.Dating{}.TableName())
	if len(tx) > 0 && tx[0] != nil {
		return tx[0].Select(ids, sql, unix)
	}
	return DB.Select(ids, sql, unix)
}

func (d *DatingDb) CheckUserTime(datingUser *common.DatingUserJoin, utId uint, tx ...*sqlx.Tx) error {
	sql := fmt.Sprintf("SELECT `ut`.`user_id`, `ut`.`dating_id`, `d`.`create_user_id` FROM `%s` AS ut JOIN `%s` AS d ON `ut`.`dating_id` = `d`.`id` AND `ut`.`id` = ? AND `ut`.`status` = 1 AND `d`.`status` = 1", db.UserTime{}.TableName(), db.Dating{}.TableName())
	if len(tx) > 0 && tx[0] != nil {
		return tx[0].Get(datingUser, sql, utId)
	}
	return DB.Get(datingUser, sql, utId)
}

func (d *DatingDb) AddDating(userId uint, tx *sqlx.Tx) (uint, error) {
	if tx == nil {
		return 0, errors.New("请使用事务[iodhja]")
	}

	time := time.Now().Unix()
	sql, args := utils.getInsertSql(db.Dating{}, map[string]interface{}{
		"create_user_id": userId,
		"status":         1,
		"result":         "",
		"add_time":       time,
		"update_time":    time,
	})

	res, err := tx.Exec(sql, args...)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	return uint(id), err
}

func (d *DatingDb) JoinDating(UserTime *db.UserTime, tx *sqlx.Tx) (uint, error) {
	if tx == nil {
		return 0, errors.New("请使用事务[iodfgh]")
	}
	var (
		Id  uint
		err error
	)
	sql := fmt.Sprintf("SELECT `id` FROM `%s` WHERE `dating_id` = ? AND `user_id` = ? AND `status` = 0", UserTime.TableName())
	if tx.Get(&Id, sql, UserTime.DatingId, UserTime.UserId) == nil && Id > 0 {
		//曾经加入会面, 用回原来的数据
		return Id, d.JoinUpdate(Id, &db.UserTime{
			Info: UserTime.Info,
		}, tx)
	}

	time := time.Now().Unix()
	sql, args := utils.getInsertSql(UserTime, map[string]interface{}{
		"dating_id":   UserTime.DatingId,
		"user_id":     UserTime.UserId,
		"info":        UserTime.Info,
		"status":      1,
		"add_time":    time,
		"update_time": time,
	})

	res, err := tx.Exec(sql, args...)

	if err != nil {
		return 0, fmt.Errorf("[fuisdu]%s", err)
	}
	id, err := res.LastInsertId()
	return uint(id), err
}

func (d *DatingDb) JoinUpdate(utId uint, UserTime *db.UserTime, tx *sqlx.Tx) (err error) {
	if tx == nil {
		return errors.New("请使用事务[odfghja]")
	}

	var sql string
	if CanLock {
		sql = fmt.Sprintf("SELECT `id` FROM `%s` WHERE `id` = ? FOR UPDATE", UserTime.TableName())
		if _, err = tx.Exec(sql, utId); err != nil {
			return fmt.Errorf("[fuisdku]%s", err)
		}
	}

	sql, args := utils.getUpdateSql(UserTime, utId, map[string]interface{}{
		"info":        UserTime.Info,
		"status":      1,
		"update_time": time.Now().Unix(),
	})
	_, err = tx.Exec(sql, args...)

	return
}

func (d *DatingDb) CloseDating(datingIds []uint, tx *sqlx.Tx) error {
	if tx == nil {
		return errors.New("请使用事务[dfgshaaja]")
	}

	var sql string
	if CanLock {
		sql = fmt.Sprintf("SELECT `id` FROM `%s` WHERE `id` IN (?) FOR UPDATE", db.Dating{}.TableName())
		query, args, err := sqlx.In(sql, datingIds)
		if err != nil {
			return err
		}
		if _, err = tx.Exec(DB.Rebind(query), args...); err != nil {
			return err
		}
	}

	sql = fmt.Sprintf("UPDATE `%s` SET `status` = 0, `update_time` = ? WHERE `id` IN (?)", db.Dating{}.TableName())
	query, args, err := sqlx.In(sql, time.Now().Unix(), datingIds)
	if err != nil {
		return err
	}
	_, err = tx.Exec(DB.Rebind(query), args...)
	return err
}

func (d *DatingDb) QuitDating(utId uint, tx *sqlx.Tx) (err error) {
	if tx == nil {
		return errors.New("请使用事务[iodsghja]")
	}

	var sql string
	if CanLock {
		sql = fmt.Sprintf("SELECT `id` FROM `%s` WHERE `id` = ? FOR UPDATE", db.UserTime{}.TableName())
		if _, err = tx.Exec(sql, utId); err != nil {
			return
		}
	}

	sql, args := utils.getUpdateSql(db.UserTime{}, utId, map[string]interface{}{
		"status":      0,
		"update_time": time.Now().Unix(),
	})
	_, err = tx.Exec(sql, args...)

	return
}

func (d *DatingDb) DatingUpdate(datingId uint, res *db.DatingResult, tx *sqlx.Tx) (err error) {

	if tx == nil {
		return errors.New("请使用事务[dghja]")
	}

	var sql string
	if CanLock {
		sql = fmt.Sprintf("SELECT `id` FROM `%s` WHERE `id` = ? FOR UPDATE", db.Dating{}.TableName())
		if _, err = tx.Exec(sql, datingId); err != nil {
			return
		}
	}

	sql, args := utils.getUpdateSql(db.Dating{}, datingId, map[string]interface{}{
		"result":      res.Marshal(),
		"update_time": time.Now().Unix(),
	})

	_, err = tx.Exec(sql, args...)
	return
}
