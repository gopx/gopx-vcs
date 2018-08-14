package handler

import (
	"net/http"

	"gopx.io/gopx-vcs/pkg/controller/cgi"
)

// CatchAll handles all vcs requests via cgi interfaces.
func CatchAll(w http.ResponseWriter, r *http.Request) {
	cgiHandler := cgi.Handler()
	cgiHandler.ServeHTTP(w, r)
}
