package db

import (
	"context"
	"testing"

	"github.com/ahmadfarhanstwn/backend-masterclass/util"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	args := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	createdAccount := createRandomAccount(t)
	fetchedAccount, err := testQueries.GetAccountForUpdate(context.Background(), createdAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, fetchedAccount)

	require.Equal(t, createdAccount.ID, fetchedAccount.ID)
	require.Equal(t, createdAccount.Owner, fetchedAccount.Owner)
	require.Equal(t, createdAccount.Balance, fetchedAccount.Balance)
	require.Equal(t, createdAccount.Currency, fetchedAccount.Currency)
}

func TestUpdateAccount(t *testing.T) {
	createdAccount := createRandomAccount(t)

	updateArg := UpdateAccountParams{
		ID:      createdAccount.ID,
		Balance: util.RandomMoney(),
	}

	updatedAccount, err := testQueries.UpdateAccount(context.Background(), updateArg)

	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)

	require.Equal(t, createdAccount.ID, updatedAccount.ID)
	require.Equal(t, createdAccount.Owner, updatedAccount.Owner)
	require.Equal(t, updateArg.Balance, updatedAccount.Balance)
	require.Equal(t, createdAccount.Currency, updatedAccount.Currency)
}

func TestDeleteAccount(t *testing.T) {
	createdAccount := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), createdAccount.ID)

	require.NoError(t, err)

	deletedAccount, err := testQueries.GetAccountForUpdate(context.Background(), createdAccount.ID)

	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, deletedAccount)
}

func TestListAccount(t *testing.T) {
	var lastAccount Account

	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}
