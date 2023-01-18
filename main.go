package main

import (
	// Standard library packages
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	// Third party packages
	"github.com/julienschmidt/httprouter"
	//"github.com/skratchdot/open-golang/open"
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
		fmt.Fprintf(w, "RemoteAddr: %s\n", r.RemoteAddr)
		fmt.Fprintf(w, "X-Forwarded-For: %s\n", r.Header.Get("X-Forwarded-For"))
		fmt.Fprintf(w, "X-Real-IP: %s\n", r.Header.Get("X-Real-IP"))
		fmt.Fprintf(w, "URL: %s\n", r.URL)
		fmt.Fprintf(w, "Host: %s\n", r.Host)
		fmt.Fprintf(w, "RequestURI: %s\n", r.RequestURI)
		fmt.Fprintf(w, "Proto: %s\n", r.Proto)
		fmt.Fprintf(w, "ProtoMajor: %d\n", r.ProtoMajor)
		fmt.Fprintf(w, "ProtoMinor: %d\n", r.ProtoMinor)
		fmt.Fprintf(w, "Header: %s\n", r.Header)
		fmt.Fprintf(w, "Body: %s\n", r.Body)
		fmt.Fprintf(w, "ContentLength: %d\n", r.ContentLength)
		fmt.Fprintf(w, "TransferEncoding: %s\n", r.TransferEncoding)
		fmt.Fprintf(w, "Close: %t\n", r.Close)
	})

	// Add a handler on /test
	r.GET("/test", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Simply write some test data for now
		fmt.Fprint(w, "<p>Welcome!<p>")
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
	fmt.Println("http://0.0.0.0:" + myport + "/ip")

	// Start the blocking server loop.
	log.Fatal(http.Serve(l, r))
}
