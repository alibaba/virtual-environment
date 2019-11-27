package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var envMark = os.Getenv("envMark")
var url = os.Getenv("url")

func printOpenTracingText(w http.ResponseWriter, r *http.Request) {
	reqEnvMark := r.Header.Get("ali-env-mark")
	if reqEnvMark == "" {
		reqEnvMark = "empty"
	}
	var requestText = ""
	if url != "" && url != "none" {
		httpReq, _ := http.NewRequest("GET", url, nil)
		httpReq.Header = r.Header
		resp, err := http.DefaultClient.Do(httpReq)
		if err != nil {
			requestText = "call " + url + " failed"
		} else {
			defer resp.Body.Close()
			if err != nil {
				requestText = "call " + url + " failed"
			} else {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					requestText = "call " + url + " failed"
				} else {
					requestText = string(body)
				}
			}
		}
		requestText += "\n"
	}

	fmt.Fprintf(w, requestText+"[go @ "+envMark+"] <-"+reqEnvMark+"\n")
}

func main() {
	log.Printf("envMark:" + envMark)
	log.Printf("url:" + url)

	http.HandleFunc("/demo", printOpenTracingText)
	log.Printf("listening to port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
