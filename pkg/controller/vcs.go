package controller

import (
	"net/http"
)

// VCS handles vcs related HTTP requests e.g. package cloning.
func VCS(w http.ResponseWriter, r *http.Request) {
	h := cgiHandler()
	h.ServeHTTP(w, r)
}
