package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/minpeter/iplogger/pkg/ip"
	"github.com/minpeter/iplogger/pkg/useragent"
)

type IndexTemplate struct {
	Ip string
}

type DetailTemplate struct {
	Detail []string
	Ip     string
}

func openLogFile() *os.File {
	// log folder must be created before running the app
	if _, err := os.Stat("log"); os.IsNotExist(err) {
		os.Mkdir("log", 0755)
	}
	f, err := os.OpenFile("log/access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func commandLineResponse(w http.ResponseWriter, r *http.Request) {
	clientIP := ip.GetIP(r)
	fmt.Fprintf(w, "Your IP is: %s\n", clientIP)
	ra, _, _ := net.SplitHostPort(r.RemoteAddr)
	fmt.Fprintf(w, "\n\ndirectly get IP is: %s\n", ra)
	if r.Header.Get("X-Forwarded-For") != "" {
		fmt.Fprintf(w, "XFF: %s\n", r.Header.Get("X-Forwarded-For"))
	}
}

func GetDetail(r *http.Request) []string {

	var result []string

	u, _ := url.Parse(r.URL.String())
	wait := u.Query().Get("wait")
	if len(wait) > 0 {
		duration, err := time.ParseDuration(wait)
		if err == nil {
			time.Sleep(duration)
		}
	}

	hostname, _ := os.Hostname()
	result = append(result, "Hostname: "+hostname)

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			result = append(result, "IP: "+ip.String())
		}
	}

	result = append(result, "RemoteAddr: "+r.RemoteAddr)

	// method + path + HTTP version
	result = append(result, r.Method+" "+r.URL.Path+" "+r.Proto)

	headerStartIdx := len(result)

	// headers
	for name, headers := range r.Header {
		for _, h := range headers {
			result = append(result, name+": "+h)
		}
	}

	sort.Strings(result[headerStartIdx:])

	return result
}

func main() {
	logFile := openLogFile()
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))

	myport := strconv.Itoa(10000)

	r := httprouter.New()

	r.GET("/", loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if useragent.IsCommandLine(r.UserAgent()) {
			commandLineResponse(w, r)
			return
		}
		clientIP := ip.GetIP(r)
		t := IndexTemplate{Ip: clientIP}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		tmpl := template.Must(template.ParseFiles("templates/index.gohtml"))
		tmpl.Execute(w, t)
	})))

	r.GET("/detail", loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := ip.GetIP(r)
		details := GetDetail(r)
		t := DetailTemplate{Ip: clientIP, Detail: details}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		tmpl := template.Must(template.ParseFiles("templates/detail.gohtml"))
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
