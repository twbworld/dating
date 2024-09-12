package db

import "encoding/json"

type UserTime struct {
	BaseField
	DatingId uint   `db:"dating_id" json:"dating_id" info:"dating表id"`
	UserId   uint   `db:"user_id" json:"user_id" info:"user表id"`
	Info     string `db:"info" json:"info" info:"空闲时间信息;{'t': [[1706978785,1706978785]]}"`
	Status   int8   `db:"status" json:"status" info:"状态; 0:已退出;1:加入"`
}

func (UserTime) TableName() string {
	return `user_time`
}

// json转结构体
func (d *UserTime) InfoUnmarshal() *UserTimeInfo {
	result := &UserTimeInfo{
		Time: make([][2]int64, 0),
	}
	if d.Info == "" {
		return result
	}
	if err := json.Unmarshal([]byte(d.Info), result); err == nil {
		return result
	}
	//避免转json后,属性为"null"
	result.Time = make([][2]int64, 0)
	return result
}

type UserTimeInfo struct {
	Time [][2]int64 `json:"t"`
}

// 结构体转json
func (u *UserTimeInfo) Marshal() string {
	if u.Time == nil {
		//避免转json后,属性为"null"
		u.Time = make([][2]int64, 0)
	}

	result, err := json.Marshal(u)
	if err != nil {
		return ``
	}
	return string(result)
}
