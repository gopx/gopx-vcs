package main

import (
	l "log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"gopx.io/gopx-vcs/pkg/config"
	"gopx.io/gopx-vcs/pkg/log"
	"gopx.io/gopx-vcs/pkg/route"
)

var serverLogger = l.New(os.Stdout, "", l.Ldate|l.Ltime|l.Lshortfile)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	startServer()
}

func startServer() {
	addr := addr()
	router := route.NewGoPXVCSAPIRouter()
	server := &http.Server{Addr: addr, Handler: router, ErrorLog: serverLogger}

	log.Info("Running API service on: %s", addr)
	err := server.ListenAndServe()
	log.Fatal("Error: %s", err) // err is always non-nill
}

func addr() string {
	return net.JoinHostPort(config.APIService.Host, strconv.Itoa(config.APIService.Port))
}
