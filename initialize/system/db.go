package system

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/twbworld/dating/dao"
	"github.com/twbworld/dating/global"
	"github.com/twbworld/dating/model/db"
	"github.com/twbworld/dating/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type mysql struct{}
type sqlite struct{}
type class interface {
	connect() error
	createTable() error
	insertData(string, *sqlx.Tx) error
	version() string
}

func DbStart() error {
	var dbRes class

	switch global.Config.Database.Type {
	case "mysql":
		dbRes = &mysql{}
	case "sqlite":
		dbRes = &sqlite{}
	default:
		dbRes = &sqlite{}
	}

	if err := dbRes.connect(); err != nil {
		return err
	}
	dbRes.createTable()
	return nil
}

// 关闭数据库连接
func DbClose() error {
	if dao.DB != nil {
		return dao.DB.Close()
	}
	return nil
}

// 连接SQLite数据库
func (s *sqlite) connect() error {
	var err error

	if dao.DB, err = sqlx.Open("sqlite3", global.Config.Database.SqlitePath); err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}
	//没有数据库会创建
	if err = dao.DB.Ping(); err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	dao.DB.SetMaxOpenConns(16)
	dao.DB.SetMaxIdleConns(8)
	dao.DB.SetConnMaxLifetime(time.Minute * 5)

	//提高并发
	if _, err = dao.DB.Exec("PRAGMA journal_mode = WAL"); err != nil {
		return fmt.Errorf("数据库设置失败: %w", err)
	}
	//超时等待
	if _, err = dao.DB.Exec("PRAGMA busy_timeout = 10000;"); err != nil {
		return fmt.Errorf("数据库设置失败: %w", err)
	}
	// 设置同步模式为 NORMAL
	if _, err = dao.DB.Exec("PRAGMA synchronous = NORMAL;"); err != nil {
		return fmt.Errorf("数据库设置失败: %w", err)
	}

	dao.CanLock = false

	global.Log.Infof("%s版本: %s; 地址: %s", global.Config.Database.Type, s.version(), global.Config.Database.SqlitePath)
	return nil
}

func (m *mysql) connect() error {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", global.Config.Database.MysqlUsername, global.Config.Database.MysqlPassword, global.Config.Database.MysqlHost, global.Config.Database.MysqlPort, global.Config.Database.MysqlDbname)

	//也可以使用MustConnect连接不成功就panic
	if dao.DB, err = sqlx.Connect("mysql", dsn); err != nil {
		return fmt.Errorf("数据库连接失败[rwbhe3]: %s\n%w", dsn, err)
	}

	dao.DB.SetMaxOpenConns(16)
	dao.DB.SetMaxIdleConns(8)
	dao.DB.SetConnMaxLifetime(time.Minute * 5) // 设置连接的最大生命周期

	if err = dao.DB.Ping(); err != nil {
		return fmt.Errorf("数据库连接失败: %s\n%w", dsn, err)
	}

	dao.CanLock = true
	global.Log.Infof("%s版本: %s; 地址: @tcp(%s:%s)/%s", global.Config.Database.Type, m.version(), global.Config.Database.MysqlHost, global.Config.Database.MysqlPort, global.Config.Database.MysqlDbname)
	return nil
}

