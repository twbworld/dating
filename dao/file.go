package dao

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/twbworld/dating/model/db"
)

type FileDb struct{}

// 新增文件记录
// 要先检查文件是否存在于目录, 以及不存在于数据库
func (u *FileDb) AddFile(f *db.File, tx *sqlx.Tx) (uint, error) {
	if tx == nil {
		return 0, errors.New(`请使用事务[iodhja]`)
	}
	sql, args := utils.getInsertSql(f, map[string]interface{}{
		"path":     f.Path,
		"ext":      f.Ext,
		"add_time": time.Now().Unix(),
	})

	res, err := tx.Exec(sql, args...)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	return uint(id), err
}

func (u *FileDb) GetFile(file *db.File, fileId uint, tx ...*sqlx.Tx) error {
	sql := fmt.Sprintf("SELECT `id`, `path`, `ext` FROM `%s` WHERE `id` = ?", file.TableName())
	if len(tx) > 0 && tx[0] != nil {
		return tx[0].Get(file, sql, fileId)
	}
	return DB.Get(file, sql, fileId)
}

func (u *FileDb) GetFileId(path string, tx ...*sqlx.Tx) (uint, error) {
	var file db.File
	sql := fmt.Sprintf("SELECT `id` FROM `%s` WHERE `path` = ?", file.TableName())
	if len(tx) > 0 && tx[0] != nil {
		err := tx[0].Get(&file, sql, path)
		return file.Id, err
	}
	err := DB.Get(&file, sql, path)
	return file.Id, err
}
