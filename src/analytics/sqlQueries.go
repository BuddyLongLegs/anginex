package analytics

var PROXY_SQL_QUERIES = map[string]string{
	"requests_per_unit_time": `
  SELECT 
      time_bucket('$1', time) AS minute,
      COUNT(*) AS call_count
    FROM 
      log_data
    WHERE 
      time >= NOW() - INTERVAL '$2'
    GROUP BY 
      minute
    ORDER BY 
      minute;
  `,
	"avg_and_max_latency_per_unit_time": `
    SELECT 
      time_bucket('$1', time) AS minute,
      AVG(latency) AS avg_latency,
      MAX(latency) AS max_latency
    FROM 
      log_data
    WHERE 
      time >= NOW() - INTERVAL '$2'
    GROUP BY 
      minute
    ORDER BY 
      minute;
  `,
	"system_usage_per_unit_time": `
      SELECT 
        time_bucket('$1', time) AS minute,
        AVG(cpu) AS avg_cpu,
        AVG(memory) AS avg_memory
      FROM 
        system_metrics
      WHERE 
        time >= NOW() - INTERVAL '$2'
      GROUP BY 
        minute
      ORDER BY 
        minute;
    `,
	"most_hit_endpoints": `
    SELECT 
      host,
      path,
      method,
      COUNT(*) AS call_count
    FROM 
      log_data
    WHERE 
      time >= NOW() - INTERVAL '$2'
    GROUP BY 
      host, path, method
    ORDER BY 
      call_count DESC
    LIMIT 10;
  `,
	"most_successful_endpoints_in_time_range": `
    SELECT 
      host,
      path,
      method,
      COUNT(*) AS call_count
    FROM 
      log_data
    WHERE 
      time >= NOW() - INTERVAL '$2'
      AND status_code >= 200
      AND status_code < 300
    GROUP BY 
      host, path, method
    ORDER BY 
      call_count DESC
    LIMIT 10;
  `,
	"most_user_errored_endpoints_in_time_range": `
    SELECT
      host,
      path,
      method,
      COUNT(*) AS call_count
    FROM
      log_data
    WHERE
      time >= NOW() - INTERVAL '$2'
      AND status_code >= 400
      AND status_code < 500
    GROUP BY
      host, path, method
    ORDER BY
      call_count DESC
    LIMIT 10;
  `,
	"most_server_errored_endpoints_in_time_range": `
    SELECT
      host,
      path,
      method,
      COUNT(*) AS call_count
    FROM
      log_data
    WHERE
      time >= NOW() - INTERVAL '$2'
      AND status_code >= 500
    GROUP BY
      host, path, method
    ORDER BY
      call_count DESC
    LIMIT 10;
  `,
}

var SYSTEM_METRICS_SQL_QUERIES = map[string]string{}
