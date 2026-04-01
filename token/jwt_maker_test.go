package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
	"github.com/surojuvenkatesh20/bank-mgmt/utils"
)

func TestJWTMakerToken(t *testing.T) {
	maker, err := NewJWTMaker(utils.GenerateRandomString(32))
	require.NoError(t, err)

	username := utils.GenerateRandomString(6)
	duration := time.Minute
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

	tokenString, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	payload, err := maker.VerifyToken(tokenString)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiresAt, payload.ExpiresAt, time.Second)
	require.NotZero(t, payload.ID)

}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(utils.GenerateRandomString(32))
	require.NoError(t, err)

	username := utils.GenerateRandomString(6)

	tokenString, err := maker.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	payload, err := maker.VerifyToken(tokenString)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTToken(t *testing.T) {
	payload, err := NewPayload(utils.GenerateRandomString(6), time.Minute)
	require.NoError(t, err)

	//Signing without algorithm
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(utils.GenerateRandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
