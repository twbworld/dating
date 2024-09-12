package middleware

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/twbworld/dating/dao"
	"github.com/twbworld/dating/global"
	"github.com/twbworld/dating/model/common"
	"github.com/twbworld/dating/model/db"
	"github.com/twbworld/dating/service"
	"github.com/twbworld/dating/utils"
)

// 为了重写Response的body
// 自定义一个结构体，实现 gin.ResponseWriter interface
type reWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// 重写 Write([]byte) (int, error) 方法
func (w reWriter) Write(b []byte) (int, error) {
	//向一个bytes.buffer中写一份数据来为获取body使用
	// w.ResponseWriter.Write(b)
	return w.body.Write(b)
}

// jwt授权验证
func JWTAuth(ctx *gin.Context) {
	authCode := ctx.GetHeader(`Authorization`)
	if authCode == "" {
		common.FailAuth(ctx, `非法请求`)
		return
	}

	token, err := jwt.ParseWithClaims(strings.TrimSpace(authCode), &common.JwtInfo{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.JwtKey), nil
	})

	if err != nil {
		//token过期
		common.FailAuth(ctx, `未登录[ojiryu]`)
		return
	}

	claims, ok := token.Claims.(*common.JwtInfo)
	if !ok || !token.Valid {
		common.FailAuth(ctx, `未登录[oftyu]`)
		return
	}

	et, err := claims.GetExpirationTime()
	now := time.Now()
	if err != nil || now.After(et.Time) {
		common.FailAuth(ctx, `登录过期[ojikyu]`)
		return
	}

	if audience, err := claims.GetAudience(); err != nil || utils.InSlice(audience, `miniProgram`) < 0 {
		common.FailAuth(ctx, `未登录[ojidu]`)
		return
	}

	userIdStr, err := claims.GetSubject()
	if err != nil {
		common.FailAuth(ctx, `未登录[ujhgio]`)
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		common.FailAuth(ctx, `未登录[ujdsgio]`)
		return
	}

	//jwt验证成功, 给业务层设置user_id
	ctx.Set(`userId`, uint(userId))

	var (
		rw       *reWriter
		newToken string
	)
	nu, eu := now.Unix(), et.Unix()
	if nb, err := claims.GetNotBefore(); (err != nil && nu+60*10 > eu) || (err == nil && 2*nu > eu+nb.Unix()) {
		// (eu - nb.Unix()) / 2 > eu - nu
		//token临近过期, 准备更新
		//也可以考虑使用双token机制(不过已经有小程序的code验证机制, 就不做复杂了)
		u := &db.User{}
		if dao.App.UserDb.GetUserById(u, uint(userId)) != sql.ErrNoRows && u.Id > 0 {
			newToken, err = service.Service.UserServiceGroup.BaseService.LoginToken(u)
			if err == nil && len(newToken) > 0 {
				rw = &reWriter{
					ctx.Writer,
					&bytes.Buffer{},
				}
				ctx.Writer = rw
			}
		}

	}

	ctx.Next() //执行业务代码

	if rw != nil {
		data, bo := common.Response{}, rw.body.Bytes()
		rw.body.Reset()

		if err := json.Unmarshal(bo, &data); err != nil {
			rw.ResponseWriter.Write(bo) //用回原来的数据
			global.Log.Errorln(err, `[dgihu]`)
			return
		}

		data.Token = newToken //给与新token
		dataJson, err := json.Marshal(&data)
		if err != nil {
			rw.ResponseWriter.Write(bo)
			global.Log.Errorln(err, `[dofigj9]`)
			return
		}
		rw.ResponseWriter.Write(dataJson)
	}
}
