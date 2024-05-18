package analytics

import (
	"log"
	"net/http"

	"github.com/BuddyLongLegs/anginex/src/config"
)

func AnalyticsAPI(config config.Config) {
	// create a db pool with a modified connection string

	dbpool := &ReadDB{}
	dbpool.Conn()
	defer dbpool.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/api/", withAuth(config, QueryHandler(dbpool)))
	mux.HandleFunc("/", withAuth(config, func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/", http.FileServer(http.Dir("./src/analytics/static"))).ServeHTTP(w, r)
	}))

	log.Default().Println("[analytics]: Server Active")
	http.ListenAndServe(":4000", mux)
}
