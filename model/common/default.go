package common

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/twbworld/dating/model/db"
)

type DatingUser struct {
	UtId         uint   `db:"ut_id" json:"ut_id"`
	Info         string `db:"info" json:"-"`
	Id           uint   `db:"id" json:"id"`
	NickName     string `db:"nick_name" json:"nick_name"`
	Path         string `db:"path" json:"avatar_url"`
	InfoResponse struct {
		InfoPost
		TimeStr []string `json:"ts"`
	} `json:"info"`
}

type DatingUserJoin struct {
	UserId       uint `db:"user_id" json:"user_id"`
	DatingId     uint `db:"dating_id" json:"dating_id"`
	CreateUserId uint `db:"create_user_id" json:"create_user_id"`
}

type DatingList struct {
	Id           uint     `db:"id" json:"id"`
	Status       int8     `db:"status" json:"status"`
	AddTime      int64    `db:"add_time" json:"-"`
	AddTimeStr   string   `json:"add_time"`
	UtId         uint     `db:"ut_id" json:"ut_id"`
	CreateUserId uint     `db:"create_user_id" json:"create_user_id"`
	AvatarUrl    []string `json:"avatar_url"`
}

type DatingUserAvatar struct {
	DatingId uint   `db:"dating_id" json:"dating_id"`
	Path     string `db:"path" json:"path"`
}

type JwtInfo struct {
	jwt.RegisteredClaims
	//可自定义数据; 加入IP地址/用户名等非敏感数据
}

// json转结构体
func (d *DatingUser) InfoUnmarshal() *db.UserTimeInfo {
	return (&db.UserTime{
		Info: d.Info,
	}).InfoUnmarshal()
}
