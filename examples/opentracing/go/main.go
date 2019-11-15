package main

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var envMark = os.Getenv("envMark")
var url = os.Getenv("url")

func printOpenTracingText(w http.ResponseWriter, r *http.Request) {
	tracer, closer := jaeger.NewTracer("demo", jaeger.NewConstSampler(false), jaeger.NewNullReporter())
	defer closer.Close()
	ctx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	var span opentracing.Span
	if err != nil {
		span = tracer.StartSpan("demo")
	} else {
		span = tracer.StartSpan("demo", opentracing.ChildOf(ctx))
	}

	var reqEnvMark, requestText string = span.BaggageItem("ali-env-mark"), ""
	if reqEnvMark == "" {
		reqEnvMark = "empty"
	}

	hdr := opentracing.HTTPHeadersCarrier{}
	err = tracer.Inject(span.Context(), opentracing.HTTPHeaders, hdr)

	if url != "" && url != "none" {

		httpReq, _ := http.NewRequest("GET", url, nil)
		err = hdr.ForeachKey(func(key, val string) error {
			httpReq.Header.Add(key, val)
			return nil
		})
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
