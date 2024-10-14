package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twbworld/dating/model/db"
	"github.com/twbworld/dating/service"
)

func TestJWTAuth(t *testing.T) {
	u := db.User{}
	u.Id = 2 //!!!!!!!!!!!!!!!!测试用户
	jwt, err := service.Service.UserServiceGroup.BaseService.LoginToken(&u)
	if err != nil {
		t.Fatal(err)
	}
	userId, newToken, err := JWTAuth(jwt)

	assert.Equal(t, err, nil, newToken, "", userId, u.Id)
}

func TestJWTAuthFail(t *testing.T) {
	userId, newToken, err := JWTAuth("abc")

	assert.NotEqual(t, err, nil)
	assert.Equal(t, newToken, "", userId, uint(0))
}
