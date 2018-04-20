package route

import (
	"net/http"
	"strings"

	"gopx.io/gopx-vcs/pkg/controller"
	"gopx.io/gopx-vcs/pkg/log"
)

// VCSRouter handles VCS HTTP routes.
type VCSRouter struct{}

func (vr VCSRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info("%s %s", strings.ToUpper(r.Method), r.RequestURI)
	controller.VCS(w, r)
}

// NewVCSRouter returns a new VCSRouter instance.
func NewVCSRouter() *VCSRouter {
	return &VCSRouter{}
}
