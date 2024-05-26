package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/BuddyLongLegs/anginex/src/config"
	"github.com/BuddyLongLegs/anginex/src/logger"
)

type ProxyHandler func(w http.ResponseWriter, r *http.Request)

func addRequestLog(r *http.Response, dbLogger logger.DBLogger) error {
	startTime, err := strconv.ParseInt(r.Request.Header.Get("X-ANG-START-TIME"), 10, 64)
	if err != nil {
		return err
	}
	latency := float64(time.Now().UnixNano()-startTime) / 1000000

	go func() {
		dbLogger.AddLog(
			time.Now().Local().Format(time.RFC3339),
			r.Request.URL.Host,
			r.Request.URL.Path,
			r.Request.Method,
			r.StatusCode,
			latency,
		)
	}()

	return nil
}

func CreateRedirectProxyHandler(route config.RouteConfig, dbLogger logger.DBLogger) ProxyHandler {
	url, err := url.Parse(route.ProxyPass)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.ModifyResponse = func(r *http.Response) error {
		return addRequestLog(r, dbLogger)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now().UnixNano()

		// Set the host and scheme of the request
		r.URL.Host = url.Host
		r.URL.Scheme = url.Scheme
		r.Host = url.Host
		r.Header.Set("X-ANG-START-TIME", fmt.Sprint(startTime))

		forwardPath := r.URL.Path

		// Trim prefix from the request path
		if route.TrimPrefix != "" && strings.HasPrefix(forwardPath, route.TrimPrefix) {
			forwardPath = forwardPath[len(route.TrimPrefix):]

			if len(forwardPath) > 0 && forwardPath[0] != '/' {
				forwardPath = "/" + forwardPath
			}
		}

		// Trim suffix from the request path
		if route.TrimSuffix != "" && strings.HasSuffix(forwardPath, route.TrimSuffix) {
			forwardPath = forwardPath[:len(forwardPath)-len(route.TrimSuffix)]

			if len(forwardPath) > 0 && forwardPath[0] != '/' {
				forwardPath = "/" + forwardPath
			}
		}

		r.URL.Path = forwardPath

		proxy.ServeHTTP(w, r)
	}
}
