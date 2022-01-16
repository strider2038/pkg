package pgx

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Connection interface {
	// Scope returns connection scope or transactional scope from context.
	Scope(ctx context.Context) Scope
	// Ping checks connection to a database.
	Ping(ctx context.Context) error
	// Close closes underlying connection.
	Close()
}

type Scope interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row
}

type Pool struct {
	pool *pgxpool.Pool
}

func NewPool(pool *pgxpool.Pool) *Pool {
	return &Pool{pool: pool}
}

func (p *Pool) Scope(ctx context.Context) Scope {
	transaction, ok := transactionFromContext(ctx)
	if ok {
		return transaction
	}

	return p.pool
}

func (p *Pool) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

func (p *Pool) Close() {
	p.pool.Close()
}

type Conn struct {
	conn *pgx.Conn
}

func NewConn(conn *pgx.Conn) *Conn {
	return &Conn{conn: conn}
}

func (c *Conn) Get() *pgx.Conn {
	return c.conn
}

func (c *Conn) Scope(ctx context.Context) Scope {
	transaction, ok := transactionFromContext(ctx)
	if ok {
		return transaction
	}

	return c.conn
}

func (c *Conn) Ping(ctx context.Context) error {
	return c.conn.Ping(ctx)
}

func (c *Conn) Close() {
	c.conn.Close(context.Background())
}
