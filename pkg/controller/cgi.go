package controller

import (
	"fmt"
	l "log"
	"net/http/cgi"
	"os"

	"gopx.io/gopx-vcs/pkg/config"
	"gopx.io/gopx-vcs/pkg/log"
)

// CGIErrorWriter provides an error stream for CGI program.
type CGIErrorWriter struct{}

func (cw *CGIErrorWriter) Write(b []byte) (n int, err error) {
	log.Info("CGI: %s", b)
	return len(b), nil
}

var cgiLogger = l.New(os.Stdout, "", l.Ldate|l.Ltime|l.Lshortfile)
var cgiErrorLogger = &CGIErrorWriter{}

func cgiHandler() *cgi.Handler {
	cgiExecPath := config.VCS.CGIPath
	dir := cgiWd()
	env := cgiEnv()
	inheritEnv := cgiInheritEnv()

	return &cgi.Handler{
		Path:       cgiExecPath,
		Dir:        dir,
		Env:        env,
		InheritEnv: inheritEnv,
		Logger:     cgiLogger,
		Stderr:     cgiErrorLogger,
	}
}

func cgiWd() string {
	wd, err := os.Getwd()
	if err != nil {
		// Use the base directory of the CGI executable as fallback
		wd = ""
	}
	return wd
}

func cgiEnv() []string {
	return []string{
		fmt.Sprintf("%s=%s", "GIT_PROJECT_ROOT", config.VCS.RepoRoot),
		fmt.Sprintf("%s=%s", "GIT_HTTP_EXPORT_ALL", ""),
		fmt.Sprintf("%s=%s", "GIT_HTTP_MAX_REQUEST_BUFFER", "10M"),
	}
}

func cgiInheritEnv() []string {
	return []string{
		"PATH",
	}
}
