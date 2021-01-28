package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Response struct {
	Method string `json:"method"`
	URL string `json:"url"`
	Headers Headers `json:"headers"`
}

type Headers struct {
	XForward string `json:"x-forwarded-proto"`
	XForwardPort string `json:"x-forwarded-port"`
	Host string `json:"host"`
	XAmazonTrace string `json:"x-amzn-trace-id"`
	AcceptEncoding string `json:"accept-encoding"`
	UserAgent string `json:"user-agent"`
	Accept string `json:"accept"`
	CacheControl string `json:"cache-control"`
	PostmanToken string `json:"postman-token"`
}

var results []interface{}

func callApi(count int, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(fmt.Sprintf("https://postman-echo.com/get?x=%s", strconv.Itoa(count)))
	if err != nil {
		log.Fatalln(err)
	}

	d := &Response{}
	d.Method = "GET"
	d.Headers.Accept = "*/*"
	d.Headers.CacheControl = "no-cache"

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	_ = json.Unmarshal(body, &d)
	results = append(results, d)
}

func call_postmanEcho(x string) {
	counter, err := strconv.Atoi(x)
	if err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup

	for i := 0; i < counter; i++ {
		wg.Add(1)
		go callApi(i, &wg)
	}

	wg.Wait()
}

func main() {
	http.HandleFunc("/count", func(rw http.ResponseWriter, r *http.Request) {
		if x := r.FormValue("x"); x != "" {
			call_postmanEcho(x)
			
			// At this point, all the calls to postman echo is done
			data, err := json.Marshal(results)
			if err != nil {
				log.Fatalln(err)
			}
			rw.Header().Add("method", "GET")
			rw.Write(data)
		}
	})
	http.ListenAndServe(":8080", nil)
}