package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"gopx.io/gopx-vcs/pkg/log"
)

// APIResponse represents response for VCS API.
type APIResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// API handles HTTP requests for API v1.
func API(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := strings.ToUpper(r.Method)

	switch path {
	case "/api/v1/package/register":
		switch method {
		case "POST":
			packageRegisterPOST(w, r)
		default:
			Error405(w, r)
		}
	default:
		Error404(w, r)
	}
}

func packageRegisterPOST(w http.ResponseWriter, r *http.Request) {
	mr, err := r.MultipartReader()
	if err != nil {
		log.Error("Error: %s", err)
		writeAPIResponse(w, r, 0, "Error: content type must be type of multipart/form-data")
		return
	}

	var metaField []byte
	var dataField *os.File

	for {
		p, err := mr.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Error("Error: %s", err)
				writeAPIResponse(w, r, 0, "Error: corrupted content received")
				return
			}
		}

		fName := p.FormName()
		switch fName {
		case "meta":
			metaField, err = ioutil.ReadAll(p)
			if err != nil {
				log.Error("Error: %s", err)
				writeAPIResponse(w, r, 0, "Error: corrupted content received")
				return
			}
		case "data":
			dataField, err = ioutil.TempFile("", "package-data-")
			if err != nil {
				Error500(w, r)
				log.Error("Error: %s", err)
				return
			}
			defer os.RemoveAll(dataField.Name())

			// TODO: Here limit the dataField file size
			_, err := io.Copy(dataField, p)
			if err != nil {
				Error500(w, r)
				log.Error("Error: %s", err)
				return
			}
		}

		err = p.Close()
		if err != nil {
			Error500(w, r)
			log.Error("Error: %s", err)
			return
		}
	}

	if metaField == nil {
		writeAPIResponse(w, r, 0, "Error: package meta not provided")
		return
	}

	if dataField == nil {
		writeAPIResponse(w, r, 0, "Error: package data not provided")
		return
	}

	_, err = dataField.Seek(io.SeekStart, 0)
	if err != nil {
		Error500(w, r)
		log.Error("Error: %s", err)
		return
	}

	pkgMeta, err := ParsePackageMeta(string(metaField))
	if err != nil {
		log.Error("Error: %s", err)
		writeAPIResponse(w, r, 0, "Error: package meta should be in valid JSON format")
		return
	}

	err = registerPackage(pkgMeta, dataField)
	if err != nil {
		log.Error("Error: %s", err)
		writeAPIResponse(w, r, 0, fmt.Sprintf("Error: %s", err))
		return
	}

	writeAPIResponse(w, r, 1, "success")
}

func writeAPIResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	apiRes := APIResponse{
		Status:  status,
		Message: message,
	}

	bytes, err := json.Marshal(apiRes)
	if err != nil {
		Error500(w, r)
		log.Error("Error: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(bytes)
	if err != nil {
		log.Error("Error: %s", err)
	}
}
