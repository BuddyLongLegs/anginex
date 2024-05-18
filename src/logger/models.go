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

SELECT add_retention_policy('log_data', INTERVAL '24 hours', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS system_metrics (
	time TIMESTAMPTZ NOT NULL,
	cpu FLOAT NOT NULL,
	memory FLOAT NOT NULL
);

SELECT create_hypertable('system_metrics', by_range('time'), if_not_exists => TRUE);

SELECT add_retention_policy('system_metrics', INTERVAL '24 hours', if_not_exists => TRUE);

-- Create the read only user if it does not exist and grant permissions
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_roles WHERE rolname = 'analytics'
    ) THEN
        CREATE USER analytics WITH PASSWORD 'password';
    END IF;
END $$;

-- Grant SELECT permissions on log_data and system_metrics to the analytics user
GRANT SELECT ON log_data TO analytics;
GRANT SELECT ON system_metrics TO analytics;
`

const INSERT_LOG = `
	INSERT INTO log_data (time, host, path, method, status_code, latency) VALUES ($1, $2, $3, $4, $5, $6);
`

const INSERT_SYSTEM_METRICS = `
	INSERT INTO system_metrics (time, cpu, memory) VALUES ($1, $2, $3);
`