func (s *sqlite) createTable() error {
	var u []string
	err := dao.DB.Select(&u, "SELECT name _id FROM sqlite_master WHERE type ='table'")
	if err != nil {
		return fmt.Errorf("查询表失败: %w", err)
	}

	sqls := map[string][]string{
		db.Dating{}.TableName(): {
			`CREATE TABLE "%s" ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "create_user_id" INTEGER NOT NULL DEFAULT 0, "result" TEXT NOT NULL DEFAULT '', "status" INTEGER(1) NOT NULL DEFAULT 0, "add_time" TEXT(10) NOT NULL DEFAULT '', "update_time" TEXT(10) NOT NULL DEFAULT '');`,
			`CREATE INDEX "idx_create_user_id" ON "%s" ("create_user_id" ASC);`,
		},
		db.User{}.TableName(): {
			`CREATE TABLE "%s" ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "nick_name" TEXT(64) NOT NULL DEFAULT '', "avatar" INTEGER NOT NULL DEFAULT 0, "gender" INTEGER(1) NOT NULL DEFAULT 0, "phone_number" TEXT(20) NOT NULL DEFAULT '', "openid" TEXT NOT NULL DEFAULT '', "unionid" TEXT NOT NULL DEFAULT '', "session_key" TEXT NOT NULL DEFAULT '', "add_time" TEXT(10) NOT NULL DEFAULT '', "update_time" TEXT(10) NOT NULL DEFAULT '');`,
			`CREATE INDEX "idx_avatar" ON "%s" ("avatar" ASC);`,
			`CREATE INDEX "idx_openid" ON "%s" ("openid" ASC);`,
			`CREATE INDEX "idx_unionid" ON "%s" ("unionid" ASC);`,
		},
		db.UserTime{}.TableName(): {
			`CREATE TABLE "%s" ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "dating_id" INTEGER NOT NULL DEFAULT 0 , "user_id" INTEGER NOT NULL DEFAULT 0, "info" TEXT NOT NULL DEFAULT '', "status" INTEGER(1) NOT NULL DEFAULT 0, "add_time" TEXT(10) NOT NULL DEFAULT '', "update_time" TEXT(10) NOT NULL DEFAULT '');`,
			`CREATE INDEX "idx_dating_id" ON "%s" ("dating_id" ASC);`,
			`CREATE INDEX "idx_user_id" ON "%s" ("user_id" ASC);`,
		},
		db.Feedback{}.TableName(): {
			`CREATE TABLE "%s" ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "user_id" INTEGER NOT NULL DEFAULT 0, "desc" TEXT(255) NOT NULL DEFAULT '', "file_id" TEXT(255) NOT NULL DEFAULT '', "add_time" TEXT(10) NOT NULL DEFAULT '', "update_time" TEXT(10) NOT NULL DEFAULT '');`,
			`CREATE INDEX "idx_feedback_user_id" ON "%s" ("user_id" ASC);`,
		},
		db.File{}.TableName(): {
			`CREATE TABLE "%s" ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "path" TEXT NOT NULL DEFAULT '', "ext" TEXT(10) NOT NULL DEFAULT '', "type" INTEGER(1) NOT NULL DEFAULT 0, "add_time" TEXT(10) NOT NULL DEFAULT '');`,
			`CREATE INDEX "idx_path" ON "%s" ("path" ASC);`,
		},
	}

	err = dao.Tx(func(tx *sqlx.Tx) (e error) {
		for k, v := range sqls {
			if utils.InSlice(u, k) < 0 {
				for _, val := range v {
					if _, e := tx.Exec(fmt.Sprintf(val, k)); e != nil {
						return fmt.Errorf("错误[ghjbcvgs]:  %s\n%w", val, e)
					}
				}
				if err := s.insertData(k, tx); err != nil {
					return fmt.Errorf("插入数据失败: %s\n%w", k, err)
				}
				global.Log.Infof("创建%s表[dkyjh]", k)
			}
		}
		return
	})
	if err != nil {
		return fmt.Errorf("事务执行失败: %w", err)
	}
	return nil
}

