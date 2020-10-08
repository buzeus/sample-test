package main

import (
	"crypto/tls"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
)

const (
	Stamp      = "Jan _2 15:04:05"
	StampMilli = "Jan _2 15:04:05.000"
	StampMicro = "Jan _2 15:04:05.000000"
	StampNano  = "Jan _2 15:04:05.000000000"
)

func main() {
	http.HandleFunc("/check/", checkConnectivity)
	http.HandleFunc("/delay/", delayedResponse)

	http.Handle("/swagger.yml", http.FileServer(http.Dir("./swagger/")))

	opts := middleware.RedocOpts{SpecURL: "/swagger.yml"}
	sh := middleware.Redoc(opts, nil)
	http.Handle("/docs", sh)

	http.HandleFunc("/favicon.ico", doNothing)

	http.HandleFunc("/", home)

	fmt.Println("Starting the server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("failed to start")
	}
}

func doNothing(w http.ResponseWriter, r *http.Request) {}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello from homepage")
	w.Write([]byte("Hello from homepage"))
}

// Sample Request - http://localhost:8080/delay/5
func delayedResponse(w http.ResponseWriter, r *http.Request) {
	arg := r.URL.Path[7:] // Excluding "delay/"
	// Delay will be zero if no argument is passed or invalid argument is passed. This means no delay
	delay, _ := strconv.Atoi(arg)

	recTime := time.Now().Format(StampMicro)
	time.Sleep(time.Duration(delay) * time.Second)
	respTime := time.Now().Format(StampMicro)
	retMsg := fmt.Sprintf("Request Received at %v \nWaited for %d seconds \nResponse returned at %v", recTime, delay, respTime)
	fmt.Println(retMsg)
	w.Write([]byte(retMsg))
}

// Sample Request - http://localhost:8080/check/https,google.com
func checkConnectivity(w http.ResponseWriter, r *http.Request) {
	args := r.URL.Path[7:]
	argsSplit := strings.Split(args, ",")

	if argsSplit != nil && len(argsSplit) > 0 && (argsSplit[0] == "https" || argsSplit[0] == "http") {
		protoParam := argsSplit[0]
		urlParam := argsSplit[1]
		fmt.Println(protoParam)
		fmt.Println(urlParam)

		urlToBeHit := fmt.Sprintf("%s://%s", protoParam, urlParam)

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client := &http.Client{Transport: tr}
		resp, err := client.Get(urlToBeHit)
		if err != nil {
			fmt.Println(err)
		}
		check(err)
		fmt.Println(w)
		resp.Write(w)
	} else {
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(requestDump)
		w.Write(requestDump)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
