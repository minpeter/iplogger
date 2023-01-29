package main

import (
	"errors"
	"fmt"
	"ipLogger/utils"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

var env utils.Env

func getIP(r *http.Request) (net.IP, error) {

	forwards := strings.Split(r.Header.Get("X-Forwarded-For"), ",")

	fmt.Println("forwards: ", forwards)
	fmt.Println("forwards len: ", len(forwards))
	if forwards[0] == "" {
		remoteAddr, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return nil, err
		}
		return net.ParseIP(remoteAddr), nil
	}
	if env.TrustMode == 1 {
		for i := 0; i < env.ProxyHopCount; i++ {
			if i >= len(forwards) {
				break
			}
			ip := net.ParseIP(forwards[i])
			if ip != nil {
				return ip, nil
			}
		}
	}

	return nil, errors.New("no ip found")

}

func main() {
	var err error
	env, err = utils.SetEnv()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("is proxy: %t\n", env.IsProxy)
	fmt.Printf("trust mode: %d\n", env.TrustMode)
	fmt.Printf("proxy hop count: %d\n", env.ProxyHopCount)
	fmt.Printf("trust address: %s\n\n", env.TrustAddress)

	myport := strconv.Itoa(10000)

	r := httprouter.New()

	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		ip, err := getIP(r)
		if err != nil {
			fmt.Fprintf(w, "Error: %s\n", err)
		}
		fmt.Fprintf(w, "Your IP is: %s\n", ip)
	})

	l, err := net.Listen("tcp", "0.0.0.0:"+myport)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("http://localhost:" + myport)
	// Start the blocking server loop.
	log.Fatal(http.Serve(l, r))
}
