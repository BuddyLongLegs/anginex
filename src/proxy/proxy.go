package proxy

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/BuddyLongLegs/anginex/src/config"
	"github.com/BuddyLongLegs/anginex/src/logger"
)

func Run(config config.Config, dbLogger *logger.DBLogger) error {
	mux := http.NewServeMux()

	proxyHandlers := make(map[string]ProxyHandler)
	endpoints := make([]string, 0, len(config.Routes))

	for endpoint, route := range config.Routes {
		proxyHandlers[endpoint] = CreateRedirectProxyHandler(route, *dbLogger)
		endpoints = append(endpoints, endpoint)
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, endpoint := range endpoints {
			if strings.HasPrefix(r.URL.Path, endpoint) {
				proxyHandlers[endpoint](w, r)
				return
			}
		}

		http.Error(w, "Not Found", http.StatusNotFound)
		dbLogger.AddLog(
			time.Now().Local().Format(time.RFC3339),
			r.URL.Host,
			r.URL.Path,
			r.Method,
			http.StatusNotFound,
			0,
		)
	})

	log.Default().Println("[anginex]: Server Active")
	if err := http.ListenAndServe(":80", mux); err != nil {
		return err
	}

	return nil
}
