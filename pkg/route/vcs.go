package route

import (
	"net/http"
	"strings"

	"gopx.io/gopx-vcs/pkg/controller"
	"gopx.io/gopx-vcs/pkg/log"
)

// GoPXVCSRouter handles GoPX VCS HTTP routes.
type GoPXVCSRouter struct{}

func (vr GoPXVCSRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info("%s %s", strings.ToUpper(r.Method), r.RequestURI)
	controller.VCS(w, r)
}

// NewGoPXVCSRouter returns a new GoPXVCSRouter instance.
func NewGoPXVCSRouter() *GoPXVCSRouter {
	return &GoPXVCSRouter{}
}
