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
	"github.com/rs/cors"

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
	_, _ = fmt.Fprintf(w, "Your IP is: %s\n", clientIP)
	details := GetDetail(r)
	for _, detail := range details {
		_, _ = fmt.Fprintf(w, "%s\n", detail)
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
			var n net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				n = v.IP
			case *net.IPAddr:
				n = v.IP
			}
			result = append(result, "IP: "+n.String())
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
		if err := tmpl.Execute(w, t); err != nil {
			log.Println("Error executing template")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})))

	r.GET("/detail", loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := ip.GetIP(r)
		details := GetDetail(r)
		t := DetailTemplate{Ip: clientIP, Detail: details}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		tmpl := template.Must(template.ParseFiles("templates/detail.gohtml"))
		if err := tmpl.Execute(w, t); err != nil {
			log.Println("Error executing template")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})))

	r.GET("/json", loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := ip.GetIP(r)
		details := GetDetail(r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, "{\n")
		fmt.Fprintf(w, "  \"ip\": \"%s\",\n", clientIP)
		fmt.Fprintf(w, "  \"detail\": [\n")
		for i, detail := range details {
			fmt.Fprintf(w, "    \"%s\"", detail)
			if i < len(details)-1 {
				fmt.Fprintf(w, ",")
			}
			fmt.Fprintf(w, "\n")
		}
		fmt.Fprintf(w, "  ]\n")
		fmt.Fprintf(w, "}\n")

	})))

	l, err := net.Listen("tcp", "0.0.0.0:"+myport)
	if err != nil {
		log.Fatal(err)
	}

	handler := cors.Default().Handler(r)

	log.Println("Starting the application... http://localhost:" + myport)
	// Start the blocking server loop.
	log.Fatal(http.Serve(l, handler))
}
