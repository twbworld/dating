package user

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/twbworld/dating/global"
	"github.com/twbworld/dating/model/common"
	"github.com/twbworld/dating/model/db"
)

const (
	tokenExpirationDuration = time.Hour
	tokenNotBeforeDuration  = -10 * time.Minute
	tokenIssuer             = "dating"
	tokenAudience           = "miniProgram"
)

type BaseService struct{}

func (b *BaseService) LoginToken(user *db.User) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, common.JwtInfo{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpirationDuration)), //过期时间
			NotBefore: jwt.NewNumericDate(time.Now().Add(tokenNotBeforeDuration)),  //生效时间
			Issuer:    tokenIssuer,                                                 //颁发者
			Subject:   strconv.FormatUint(uint64(user.Id), 10),                     //唯一标识
			Audience:  []string{tokenAudience},                                     //平台标识
		},
	}).SignedString([]byte(global.Config.JwtKey))
}
