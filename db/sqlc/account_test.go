package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/surojuvenkatesh20/bank-mgmt/utils"
)

func createTestAccount(t *testing.T) Account {
	user := createTestUser(t)
	params := CreateAccountParams{
		Owner:    user.Username,
		Currency: utils.GenerateRandomCurrency(),
		Balance:  utils.GenerateRandomMoney(),
	}

	account, err := testQueries.CreateAccount(context.Background(), params)

	//Check if error is nil and account is not empty
	require.NoError(t, err)
	require.NotEmpty(t, account)

	//Check if params values are matching with response
	require.Equal(t, params.Owner, account.Owner)
	require.Equal(t, params.Currency, account.Currency)
	require.Equal(t, params.Balance, account.Balance)

	//Check if id and created_at are not zero of their type, ex: int=0, string="", bool=false
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func deleteTestAccount(accountID int64) error {
	err := testQueries.DeleteAccount(context.Background(), accountID)
	return err
}

func TestCreateAccount(t *testing.T) {
	createTestAccount(t)
}

func TestGetAcount(t *testing.T) {
	account1 := createTestAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	//Check if err is nil and response is not empty
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	//Check if account details in sql response is matching with actual account
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt, time.Second)
	require.Equal(t, account1.Currency, account2.Currency)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createTestAccount(t)

	params := UpdateAccountParams{
		ID:      account1.ID,
		Balance: utils.GenerateRandomMoney(),
	}
	account2, err := testQueries.UpdateAccount(context.Background(), params)

	//Check if err is nil and SQL response is not empty
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	//Check if balance in SQL is matching with updated balance
	require.Equal(t, params.Balance, account2.Balance)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt, time.Second)
	require.Equal(t, account1.Owner, account2.Owner)
}

func TestDeleteAccount(t *testing.T) {
	account := createTestAccount(t)

	//Delete account
	err := deleteTestAccount(account.ID)
	require.NoError(t, err)

	//Get that deleted account.
	account1, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account1)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createTestAccount(t)
	}
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}
