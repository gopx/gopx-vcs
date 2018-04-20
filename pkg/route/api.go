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
	controller.API(w, r)
}

// NewAPIRouter returns a new APIRouter instance.
func NewAPIRouter() *APIRouter {
	return &APIRouter{}
}
