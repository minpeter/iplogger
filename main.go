package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	ipextractor "ipLogger/pkg/IPExtractor"
)

func main() {

	var err error
	if err != nil {
		log.Fatal(err)
	}

	myport := strconv.Itoa(10000)

	r := httprouter.New()

	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		ip := ipextractor.ExtractIPFromXFF()
		fmt.Fprintf(w, "implemented by myself Get IP is: %s\n", ip(r))

		ip = ipextractor.ExtractIPFromXFFHeader()
		fmt.Fprintf(w, "implemented by echo framework Get IP is: %s\n", ip(r))

		dip := ipextractor.ExtractIPDirect()
		fmt.Fprintf(w, "\n\ndirectly get IP is: %s\n", dip(r))
		fmt.Fprintf(w, "XFF: %s", r.Header.Get("X-Forwarded-For"))

	})

	l, err := net.Listen("tcp", "0.0.0.0:"+myport)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("http://localhost:" + myport)
	// Start the blocking server loop.
	log.Fatal(http.Serve(l, r))
}
