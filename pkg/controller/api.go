package controller

import (
	"net/http"
	"os"

	"gopx.io/gopx-vcs/pkg/constants"
)

// API handles api related HTTP requests.
func API(w http.ResponseWriter, r *http.Request) {

}

func checkAPIAuth(username, password string) bool {
	validUser, isUserSet := os.LookupEnv(constants.ENV_GOPX_VCS_API_AUTH_USER)
	validPass, isPassSet := os.LookupEnv(constants.ENV_GOPX_VCS_API_AUTH_PASSWORD)

	return isUserSet && isPassSet && username == validUser && password == validPass
}
