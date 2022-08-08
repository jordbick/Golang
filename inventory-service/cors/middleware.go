package cors

import "net/http"

// add calls to set the necessary HTTP headers to our requests before the call to serve HTTP

// Middleware
func Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorisation, X-CSRF-Token, Accept-Encoding")
		handler.ServeHTTP(w, r)
	})
}
