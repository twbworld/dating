package user

import (
	"context"
	"database/sql"

	"errors"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram/auth/response"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/twbworld/dating/dao"
	"github.com/twbworld/dating/global"
	"github.com/twbworld/dating/model/common"
	"github.com/twbworld/dating/model/db"
	"github.com/twbworld/dating/utils"
	"golang.org/x/sync/errgroup"
)

type DatingService struct{}

// 用Code向微信换取用户信息
func (d *DatingService) GetUserByCode(code string, u *db.User) (err error) {
	var responseCode *response.ResponseCode2Session
	if gin.Mode() == gin.TestMode {
		responseCode = &response.ResponseCode2Session{OpenID: "abc"}
	} else {
		responseCode, err = utils.AuthWxCode(code)
		if err != nil {
			return
		}
	}

	if err = dao.App.UserDb.GetUserByOpenId(u, responseCode.OpenID); err != nil {
		if err == sql.ErrNoRows {
			//新用户,下一步注册所需信息
			u.OpenId = responseCode.OpenID
			u.UnionId = responseCode.UnionID
			u.SessionKey = responseCode.SessionKey
			return nil
		}
		return
	}
	if u.Id == dao.BaseUserId {
		return errors.New("出现测试账号[rtyoij]")
	}

	if u.SessionKey != responseCode.SessionKey {
		go d.updateSessionKeyAsync(u.Id, responseCode.SessionKey)
	}

	return
}

func (d *DatingService) updateSessionKeyAsync(userID uint, newSessionKey string) {
	defer func() {
		// 避免协程内 panic 影响到外层主程序
		if p := recover(); p != nil {
			global.Log.Error(p)
		}
	}()

	err := dao.Tx(func(tx *sqlx.Tx) error {
		return dao.App.UserDb.UpdateSessionKey(userID, newSessionKey, tx)
	})
	if err != nil {
		panic(err)
	}
}

// 获取会面详情
func (d *DatingService) GetDating(data *common.GetDatingPost, userId uint) (interface{}, error) {
	var (
		datingUsers []common.DatingUser = make([]common.DatingUser, 0)
		dating      db.Dating           = db.Dating{}
		err         error
	)

	if data.Id < 1 {
		return nil, errors.New("参数错误[oihuiu]")
	}

	g, c := errgroup.WithContext(context.Background())
	g.Go(func() error {
		select {
		case <-c.Done(): //发现其他goroutine报错,当前直接退出
			return nil
		default:
			if err := dao.App.DatingDb.GetDating(&dating, data.Id); err != nil {
				return err
			}
		}
		return nil
	})
	g.Go(func() error {
		select {
		case <-c.Done():
			return nil
		default:
			if err = dao.App.DatingDb.GetDatingUsers(&datingUsers, data.Id); err != nil {
				return err
			}
		}
		return nil
	})
	if err = g.Wait(); err != nil {
		return nil, errors.New("参数错误[oi7ja]")
	}

	if userId != 0 {
		isset := false
		for _, value := range datingUsers {
			if value.Id == userId {
				isset = true
				break
			}
		}
		if !isset {
			return nil, errors.New("数据不存在[thhgi]")
		}
	}

	dr := dating.ResultUnmarshal()

	for key, value := range datingUsers {
		if len(value.Info) < 1 {
			continue
		}
		info := *value.InfoUnmarshal()
		if datingUsers[key].InfoResponse == nil {
			datingUsers[key].InfoResponse = make([]common.InfoResponse, 0, len(info.Time))
		}
		for _, val := range info.Time {
			tlist := utils.SpreadPeriodToHour(int(val[0]), int(val[1]))
			res := uint8(0)
			if union := utils.Union(dr.Date, tlist); len(union) > 0 {
				if len(union) == len(tlist) {
					res = 1
				} else {
					res = 2
				}
			}
			t := d.SimplePeriod(tlist)
			datingUsers[key].InfoResponse = append(datingUsers[key].InfoResponse, common.InfoResponse{
				Tag:  t[0],
				Time: [2]string{utils.TimeFormat(val[0]), utils.TimeFormat(val[1])},
				Res:  res,
			})
		}
	}

	return gin.H{
		"dating": gin.H{
			"create_user_id": dating.CreateUserId,
			"id":             dating.Id,
			"status":         dating.Status,
			"result": common.DatingResult{
				Res:  dr.Res,
				Date: d.SimplePeriod(dr.Date),
			},
		},
		"users": datingUsers,
	}, nil
}
