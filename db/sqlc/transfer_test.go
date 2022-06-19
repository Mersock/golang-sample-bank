package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	amount := int64(10)
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        amount,
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, acc1.ID)
	require.Equal(t, arg.ToAccountID, acc2.ID)
	require.Equal(t, arg.Amount, amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	trans1 := createRandomTransfer(t)
	trans2, err := testQueries.GetTransfer(context.Background(), trans1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, trans2)

	require.Equal(t, trans1.ID, trans2.ID)
	require.Equal(t, trans1.FromAccountID, trans2.FromAccountID)
	require.Equal(t, trans1.ToAccountID, trans2.ToAccountID)
	require.Equal(t, trans1.Amount, trans2.Amount)
	require.WithinDuration(t, trans1.CreatedAt, trans2.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	trans1 := createRandomTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(), trans1.ID)
	require.NoError(t, err)

	trans2, err := testQueries.GetTransfer(context.Background(), trans1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, trans2)
}

func TestListTransfer(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}

	arg := ListTransferParams{
		Limit:  5,
		Offset: 5,
	}

	acc, err := testQueries.ListTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, acc, 5)

	for _, v := range acc {
		require.NotEmpty(t, v)
	}
}
