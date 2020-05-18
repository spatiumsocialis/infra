package common

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

// Logger wraps an http.Handler and logs requests and their response time
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		// Save a copy of this request for debugging.
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(requestDump))

		log.Printf(
			"request: %s %s %s %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)

	})
}
