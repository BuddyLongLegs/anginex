package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type addLogFunc func(time string, host string, path string, method string, status_code int, latency float64)

func Run(addLog addLogFunc) error {
	mux := http.NewServeMux()

	url, _ := url.Parse("http://host.docker.internal:8000")
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ModifyResponse = func(r *http.Response) error {
		startTime, err := strconv.ParseInt(r.Request.Header.Get("X-ANG-START-TIME"), 10, 64)
		if err != nil {
			return err
		}
		latency := float64(time.Now().UnixNano()-startTime) / 1000000

		go func() {
			addLog(time.Now().Local().Format(time.RFC3339), r.Request.URL.Host, r.Request.URL.Path, r.Request.Method, r.StatusCode, latency)
		}()

		return nil
	}
	endpoint := "/test"

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, endpoint) {
			// Handle requests starting with endpoint
			startTime := time.Now().UnixNano()
			r.Header.Set("X-ANG-START-TIME", fmt.Sprint(startTime))
			r.URL.Host = url.Host
			r.URL.Scheme = url.Scheme
			r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
			r.Host = url.Host

			// Remove endpoint prefix from the request path
			newPath := strings.TrimPrefix(r.URL.Path, endpoint)
			if newPath == "" || newPath[0] != '/' {
				newPath = "/" + newPath
			}
			r.URL.Path = newPath

			proxy.ServeHTTP(w, r)
		} else {
			// Handle other requests
			http.NotFound(w, r)
		}
	})

	if err := http.ListenAndServe(":80", mux); err != nil {
		return err
	}

	return nil
}