func (m *mysql) createTable() error {
	var u []string
	err := dao.DB.Select(&u, "SHOW TABLES")
	if err != nil {
		return fmt.Errorf("插入数据失败: %w", err)
	}

	sqls := map[string]string{
		db.Dating{}.TableName():   "CREATE TABLE `%s` (`id` int unsigned NOT NULL AUTO_INCREMENT, `create_user_id` int unsigned NOT NULL DEFAULT '0' COMMENT '会面创建者', `result` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '推荐结果', `status` tinyint NOT NULL DEFAULT '0' COMMENT '会面状态; 0:结束;1:进行中;', `add_time` int unsigned NOT NULL DEFAULT '0', `update_time` int unsigned NOT NULL DEFAULT '0', PRIMARY KEY (`id`), KEY `idx_create_user_id` (`create_user_id`) USING BTREE ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='会面表';",
		db.User{}.TableName():     "CREATE TABLE `%s` (`id` int unsigned NOT NULL AUTO_INCREMENT, `nick_name` varchar(64) NOT NULL DEFAULT '' COMMENT '昵称', `avatar` int unsigned NOT NULL DEFAULT '0' COMMENT '头像, 关联file表', `gender` tinyint NOT NULL DEFAULT '0' COMMENT '性别;0:未知;1:男;2:女', `phone_number` varchar(20) NOT NULL DEFAULT '' COMMENT '手机号', `openid` varchar(255) NOT NULL DEFAULT '' COMMENT '小程序平台的用户识别码', `unionid` varchar(255) NOT NULL DEFAULT '' COMMENT '微信用户识别码', `session_key` varchar(255) NOT NULL DEFAULT '' COMMENT '微信Session_Key', `add_time` int unsigned NOT NULL DEFAULT '0', `update_time` int unsigned NOT NULL DEFAULT '0', PRIMARY KEY (`id`), KEY `idx_avatar` (`avatar`), KEY `idx_openid` (`openid`), KEY `idx_unionid` (`unionid`) ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';",
		db.UserTime{}.TableName(): "CREATE TABLE `%s` (`id` int unsigned NOT NULL AUTO_INCREMENT, `dating_id` int unsigned NOT NULL DEFAULT '0' COMMENT 'dating表id', `user_id` int unsigned NOT NULL DEFAULT '0' COMMENT 'user表id', `info` text NOT NULL COMMENT '空闲时间信息;{''t'': [[1706978785,1706978785]]}', `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态; 0:已退出;1:加入', `add_time` int unsigned NOT NULL DEFAULT '0', `update_time` int unsigned NOT NULL DEFAULT '0', PRIMARY KEY (`id`), KEY `idx_dating_id` (`dating_id`), KEY `idx_user_id` (`user_id`) ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户时间表';",
		db.Feedback{}.TableName(): "CREATE TABLE `%s` (`id` int unsigned NOT NULL AUTO_INCREMENT, `user_id` int unsigned NOT NULL DEFAULT '0' COMMENT '反馈人', `desc` varchar(255) NOT NULL DEFAULT '' COMMENT '反馈详情', `file_id` varchar(255) NOT NULL DEFAULT '' COMMENT 'file表id, 逗号相隔', `add_time` int unsigned NOT NULL DEFAULT '0', `update_time` int unsigned NOT NULL DEFAULT '0', PRIMARY KEY (`id`), KEY `idx_user_id` (`user_id`) ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='反馈信息表';",
		db.File{}.TableName():     "CREATE TABLE `%s` (`id` int unsigned NOT NULL AUTO_INCREMENT, `path` varchar(255) NOT NULL DEFAULT '' COMMENT '文件路径', `ext` varchar(10) NOT NULL DEFAULT '' COMMENT '文件类型,如jpg/mp4等', `type` tinyint NOT NULL DEFAULT '0' COMMENT '类型;0:本地;1:远程(如cdn)', `add_time` int unsigned NOT NULL DEFAULT '0', PRIMARY KEY (`id`), KEY `idx_path` (`path`) ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='媒体文件表';",
	}

	err = dao.Tx(func(tx *sqlx.Tx) (e error) {
		for k, v := range sqls {
			if utils.InSlice(u, k) < 0 {
				if _, e := tx.Exec(fmt.Sprintf(v, k)); e != nil {
					return fmt.Errorf("插入数据失败: %s\n%w", k, err)
				}
				global.Log.Infof("创建%s表[dfsjh]", k)
				if err := m.insertData(k, tx); err != nil {
					return fmt.Errorf("插入数据失败[fnko9]: %s\n%w", k, err)
				}
			}
		}
		return
	})
	if err != nil {
		return fmt.Errorf("事务执行失败: %w", err)
	}
	return nil
}

func (m *mysql) insertData(t string, tx *sqlx.Tx) error {
	return insert(t, tx)
}

func (m *sqlite) insertData(t string, tx *sqlx.Tx) error {
	return insert(t, tx)
}

