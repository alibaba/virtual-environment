package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:\n" + os.Args[0] + " status\n" + os.Args[0] + " version")
		os.Exit(0)
	}
	resp, err := http.Get("http://127.0.0.1:8000/" + os.Args[1])
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
