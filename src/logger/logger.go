package logger

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBLogger struct {
	ctx    context.Context
	dbPool *pgxpool.Pool
}

func (d *DBLogger) Conn() {
	ctx := context.Background()
	connStr := os.Getenv("DB_CONN_STRING")
	dbpool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	d.ctx = ctx
	d.dbPool = dbpool
}

func (d *DBLogger) Close() {
	d.dbPool.Close()
}

func (d *DBLogger) InitDatabase() {
	_, err := d.dbPool.Exec(d.ctx, INIT_DB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create the `log_data` hypertable: %v\n", err)
		os.Exit(1)
	}
}

func (d *DBLogger) AddLog(time string, host string, path string, method string, status_code int, latency float64) {
	_, err := d.dbPool.Exec(d.ctx, INSERT_LOG, time, host, path, method, status_code, latency)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert log data: %v\n", err)
		os.Exit(1)
	}
}

func (d *DBLogger) ExecQuery(query string, args ...interface{}) (pgx.Rows, error) {
	rows, err := d.dbPool.Query(d.ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (d *DBLogger) AddSystemMetrics(time string, cpu float64, memory float64) {
	if cpu == 0 || memory == 0 {
		return
	}

	_, err := d.dbPool.Exec(d.ctx, INSERT_SYSTEM_METRICS, time, cpu, memory)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert system metrics: %v\n", err)
		os.Exit(1)
	}
}
