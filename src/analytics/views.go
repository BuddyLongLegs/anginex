package analytics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func QueryHandler(dbPool *ReadDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vals := r.URL.Query()["q"]
		sqlQueryIndex := strings.TrimPrefix(r.URL.Path, "/api/")
		sqlQuery := PROXY_SQL_QUERIES[sqlQueryIndex]

		for i, v := range vals {
			sqlQuery = ReplaceValsInSQLQuery(sqlQuery, i+1, v)
		}

		rows, err := dbPool.ExecQuery(sqlQuery)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(
				[]byte(fmt.Sprintf("Error: %v", err)),
			)
			return
		}

		defer rows.Close()

		// send the rows to the client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		var results []map[string]interface{}

		for rows.Next() {
			cols := rows.FieldDescriptions()
			row := make(map[string]interface{})
			values, _ := rows.Values()
			for i, col := range cols {
				row[col.Name] = values[i]
			}
			results = append(results, row)
		}

		json, err := json.Marshal(results)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(
				[]byte(fmt.Sprintf("Error: %v", err)),
			)
			return
		}

		w.Write([]byte(json))
	}
}
