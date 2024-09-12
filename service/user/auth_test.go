package user

import (
    "testing"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/stretchr/testify/assert"
    "github.com/twbworld/dating/global"
    "github.com/twbworld/dating/model/common"
    "github.com/twbworld/dating/model/db"
)

func TestTokenExpirationDuration(t *testing.T) {
	// Setup
    global.Config.JwtKey = "testsecretkey"
    user := &db.User{}
	user.Id = 1
    service := &BaseService{}

    // Generate token
    tokenString, err := service.LoginToken(user)
    assert.NoError(t, err)

    // Parse token
    token, err := jwt.ParseWithClaims(tokenString, &common.JwtInfo{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(global.Config.JwtKey), nil
    })
    assert.NoError(t, err)
    assert.True(t, token.Valid)

    // Extract claims
    claims, ok := token.Claims.(*common.JwtInfo)
    assert.True(t, ok)

    // Validate expiration
    expectedExpiration := time.Now().Add(tokenExpirationDuration).Unix()
    actualExpiration := claims.ExpiresAt.Time.Unix()
    assert.InDelta(t, expectedExpiration, actualExpiration, 1, "Token expiration time should be within 1 second of the expected time")
}
