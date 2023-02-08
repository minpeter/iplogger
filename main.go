package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/minpeter/iplogger/pkg/ip"
	"github.com/minpeter/iplogger/pkg/useragent"
)

type IpTemplate struct {
	Ip string
}

func openLogFile() *os.File {
	f, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

// logging middleware
func loggingMiddleware(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// log.Println(r.Method, r.URL.Path, "IP: "+ip.GetIP(r))
		log.Printf("%s %s %s %s", r.Method, r.URL.Path, "IP: "+ip.GetIP(r), r.UserAgent())
		next.ServeHTTP(w, r)
	}
}

func main() {
	logFile := openLogFile()
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))

	myport := strconv.Itoa(10000)

	r := httprouter.New()

	r.GET("/", loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := ip.GetIP(r)

		if useragent.IsCommandLine(r.UserAgent()) {
			fmt.Fprintf(w, "Your IP is: %s\n", clientIP)
			ra, _, _ := net.SplitHostPort(r.RemoteAddr)
			fmt.Fprintf(w, "\n\ndirectly get IP is: %s\n", ra)
			if r.Header.Get("X-Forwarded-For") != "" {
				fmt.Fprintf(w, "XFF: %s\n", r.Header.Get("X-Forwarded-For"))
			}
			return
		}

		t := IpTemplate{Ip: clientIP}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		//return clientIP to html
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, t)
	})))

	l, err := net.Listen("tcp", "0.0.0.0:"+myport)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting the application... http://localhost:" + myport)
	// Start the blocking server loop.
	log.Fatal(http.Serve(l, r))
}
