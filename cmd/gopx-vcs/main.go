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
	switch {
	case config.VCSService.UseHTTP && config.VCSService.UseHTTPS:
		go startHTTP()
		startHTTPS()
	case config.VCSService.UseHTTP:
		startHTTP()
	case config.VCSService.UseHTTPS:
		startHTTPS()
	default:
		log.Fatal("Error: no listener is specified in VCS service config file")
	}
}

func startHTTPS() {
	addr := httpsAddr()
	router := route.NewVCSRouter()
	server := &http.Server{Addr: addr, Handler: router, ErrorLog: serverLogger}

	log.Info("Running HTTPS server on: %s", addr)
	err := server.ListenAndServeTLS(config.VCSService.CertFile, config.VCSService.KeyFile)
	log.Fatal("Error: %s", err) // err is always non-nill
}

func startHTTP() {
	addr := httpAddr()
	router := route.NewVCSRouter()
	server := &http.Server{Addr: addr, Handler: router, ErrorLog: serverLogger}

	log.Info("Running HTTP server on: %s", addr)
	err := server.ListenAndServe()
	log.Fatal("Error: %s", err) // err is always non-nill
}

func httpAddr() string {
	return net.JoinHostPort(config.VCSService.Host, strconv.Itoa(config.VCSService.HTTPPort))
}

func httpsAddr() string {
	return net.JoinHostPort(config.VCSService.Host, strconv.Itoa(config.VCSService.HTTPSPort))
}
