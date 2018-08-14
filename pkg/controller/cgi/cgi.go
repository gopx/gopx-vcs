package cgi

import (
	"fmt"
	golog "log"
	"net/http/cgi"
	"os"

	"gopx.io/gopx-common/log"
	"gopx.io/gopx-vcs/pkg/config"
	"gopx.io/gopx-vcs/pkg/constants"
)

// cgiErrorWriter provides an error stream for CGI program.
type cgiErrorWriter struct{}

func (cew *cgiErrorWriter) Write(b []byte) (n int, err error) {
	log.Error("CGI: %s", string(b))
	return len(b), nil
}

var cgiLogger = golog.New(os.Stdout, "", golog.Ldate|golog.Ltime|golog.Lshortfile)
var cgiErrorLogger = &cgiErrorWriter{}

// Handler returns the cgi handler for VCS.
func Handler() (h *cgi.Handler) {
	cgiExecPath := config.VCS.CGIPath
	dir := cgiWd()
	env := cgiEnv()
	inheritEnv := cgiInheritEnv()

	h = &cgi.Handler{
		Path:       cgiExecPath,
		Dir:        dir,
		Env:        env,
		InheritEnv: inheritEnv,
		Logger:     cgiLogger,
		Stderr:     cgiErrorLogger,
	}

	return
}

func cgiWd() (wd string) {
	wd, _ = os.Getwd()
	return
}

func cgiEnv() []string {
	return []string{
		fmt.Sprintf("%s=%s", "GIT_PROJECT_ROOT", config.VCS.RepoRoot),
		fmt.Sprintf("%s=%s", "GIT_HTTP_MAX_REQUEST_BUFFER", constants.GitHTTPMaxRequestBuffer),
	}
}

func cgiInheritEnv() []string {
	return []string{
		"PATH",
	}
}
