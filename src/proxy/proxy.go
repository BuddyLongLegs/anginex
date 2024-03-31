package proxy

import (
	"net/http"
)

func Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	if err := http.ListenAndServe(":80", mux); err != nil {
		return err
	}

	return nil
}
