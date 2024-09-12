package common

import (
	"errors"
	"time"

	"github.com/twbworld/dating/global"
	"github.com/twbworld/dating/model/db"
)

type InfoPost struct {
	Time [][2]string `json:"t"`
}

type DatingPost struct {
	Id   uint     `json:"id"`
	Info InfoPost `json:"info"`
}

type UerTimePost struct {
	UtId uint     `json:"ut_id"`
	Info InfoPost `json:"info"`
}

type LoginPost struct {
	Code string `json:"code"`
}

type GetDatingPost struct {
	Id uint `json:"id"`
}

type GetDatingListPost struct {
	Page   uint `json:"page"`
	LastId uint `json:"last_id"` //上次获取最后的ut_id
}

type QuitDatingPost struct {
	Id   uint `json:"id"` //优先此值
	UtId uint `json:"ut_id"`
}

type UserInfoPost struct {
	LoginPost
	NickName string `json:"nick_name"` //昵称
	// File     string `json:"file"`      //头像文件
}

type FeedbackPost struct {
	Desc string `json:"desc"`
	// file []string `json:"file"`
}

func (i *InfoPost) Marshal() (string, error) {
	var info db.UserTimeInfo
	for _, value := range i.Time {
		te, err := time.ParseInLocation(time.DateTime, value[1], global.Tz)
		if err != nil {
			return ``, errors.New(`参数错误[ds6hgfesd]`)
		}
		ts, err := time.ParseInLocation(time.DateTime, value[0], global.Tz)
		if err != nil {
			return ``, errors.New(`参数错误[ds6hgfesd]`)
		}
		info.Time = append(info.Time, [2]int64{ts.Unix(), te.Unix()})
	}

	if len(info.Time) == 0 {
		return ``, errors.New(`参数错误[ds6sesd]`)
	}

	//时间排序
	for i := range len(info.Time) {
		val := info.Time[i]
		j := i - 1
		for ; j >= 0 && info.Time[j][0] > val[0]; j-- {
			info.Time[j+1] = info.Time[j]
		}
		info.Time[j+1] = val
	}

	return info.Marshal(), nil
}
