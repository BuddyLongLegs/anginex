package logger

const INIT_DB = `
CREATE TABLE IF NOT EXISTS log_data (
	time TIMESTAMPTZ NOT NULL,
	host TEXT NOT NULL,
	path TEXT NOT NULL,
	method TEXT NOT NULL,
	status_code INT NOT NULL,
	latency FLOAT NOT NULL
);

SELECT create_hypertable('log_data', by_range('time'), if_not_exists => TRUE);
`

const INSERT_LOG = `
	INSERT INTO log_data (time, host, path, method, status_code, latency) VALUES ($1, $2, $3, $4, $5, $6);
`
