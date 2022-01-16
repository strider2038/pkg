package pgx

import (
	"context"
	"fmt"

	"github.com/strider2038/pkg/persistence"
)

type Transaction struct {
	ctx context.Context
}

type TransactionManager struct {
	connection Connection
}

func NewTransactionManager(connection Connection) *TransactionManager {
	return &TransactionManager{connection: connection}
}

func (manager *TransactionManager) Begin(ctx context.Context) (persistence.Transaction, error) {
	transaction, err := manager.connection.Scope(ctx).Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &Transaction{ctx: contextWithTransaction(ctx, transaction)}, nil
}

func (manager *TransactionManager) DoTransactionally(
	ctx context.Context,
	action func(ctx context.Context) error,
) error {
	transaction, err := manager.Begin(ctx)
	if err != nil {
		return err
	}

	actionError := action(transaction.Context())
	if actionError != nil {
		err = transaction.Rollback()
		if err != nil {
			return fmt.Errorf("failed to rollback transaction: %w", err)
		}
	} else {
		err = transaction.Commit()
		if err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
	}

	return actionError
}

func (tx *Transaction) Context() context.Context {
	return tx.ctx
}

func (tx *Transaction) Commit() error {
	transaction, _ := transactionFromContext(tx.ctx)

	return transaction.Commit(tx.ctx)
}

func (tx *Transaction) Rollback() error {
	transaction, _ := transactionFromContext(tx.ctx)

	return transaction.Rollback(tx.ctx)
}
