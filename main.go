package main

import (
	// Standard library packages
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	// Third party packages
	"github.com/julienschmidt/httprouter"
)

// https://blog.golang.org/context/userip/userip.go
func getIP(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ip, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		//return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)

		fmt.Fprintf(w, "userip: %q is not IP:port", r.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		//return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
		fmt.Fprintf(w, "userip: %q is not IP:port", r.RemoteAddr)
		return
	}

	// This will only be defined when site is accessed via non-anonymous proxy
	// and takes precedence over RemoteAddr
	// Header.Get is case-insensitive
	forward := r.Header.Get("X-Forwarded-For")

	fmt.Fprintf(w, "<p>IP: %s</p>", ip)
	fmt.Fprintf(w, "<p>Port: %s</p>", port)
	fmt.Fprintf(w, "<p>Forwarded for: %s</p>", forward)

	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<p><a href='/'>back to home</a></p>")
}

func main() {
	myport := strconv.Itoa(10000)

	// Instantiate a new router
	r := httprouter.New()

	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprintf(w, "<h1>ipLogger</h1>")
		fmt.Fprintf(w, "<p>get ip address: <a href='./ip'>%s/ip</a></p>", r.Host)
		fmt.Fprintf(w, "<p>whoami: <a href='./whoami'>%s/whoami</a></p>", r.Host)
		fmt.Fprintf(w, "<p>test: <a href='./test'>%s/test</a></p>", r.Host)
	})

	r.GET("/ip", getIP)

	r.GET("/whoami", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		hostname, err := os.Hostname()
		if err != nil {
			fmt.Fprintf(w, "Hostname err: %s\n", err)
		}
		fmt.Fprintf(w, "Hostname: %s\n", hostname)
		fmt.Fprintf(w, "RemoteAddr: %s\n", r.RemoteAddr)
		fmt.Fprintf(w, "URL: %s\n", r.URL)
		fmt.Fprintf(w, "Host: %s\n", r.Host)
		fmt.Fprintf(w, "RequestURI: %s\n", r.RequestURI)
		fmt.Fprintf(w, "Proto: %s\n", r.Proto)
		fmt.Fprintf(w, "Header: ...\n")
		for k, v := range r.Header {
			if k == "Cookie" || k == "Sec-Ch-Ua" {
				fmt.Fprintf(w, "\t%s: ...\n", k)
				for _, vv := range strings.Split(strings.Trim(v[0], " "), ";") {
					fmt.Fprintf(w, "\t\t%s\n", vv)
				}
			} else {
				fmt.Fprintf(w, "\t%s: %s \n", k, v)
			}
		}
	})

	// Add a handler on /test
	r.GET("/test", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		forward := r.Header.Get("X-Forwarded-For")
		var ip string
		var err error
		if forward != "" {
			ip = strings.Split(forward, ",")[0]
		} else {
			ip, _, err = net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, "Error parsing remote address ["+r.RemoteAddr+"]", http.StatusInternalServerError)
				return
			}
		}
		fmt.Fprintf(w, "<p>IP: %s</p>", ip)
		fmt.Fprintf(w, "The IP above always refers to your IP.")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<p><a href='/'>back to home</a></p>")
	})

	l, err := net.Listen("tcp", "0.0.0.0:"+myport)
	if err != nil {
		log.Fatal(err)
	}
	// The browser can connect now because the listening socket is open.

	//err = open.Start("http://localhost:"+ myport + "/test")
	//err = open.Start("http://localhost:"+ myport + "/ip")
	//if err != nil {
	//     log.Println(err)
	//}
	fmt.Println("http://localhost:" + myport + "/ip")

	// Start the blocking server loop.
	log.Fatal(http.Serve(l, r))
}
