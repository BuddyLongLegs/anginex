package analytics

import (
	"encoding/base64"
	"net/http"
)

type handler func(http.ResponseWriter, *http.Request)

func withAuth(handler handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		requestAuthentication := func() {
			w.Header().Set("WWW-Authenticate", `Basic realm="User Visible Realm", charset="UTF-8"`)
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
		}

		// If auth header not present, then request to authenticate
		if r.Header.Get("Authorization") == "" {
			requestAuthentication()
			return
		}

		creds := r.Header.Get("Authorization")
		creds = creds[len("Basic "):]

		// decode the base64 encoded string
		_, err := base64.StdEncoding.DecodeString(creds)
		if err != nil {
			requestAuthentication()
			return
		}

		handler(w, r)
	}
}
