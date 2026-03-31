package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/surojuvenkatesh20/bank-mgmt/utils"
)

func createTestUser(t *testing.T) User {
	hashedPassword, err := utils.HashPassword(utils.GenerateRandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       utils.GenerateRandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       utils.GenerateRandomOwner(),
		Email:          utils.GenerateRandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	//Check if error is nil and user is not empty
	require.NoError(t, err)
	require.NotEmpty(t, user)

	//Check if arg values are matching with response
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	//Check if id and created_at are not zero of their type, ex: int=0, string="", bool=false
	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	createTestUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createTestUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	//Check if err is nil and response is not empty
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	//Check if account details in sql response is matching with actual account
	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.Equal(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	// require.Equal(t, user1.Currency, user2.Currency)
}
