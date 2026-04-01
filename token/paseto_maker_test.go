package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/surojuvenkatesh20/bank-mgmt/utils"
)

func TestPasetoMakerToken(t *testing.T) {
	maker, err := NewPasetoMaker(utils.GenerateRandomString(32))
	require.NoError(t, err)

	username := utils.GenerateRandomString(6)
	duration := time.Minute
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiresAt, payload.ExpiresAt, time.Second)
	require.NotZero(t, payload.ID)
}

func TestPasetoMakerExpiredtoken(t *testing.T) {
	maker, err := NewPasetoMaker(utils.GenerateRandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)
	username := utils.GenerateRandomString(6)

	token, err := maker.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)

}


func TestInvalidPasetoToken(t *testing.T) {
		
}