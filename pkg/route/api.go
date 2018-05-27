package route

import (
	"net/http"
	"path"
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
	r.URL.Path = sanitizeAPIRoute(r.URL.Path)
	controller.API(w, r)
}

// Here requested route needs to be converted to lower case,
// which enables "/ap1/V1" is equivalent to "/ap1/v1" etc.
// and finally cleans the path e.g. end slashes would be removed from path
// e.g. "/api/" -> "/api" etc.
func sanitizeAPIRoute(route string) string {
	return path.Clean(strings.ToLower(route))
}

// NewAPIRouter returns a new APIRouter instance.
func NewAPIRouter() *APIRouter {
	return &APIRouter{}
}
