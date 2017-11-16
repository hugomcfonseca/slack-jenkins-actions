package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"strconv"
	"time"
)

// HTTPResponse Structure used to define response object of every route request
type HTTPResponse struct {
	Status  bool   `json:"status"`
	Content string `json:"content"`
}

// ResponseType Constant
const ResponseType = "application/json"

// ContentType Constant
const ContentType = "Content-Type"

var lag int

func main() {
	
		var portstring string
	
		flag.StringVar(&portstring, "port", "10010", "Listening port")
		flag.Parse()
			
		router := http.NewServeMux()
	
		router.HandleFunc("/job/input", RouteRemoteJenkinsInput)
	
		log.Printf("Listening on port %s ...", portstring)
	
		err2 := http.ListenAndServe(":"+portstring, LogRequests(CheckURL(router)))
		log.Fatal(err2)
}


// LogRequests Middleware level to log API requests
func LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf(
			"[%s]\t%s\t%s",
			r.Method,
			r.URL.String(),
			time.Since(start),
		)
	})
}

// CheckURL Middleware level to validate requested URI
func CheckURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.String()
		pathLength := len(path)
		matchPath := "/role/replica/"
		matchLength := len(matchPath)

		if strings.Contains(path, matchPath) && pathLength > matchLength {
			lag, _ = strconv.Atoi(strings.Trim(path, matchPath))
		} else if strings.Compare(path, strings.TrimRight(path, "/")) != 0 {
			return
		}

		w.Header().Set(ContentType, ResponseType)

		next.ServeHTTP(w, r)
	})
}


// routeResponse Used to build response to API requests
func routeResponse(w http.ResponseWriter, httpStatus bool, contents string) {
	res := new(HTTPResponse)

	if httpStatus {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(403)
	}

	res.Status = httpStatus
	res.Content = contents
	response, _ := json.Marshal(res)
	fmt.Fprintf(w, "%s", response)
}

func RouteRemoteJenkinsInput(w http.ResponseWriter, r *http.Request)  {
	log.Printf("")

	routeResponse(w, true, "")
}