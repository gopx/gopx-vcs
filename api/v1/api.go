package v1

import "net/http"

// API handles HTTP requests for API V1.
func API(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch path {
	case "/api/v1/package/new":

	}
}
