package persistence

import "context"

type TransactionManager interface {
	Begin(ctx context.Context) (Transaction, error)
	DoTransactionally(ctx context.Context, action func(ctx context.Context) error) error
}

type Transaction interface {
	Context() context.Context
	Commit() error
	Rollback() error
}

type NilTransactionManager struct{}

func (NilTransactionManager) Begin(ctx context.Context) (Transaction, error) {
	return NilTransaction{}, nil
}

func (NilTransactionManager) DoTransactionally(ctx context.Context, action func(ctx context.Context) error) error {
	return action(ctx)
}

type NilTransaction struct{}

func (NilTransaction) Context() context.Context {
	return context.Background()
}

func (NilTransaction) Commit() error {
	return nil
}

func (NilTransaction) Rollback() error {
	return nil
}
