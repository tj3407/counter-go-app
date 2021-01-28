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

var results[]string

func callApi(count int, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(fmt.Sprintf("https://postman-echo.com/get?x=%s", strconv.Itoa(count)))
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)
	results = append(results, sb)
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
			rw.Write(data)
		}
	})
	http.ListenAndServe(":8080", nil)
}