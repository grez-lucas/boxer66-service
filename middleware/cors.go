package middleware

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
)

var (
	originAllowList = []string{"http://localhost:9000"}
	methodAllowList = []string{"GET", "POST", "DELETE", "OPTIONS", "PATCH"}
	allowedHeaders  = []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"}
)

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPreflight(r) {
			fmt.Println("Detected pre-flight request")
			origin := r.Header.Get("Origin")
			method := r.Header.Get("Access-Control-Request-Method")
			if slices.Contains(originAllowList, origin) && slices.Contains(methodAllowList, method) {
				// Preflight request (OPTIONS)
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(methodAllowList, ", "))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))
				w.WriteHeader(http.StatusNoContent)
				return
			}
		} else {
			// Not pre-flight
			origin := r.Header.Get("Origin")
			if slices.Contains(originAllowList, origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
		}
		w.Header().Add("Vary", "Origin")
		next.ServeHTTP(w, r)
	})
}

func isPreflight(r *http.Request) bool {
	return r.Method == "OPTIONS" &&
		r.Header.Get("Origin") != "" &&
		r.Header.Get("Access-Control-Request-Method") != ""
}
