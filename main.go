package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"

	ipextractor "ipLogger/pkg/IPExtractor"
)

func main() {
	logFile := openLogFile()
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))

	myport := strconv.Itoa(10000)

	r := httprouter.New()

	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		clientIP := ipextractor.IP(r)
		fmt.Fprintf(w, "Your IP is: %s\n", clientIP)
		ra, _, _ := net.SplitHostPort(r.RemoteAddr)
		fmt.Fprintf(w, "\n\ndirectly get IP is: %s\n", ra)
		fmt.Fprintf(w, "XFF: %s", r.Header.Get("X-Forwarded-For"))

		log.Println("IP: " + clientIP + " - " + r.Header.Get("X-Forwarded-For") + " - " + ra)

	})

	l, err := net.Listen("tcp", "0.0.0.0:"+myport)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting the application... http://localhost:" + myport)
	// Start the blocking server loop.
	log.Fatal(http.Serve(l, r))
}

func openLogFile() *os.File {
	f, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return f
}
