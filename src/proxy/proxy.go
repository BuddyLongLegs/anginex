package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	url, _ := url.Parse("http://test-api:8000")
	proxy := httputil.NewSingleHostReverseProxy(url)
	endpoint := "/test"

	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		r.URL.Host = url.Host
		r.URL.Scheme = url.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = url.Host

		path := r.URL.Path
		r.URL.Path = strings.TrimLeft(path, endpoint)

		proxy.ServeHTTP(w, r)
	})

	if err := http.ListenAndServe(":80", mux); err != nil {
		return err
	}

	return nil
}
