package controller

import (
	"net/http"
)

// Error403 handles forbidden request.
func Error403(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}

// Error404 handles HTTP request on non-existing routes.
func Error404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

// Error405 handles not allowed http method.
func Error405(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// Error500 handles internal server error.
func Error500(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}
