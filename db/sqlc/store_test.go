package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	account1, account2 := createTestAccount(t), createTestAccount(t)
	fmt.Println(">>before: ", account1.Balance, account2.Balance)

	n := 3
	amount := int64(10)

	results := make(chan TransferTxResponse)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxRequest{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			})

			results <- result
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		result := <-results
		err := <-errs

		require.NotEmpty(t, result)
		require.Empty(t, err)

		//check transfer details
		transfer := result.Transfer
		require.Equal(t, account1.ID, transfer.FromAccountID.Int64)
		require.Equal(t, account2.ID, transfer.ToAccountID.Int64)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testQueries.GetTransfer(context.Background(), transfer.ID)
		require.Empty(t, err)

		//check entries
		fromentry := result.FromEntry
		require.NotEmpty(t, fromentry)
		require.Equal(t, account1.ID, fromentry.AccountID.Int64)
		require.Equal(t, -amount, fromentry.Amount)
		require.NotZero(t, fromentry.ID)
		require.NotZero(t, fromentry.CreatedAt)

		_, err = testQueries.GetEntry(context.Background(), fromentry.ID)
		require.Empty(t, err)

		toentry := result.ToEntry
		require.NotEmpty(t, toentry)
		require.Equal(t, account2.ID, toentry.AccountID.Int64)
		require.Equal(t, amount, toentry.Amount)
		require.NotZero(t, toentry.ID)
		require.NotZero(t, toentry.CreatedAt)

		_, err = testQueries.GetEntry(context.Background(), toentry.ID)
		require.Empty(t, err)

		//check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)
		fmt.Println("Tx: >>", account1.Balance, account2.Balance)

		//check for diffs
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)
		require.Equal(t, diff1, diff2)

		k := diff1 / amount
		require.True(t, k >= 1 && k <= int64(n))

	}
	//check the final updated balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NotEmpty(t, updatedAccount1)
	require.Empty(t, err)
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NotEmpty(t, updatedAccount2)
	require.Empty(t, err)

	fmt.Println(">>after: ", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-(amount*int64(n)), updatedAccount1.Balance)
	require.Equal(t, account2.Balance+(amount*int64(n)), updatedAccount2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)
	account1, account2 := createTestAccount(t), createTestAccount(t)
	fmt.Println(">>before: ", account1.Balance, account2.Balance)

	n := 10
	amount := int64(10)

	// results := make(chan TransferTxResponse)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		fromAccountId := account1.ID
		toAccountId := account2.ID
		if i%2 == 1 {
			fromAccountId = account2.ID
			toAccountId = account1.ID
			// account1, account2 = account2, account1
		}
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			_, err := store.TransferTx(ctx, TransferTxRequest{
				FromAccountId: fromAccountId,
				ToAccountId:   toAccountId,
				Amount:        amount,
			})

			// results <- result
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		// result := <-results
		err := <-errs
		require.Empty(t, err)
	}
	//check the final updated balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NotEmpty(t, updatedAccount1)
	require.Empty(t, err)
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NotEmpty(t, updatedAccount2)
	require.Empty(t, err)

	fmt.Println(">>after: ", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)

}
