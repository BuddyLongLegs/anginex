package analytics

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReadDB struct {
	ctx    context.Context
	dbPool *pgxpool.Pool
}

func (d *ReadDB) Conn() {
	ctx := context.Background()
	connStr := os.Getenv("DB_CONN_STRING")
	connStr = ReplaceAuthInfo(connStr)
	dbpool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	d.ctx = ctx
	d.dbPool = dbpool
}

func (d *ReadDB) Close() {
	d.dbPool.Close()
}

func (d *ReadDB) ExecQuery(query string, args ...interface{}) (pgx.Rows, error) {
	rows, err := d.dbPool.Query(d.ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