func insert(t string, tx *sqlx.Tx) error {

	now := time.Now()
	ti := now.Unix()

	var sqls []string

	switch t {
	case db.Dating{}.TableName():
		if gin.Mode() == gin.TestMode {
			ti2 := now.AddDate(0, 0, -1).Unix()
			sqls = []string{
				fmt.Sprintf("INSERT INTO `%s`(`create_user_id`, `status`, `result`, `add_time`, `update_time`) VALUES(3, 1, '', %d, %d)", db.Dating{}.TableName(), ti2, ti2),
			}
		} else if gin.Mode() == gin.DebugMode {
			ti2 := now.AddDate(0, 0, -1).Unix()
			sqls = []string{
				fmt.Sprintf("INSERT INTO `%s`(`create_user_id`, `status`, `result`, `add_time`, `update_time`) VALUES(2, 1, '', %d, %d)", db.Dating{}.TableName(), ti2, ti2),
				fmt.Sprintf("INSERT INTO `%s`(`create_user_id`, `status`, `result`, `add_time`, `update_time`) VALUES(3, 1, '', %d, %d)", db.Dating{}.TableName(), ti2, ti2),
			}
		}
	case db.User{}.TableName():
		n := db.User{}.TableName()
		sqls = []string{
			fmt.Sprintf("INSERT INTO `%s`(`id`, `nick_name`, `avatar`, `add_time`, `update_time`) VALUES(%d, '手动', 1, %d, %d)", n, dao.BaseUserId, ti, ti),
		}
		if gin.Mode() == gin.TestMode {
			sqls = append(sqls,
				fmt.Sprintf("INSERT INTO `%s`(`nick_name`, `avatar`, `add_time`, `update_time`) VALUES('test', 1, %d, %d)", n, ti, ti),
			)
		} else if gin.Mode() == gin.DebugMode {
			sqls = append(sqls,
				fmt.Sprintf("INSERT INTO `%s`(`nick_name`, `avatar`, `add_time`, `update_time`) VALUES('test', 1, %d, %d)", n, ti, ti),
				fmt.Sprintf("INSERT INTO `%s`(`nick_name`, `avatar`, `add_time`, `update_time`) VALUES('test2', 1, %d, %d)", n, ti, ti),
			)
		}
	case db.UserTime{}.TableName():
		if gin.Mode() == gin.TestMode {

			type jt struct {
				Time [][2]int64 `json:"t"`
			}

			n := db.UserTime{}.TableName()
			t1, _ := time.ParseInLocation("2006-01-02 15", "2024-02-15 09", global.Tz)
			t2, _ := time.ParseInLocation("2006-01-02 15", "2024-02-15 22", global.Tz)
			t3, _ := time.ParseInLocation("2006-01-02 15", "2024-02-18 08", global.Tz)
			t4, _ := time.ParseInLocation("2006-01-02 15", "2024-02-18 23", global.Tz)
			a, _ := json.Marshal(jt{Time: [][2]int64{{t1.Unix(), t2.Unix()}, {t3.Unix(), t4.Unix()}}})

			t5, _ := time.ParseInLocation("2006-01-02 15", "2024-02-15 08", global.Tz)
			t6, _ := time.ParseInLocation("2006-01-02 15", "2024-02-15 21", global.Tz)
			b, _ := json.Marshal(jt{Time: [][2]int64{{t5.Unix(), t6.Unix()}}})

			t7, _ := time.ParseInLocation("2006-01-02 15", "2024-02-18 08", global.Tz)
			t8, _ := time.ParseInLocation("2006-01-02 15", "2024-02-18 23", global.Tz)
			t9, _ := time.ParseInLocation("2006-01-02 15", "2024-02-19 10", global.Tz)
			t10, _ := time.ParseInLocation("2006-01-02 15", "2024-02-19 22", global.Tz)
			c, _ := json.Marshal(jt{Time: [][2]int64{{t7.Unix(), t8.Unix()}, {t9.Unix(), t10.Unix()}}})

			t11, _ := time.ParseInLocation("2006-01-02 15", "2024-02-19 22", global.Tz)
			t12, _ := time.ParseInLocation("2006-01-02 15", "2024-02-19 23", global.Tz)
			d, _ := json.Marshal(jt{Time: [][2]int64{{t11.Unix(), t12.Unix()}}})

			sqls = []string{
				fmt.Sprintf("INSERT INTO `%s`(`dating_id`, `user_id`, `info`, `status`, `add_time`, `update_time`) VALUES(1, 1, '%s', 1, %d, %d)", n, a, ti, ti),
				fmt.Sprintf("INSERT INTO `%s`(`dating_id`, `user_id`, `info`, `status`, `add_time`, `update_time`) VALUES(1, 1, '%s', 1, %d, %d)", n, b, ti, ti),
				fmt.Sprintf("INSERT INTO `%s`(`dating_id`, `user_id`, `info`, `status`, `add_time`, `update_time`) VALUES(1, 3, '%s', 1, %d, %d)", n, d, ti, ti),
				fmt.Sprintf("INSERT INTO `%s`(`dating_id`, `user_id`, `info`, `status`, `add_time`, `update_time`) VALUES(1, 1, '%s', 1, %d, %d)", n, c, ti, ti),
			}
		}
	case db.File{}.TableName():
		n := db.File{}.TableName()
		if gin.Mode() == gin.TestMode {
			sqls = []string{
				fmt.Sprintf("INSERT INTO `%s`(`path`, `ext`, `add_time`) VALUES('static/favicon.ico', 'png', %d)", n, ti),
			}
		} else if gin.Mode() == gin.DebugMode {
		}
	}

	for _, v := range sqls {
		global.Log.Infof("创建数据[dfskkjh]%s", v)
		if _, e := tx.Exec(v); e != nil {
			return fmt.Errorf("错误[gh90iggs]: %s\n%w", v, e)
		}
	}
	return nil
}

func (*sqlite) version() (t string) {
	dao.DB.Get(&t, `SELECT sqlite_version()`)
	return
}

func (*mysql) version() (t string) {
	dao.DB.Get(&t, `SELECT version()`)
	return
}
