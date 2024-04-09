package analytics

import (
	"log"
	"net/http"
)

func AnalyticsAPI() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", withAuth(func(w http.ResponseWriter, r *http.Request) {
	}))

	log.Default().Println("[analytics]: Server Active")
	http.ListenAndServe(":4000", mux)
}
