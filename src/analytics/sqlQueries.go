package analytics

var PROXY_SQL_QUERIES = map[string]string{
	"each_query_per_unit_time": `
    SELECT 
      time_bucket('$1', time) AS minute,
      host,
      path,
      method,
      COUNT(*) AS call_count
    FROM 
      log_data
    WHERE 
      time >= NOW() - INTERVAL '$2'
    GROUP BY 
      minute, host, path, method
    ORDER BY 
      minute DESC, host, path, method;
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
      minute DESC;
  `,
}

var SYSTEM_METRICS_SQL_QUERIES = map[string]string{}
