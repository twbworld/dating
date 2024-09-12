package user

import (
	"database/sql"

	"errors"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram/auth/response"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/twbworld/dating/dao"
	"github.com/twbworld/dating/global"
	"github.com/twbworld/dating/model/db"
	"github.com/twbworld/dating/utils"
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
