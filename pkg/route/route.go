package route

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"gopx.io/gopx-common/log"
	"gopx.io/gopx-vcs/pkg/handler"
)

// Router registers the vcs routes.
func Router() *mux.Router {
	r := mux.NewRouter()

	r.Use(loggingMiddleware)

	r.PathPrefix("/").
		HandlerFunc(handler.CatchAll)

	return r
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("%s %s", strings.ToUpper(r.Method), r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
