package pgx

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type contextKey string

const transactionKey contextKey = "tx"

func contextWithTransaction(ctx context.Context, transaction pgx.Tx) context.Context {
	return context.WithValue(ctx, transactionKey, transaction)
}

func transactionFromContext(ctx context.Context) (pgx.Tx, bool) {
	transaction, exists := ctx.Value(transactionKey).(pgx.Tx)

	return transaction, exists
}
