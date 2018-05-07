package route

import (
	"net/http"
	"strings"

	"gopx.io/gopx-vcs/pkg/controller"
	"gopx.io/gopx-vcs/pkg/log"
)

// APIRouter handles API HTTP routes.
type APIRouter struct{}

func (vr APIRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info("%s %s", strings.ToUpper(r.Method), r.RequestURI)
	processAPIRoute(w, r)
}

func processAPIRoute(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = sanitizeAPIPath(r.URL.Path)
	controller.API(w, r)
}

func sanitizeAPIPath(path string) string {
	return strings.ToLower(path)
}

// NewAPIRouter returns a new APIRouter instance.
func NewAPIRouter() *APIRouter {
	return &APIRouter{}
}
