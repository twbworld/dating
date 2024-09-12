package user

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/twbworld/dating/dao"
	"github.com/twbworld/dating/global"
	"github.com/twbworld/dating/model/db"
	"github.com/twbworld/dating/utils"
)

const (
	minTime, maxTime = 8, 23
)

func (d *DatingService) MatchGoroutine(datingId uint) {
	if datingId < 1 {
		return
	}
	go func() {
		defer func() {
			//避免协程内panic影响到外层主程序
			if p := recover(); p != nil {
				global.Log.Error("[gji075]", datingId, p)
			}
		}()
		if _, err := d.Match(datingId); err != nil {
			global.Log.Error("[gjs075]", datingId, err)
		}
	}()
}

// 匹配最合适的时间
func (d *DatingService) Match(datingId uint) (res db.DatingResult, err error) {
	if datingId < 1 {
		return res, errors.New("[odhkjgo]")
	}
	var du []db.UserTime
	if err := dao.App.DatingDb.GetUserTimeByDatingId(&du, datingId); err != nil {
		return res, fmt.Errorf("[5634go]%s", err)
	}

	duLen := len(du)
	if duLen < 1 {
		return res, errors.New("[iudha09]")
	}

	timeIds := d.aggregateUserTimes(du)
	if len(timeIds) < 1 {
		return res, errors.New("[iudhla09]")
	}

	numTime := d.aggregateTimeSlots(timeIds)
	if len(numTime) < 1 {
		return res, errors.New("[iud909]")
	}

	num := d.findMaxMatchingUsers(numTime, len(du))
	//有1个以上有相同时间的用户数, 亦或者, 用户总数为1时, 能匹配结果; 否则匹配失败
	if num > 1 || (duLen == 1 && num > 0) {
		res.Date = d.SimplePeriod(numTime[num])
		res.Res = num == duLen //能匹配最多的用户数==用户数, 证明所有用户都有共同的空闲时间, 匹配成功 !
	}

	err = dao.Tx(func(tx *sqlx.Tx) (e error) {
		return dao.App.DatingDb.DatingUpdate(datingId, &res, tx)
	})
	if err != nil {
		return res, fmt.Errorf("[563fsd]%s", err)
	}

	return
}

func (d *DatingService) aggregateUserTimes(du []db.UserTime) map[int64][]uint {
	timeIds := make(map[int64][]uint)
	maxSecond := (maxTime - minTime) * 3600

	//汇总每个时间段(小时)下, 有哪些人空闲; 如: {"下午13点" : ["张三"]}
	for _, value := range du {
		if len(value.Info) < 1 {
			continue
		}
		info := *value.InfoUnmarshal()
		if info.Time == nil || len(info.Time) < 1 {
			continue
		}

		for _, val := range info.Time {
			if val[0] >= val[1] || int(val[1]-val[0]) > maxSecond {
				continue
			}
			ts := utils.SpreadPeriodToHour(val[0], val[1])
			for _, v := range ts {
				if _, ok := timeIds[v]; !ok {
					timeIds[v] = make([]uint, 0, len(du))
				}
				timeIds[v] = append(timeIds[v], value.Id)
			}

		}
	}
	return timeIds
}
func (d *DatingService) aggregateTimeSlots(timeIds map[int64][]uint) map[int][]int {
	numTime := make(map[int][]int)
	//汇总空闲时间段相同的用户数下, 有哪些时间段(小时); 如: {"有三人空闲" : ["下午13点"]}
	for key, value := range timeIds {
		n := len(value)
		if _, ok := numTime[n]; !ok {
			numTime[n] = make([]int, 0)
		}
		numTime[n] = append(numTime[n], int(key))
	}
	return numTime
}
func (d *DatingService) findMaxMatchingUsers(numTime map[int][]int, duLen int) int {
	for i := duLen; i > 0; i-- {
		if _, ok := numTime[i]; ok && len(numTime[i]) > 0 {
			return i
		}
	}
	return 0
}

// 简化时间表达; 如: ["02-15(8-10时|12-13时)"]
func (d *DatingService) SimplePeriod(unixTimes []int) []string {
	if len(unixTimes) < 1 {
		return []string{}
	}

	dateTimeGroups := utils.UnixGroup(unixTimes)
	dates := make([]string, 0, len(dateTimeGroups))

	for _, val := range dateTimeGroups {

		date := time.Unix(int64(val[0]), 0).In(global.Tz)
		fo := "01-02"
		if time.Now().In(global.Tz).Year() != date.Year() {
			fo = time.DateOnly
		}

		var recommendText bytes.Buffer
		recommendText.WriteString(date.Format(fo))

		//如果是具体时间段才空闲(非一整天都空闲), 需要括号内标出具体时间(小时)
		if l := len(val); l < maxTime-minTime {
			recommendText.WriteString("(")
			recommendText.WriteString(strconv.Itoa(date.Hour()))
			recommendText.WriteString("-")
			for k, v := range val {
				tv := time.Unix(int64(v), 0).In(global.Tz)
				if l == k+1 {
					//如果是最后一个(小时)值, 直接拼接
					recommendText.WriteString(strconv.Itoa(tv.Hour() + 1))
					recommendText.WriteString("时")
				} else if tvlast := time.Unix(int64(val[k+1]), 0).In(global.Tz); tv.Hour()+1 != tvlast.Hour() {
					//判断当前(小时)值和下一个(小时)值是否相同, 避免数据错误
					recommendText.WriteString(strconv.Itoa(tv.Hour() + 1))
					recommendText.WriteString("时")
					recommendText.WriteString("|") //因为不是最后一个(小时)值, 所以还有其他时间段拼接
					recommendText.WriteString(strconv.Itoa(tvlast.Hour()))
					recommendText.WriteString("-")
				}
			}
			recommendText.WriteString(")")
		}

		dates = append(dates, recommendText.String())

	}

	return dates
}
