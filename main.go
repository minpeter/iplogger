package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	ipextractor "ipLogger/pkg/IPExtractor"
	echo "ipLogger/pkg/echo"
)

func main() {

	var err error
	if err != nil {
		log.Fatal(err)
	}

	myport := strconv.Itoa(10000)

	r := httprouter.New()

	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprintf(w, "implemented by myself Get IP is: %s\n", ipextractor.IP(r))

		e := echo.ExtractIPFromXFFHeader(
			echo.TrustIPRange(&net.IPNet{Mask: net.CIDRMask(20, 32), IP: net.ParseIP("173.245.48.0")}),
			echo.TrustIPRange(&net.IPNet{Mask: net.CIDRMask(22, 32), IP: net.ParseIP("103.21.244.0")}),
			echo.TrustIPRange(&net.IPNet{Mask: net.CIDRMask(22, 32), IP: net.ParseIP("103.22.200.0")}),
			echo.TrustIPRange(&net.IPNet{Mask: net.CIDRMask(22, 32), IP: net.ParseIP("103.31.4.0")}),
			echo.TrustIPRange(&net.IPNet{Mask: net.CIDRMask(18, 32), IP: net.ParseIP("141.101.64.0")}),
			echo.TrustIPRange(&net.IPNet{Mask: net.CIDRMask(18, 32), IP: net.ParseIP("108.162.192.0")}),
			echo.TrustIPRange(&net.IPNet{Mask: net.CIDRMask(20, 32), IP: net.ParseIP("190.93.240.0")}),
			echo.TrustIPRange(&net.IPNet{Mask: net.CIDRMask(20, 32), IP: net.ParseIP("188.114.96.0")}),
			echo.TrustIPRange(&net.IPNet{Mask: net.CIDRMask(22, 32), IP: net.ParseIP("197.234.240.0")}),
			echo.TrustIPRange(&net.IPNet{Mask: net.CIDRMask(17, 32), IP: net.ParseIP("198.41.128.0")}),
			echo.TrustIPRange(&net.IPNet{Mask: net.CIDRMask(15, 32), IP: net.ParseIP("162.158.0.0")}),
			echo.TrustIPRange(&net.IPNet{Mask: net.CIDRMask(13, 32), IP: net.ParseIP("104.16.0.0")}),
			echo.TrustIPRange(&net.IPNet{Mask: net.CIDRMask(14, 32), IP: net.ParseIP("104.24.0.0")}),
			echo.TrustIPRange(&net.IPNet{Mask: net.CIDRMask(13, 32), IP: net.ParseIP("172.64.0.0")}),
			echo.TrustIPRange(&net.IPNet{Mask: net.CIDRMask(22, 32), IP: net.ParseIP("131.0.72.0")}),
			echo.TrustIPRange(&net.IPNet{Mask: net.CIDRMask(32, 128), IP: net.ParseIP("2400:cb00::")}),
			echo.TrustIPRange(&net.IPNet{Mask: net.CIDRMask(32, 128), IP: net.ParseIP("2606:4700::")}),
			echo.TrustIPRange(&net.IPNet{Mask: net.CIDRMask(32, 128), IP: net.ParseIP("2803:f800::")}),
			echo.TrustIPRange(&net.IPNet{Mask: net.CIDRMask(32, 128), IP: net.ParseIP("2405:b500::")}),
			echo.TrustIPRange(&net.IPNet{Mask: net.CIDRMask(32, 128), IP: net.ParseIP("2405:8100::")}),
			echo.TrustIPRange(&net.IPNet{Mask: net.CIDRMask(29, 128), IP: net.ParseIP("2a06:98c0::")}),
			echo.TrustIPRange(&net.IPNet{Mask: net.CIDRMask(32, 128), IP: net.ParseIP("2c0f:f248::")}),
		)
		fmt.Fprintf(w, "implemented by echo framework Get IP is: %s\n", e(r))

		ra, _, _ := net.SplitHostPort(r.RemoteAddr)
		fmt.Fprintf(w, "\n\ndirectly get IP is: %s\n", ra)
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
