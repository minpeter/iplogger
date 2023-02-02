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

		ip = ipextractor.ExtractIPFromXFFHeader(
			ipextractor.TrustIPRange(&net.IPNet{IP: net.IPv4(173, 245, 48, 0), Mask: net.CIDRMask(20, 32)}),
			ipextractor.TrustIPRange(&net.IPNet{IP: net.IPv4(103, 21, 244, 0), Mask: net.CIDRMask(22, 32)}),
			ipextractor.TrustIPRange(&net.IPNet{IP: net.IPv4(103, 22, 200, 0), Mask: net.CIDRMask(22, 32)}),
			ipextractor.TrustIPRange(&net.IPNet{IP: net.IPv4(103, 31, 4, 0), Mask: net.CIDRMask(22, 32)}),
			ipextractor.TrustIPRange(&net.IPNet{IP: net.IPv4(141, 101, 64, 0), Mask: net.CIDRMask(18, 32)}),
			ipextractor.TrustIPRange(&net.IPNet{IP: net.IPv4(108, 162, 192, 0), Mask: net.CIDRMask(18, 32)}),
			ipextractor.TrustIPRange(&net.IPNet{IP: net.IPv4(190, 93, 240, 0), Mask: net.CIDRMask(20, 32)}),
			ipextractor.TrustIPRange(&net.IPNet{IP: net.IPv4(188, 114, 96, 0), Mask: net.CIDRMask(20, 32)}),
			ipextractor.TrustIPRange(&net.IPNet{IP: net.IPv4(197, 234, 240, 0), Mask: net.CIDRMask(22, 32)}),
			ipextractor.TrustIPRange(&net.IPNet{IP: net.IPv4(198, 41, 128, 0), Mask: net.CIDRMask(17, 32)}),
			ipextractor.TrustIPRange(&net.IPNet{IP: net.IPv4(162, 158, 0, 0), Mask: net.CIDRMask(15, 32)}),
			ipextractor.TrustIPRange(&net.IPNet{IP: net.IPv4(104, 16, 0, 0), Mask: net.CIDRMask(13, 32)}),
			ipextractor.TrustIPRange(&net.IPNet{IP: net.IPv4(104, 24, 0, 0), Mask: net.CIDRMask(14, 32)}),
			ipextractor.TrustIPRange(&net.IPNet{IP: net.IPv4(172, 64, 0, 0), Mask: net.CIDRMask(13, 32)}),
			ipextractor.TrustIPRange(&net.IPNet{IP: net.IPv4(131, 0, 72, 0), Mask: net.CIDRMask(22, 32)}),
			ipextractor.TrustIPRange(&net.IPNet{IP: net.ParseIP("2400:cb00::"), Mask: net.CIDRMask(32, 128)}),
			ipextractor.TrustIPRange(&net.IPNet{IP: net.ParseIP("2606:4700::"), Mask: net.CIDRMask(32, 128)}),
			ipextractor.TrustIPRange(&net.IPNet{IP: net.ParseIP("2803:f800::"), Mask: net.CIDRMask(32, 128)}),
			ipextractor.TrustIPRange(&net.IPNet{IP: net.ParseIP("2405:b500::"), Mask: net.CIDRMask(32, 128)}),
			ipextractor.TrustIPRange(&net.IPNet{IP: net.ParseIP("2405:8100::"), Mask: net.CIDRMask(32, 128)}),
			ipextractor.TrustIPRange(&net.IPNet{IP: net.ParseIP("2a06:98c0::"), Mask: net.CIDRMask(29, 128)}),
			ipextractor.TrustIPRange(&net.IPNet{IP: net.ParseIP("2c0f:f248::"), Mask: net.CIDRMask(32, 128)}),
		)
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
