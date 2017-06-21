//
// This simple backend web server responds to requests for /.  In response to requests for /, it generally returns
// the time but sometimes it gets an error and returns "ERROR".
//
package main

import (
	"io"
	"flag"
	"net/http"
	"fmt"
	"time"
	"os"
	"math/rand"
)

var fo *os.File;
const successRate = 0.98	// Requests are 98% successful.

func logme(w http.ResponseWriter, r *http.Request) {
	t := time.Now().Format(time.RFC850)

	response := t

	if rand.Float64() >= successRate {
		response = "ERROR"
	}
	// Return an error every N responses
	io.WriteString(w, response)

	// log all responses to the log file
	_, err := fo.WriteString(fmt.Sprintf("%s,%s,%s\n", t, r.RemoteAddr, response))
	if err != nil {
		panic(err)
	}
}

func main() {
	const (
		defaultPort = "80"
		portUsage   = "Port to listen/serve on"
		defaultLogFile = "/dev/console"
		fileUsage   = "Log file name"
	)

	var (
		port string
		file string
		err error
	)

	flag.StringVar(&port, "port", defaultPort, portUsage)
	flag.StringVar(&file, "log", defaultLogFile, fileUsage)

	flag.Parse()

	// open output file
	fo, err = os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	fmt.Println("Backend starting on port " + port + " logging to " + file)

	rand.Seed(time.Now().Unix())

	http.HandleFunc("/", logme)
	http.ListenAndServe(":" + port, nil)
}