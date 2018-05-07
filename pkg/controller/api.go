package controller

import (
	"io"
	"net/http"
	"os"
	"strings"

	"gopx.io/gopx-vcs/api/v1"
	"gopx.io/gopx-vcs/pkg/constants"
	"gopx.io/gopx-vcs/pkg/utils"
)

// API handles api related HTTP requests.
func API(w http.ResponseWriter, r *http.Request) {
	ok := validateAPIRequest(w, r)
	if !ok {
		return
	}

	path := r.URL.Path
	switch {
	case strings.HasPrefix(path, "/api/v1"):
		v1.API(w, r)
	default:
		Error404(w, r)
	}
}

func validateAPIRequest(w http.ResponseWriter, r *http.Request) bool {
	if user, pass, ok := r.BasicAuth(); ok {
		if checkAPIAuth(user, pass) {
			return true
		}
		w.WriteHeader(http.StatusForbidden)
		_, err := io.WriteString(w, "403 Forbidden")
		utils.LogError(err)
		return false
	} else {
		w.Header().Add("WWW-Authenticate", `Basic realm="Access to the GoPX VCS API service", charset="UTF-8"`)
		w.WriteHeader(http.StatusUnauthorized)
		_, err := io.WriteString(w, "401 Unauthorized")
		utils.LogError(err)
		return false
	}
}

func checkAPIAuth(username, password string) bool {
	validUser, isUserSet := os.LookupEnv(constants.ENV_GOPX_VCS_API_AUTH_USER)
	validPass, isPassSet := os.LookupEnv(constants.ENV_GOPX_VCS_API_AUTH_PASSWORD)

	return isUserSet && isPassSet && username == validUser && password == validPass
}
