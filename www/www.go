//
// This simple web server responds to requests for /.  In response to requests for /, it returns whatever
// a backend service returns.
//
// This web server is also a little buggy and if the backend return "ERROR" it exits. In that case the web server
// has "crashed".
//
// Furthermore this web server is a little slow and handles at most 10-requests per second -- usually
// a little less.
//
package main

import (
	"net/http"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"sync"
)

var (
	backend string
	backendPort string
	mutex sync.Mutex
)

func hello(w http.ResponseWriter, r *http.Request) {
	//
	// Handle just one request at a time -- hacky, huh.
	//
	mutex.Lock()
	defer mutex.Unlock()

	defer func() {
		if r := recover(); r != nil {
			os.Exit(-1)
		}
	}()
	//
	// simulate a maximum throughput of ten requests per second.
	//
	time.Sleep(100 * time.Millisecond)

	response, err := http.Get(fmt.Sprintf("http://%s:%s", backend, backendPort))
	if err != nil {
		panic(err)
	} else {
		defer response.Body.Close()

		body, _ := ioutil.ReadAll(response.Body)

		if string(body) == "ERROR" {
			panic("www crashed")
		}
		w.Write(body)
	}
}

func main() {
	const (
		defaultPort = "80"
		portUsage   = "Port to listen on"
		defaultBackendPort = "80"
		backendPortUsage   = "Port for backend server"
		defaultBackend = ""
		backendUsage   = "Name of backend server"
	)

	var (
		port string
	)

	flag.StringVar(&port, "port", defaultPort, portUsage)
	flag.StringVar(&backend, "backend", defaultBackend, backendUsage)
	flag.StringVar(&backendPort, "backendPort", defaultBackendPort, backendPortUsage)

	flag.Parse()

	if len(backend) == 0 {
		panic("backend not specified")
	}
	fmt.Println("WWW starting on port " + port + " backend at " + backend + ":" + backendPort)

	http.HandleFunc("/", hello)
	http.ListenAndServe(":" + port, nil)
}