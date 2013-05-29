package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var opts struct {
	NumRequests int32 `short:"r" long:"num-requests" description:"Number of requests to make" default:"1"`
}

func main() {
	var err error
	var args []string

	parser := flags.NewParser(&opts, flags.Default)
	if args, err = parser.Parse(); err != nil {
		return
	}

	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Requires one target URL")
		return
	}

	target := args[0]

	var response *http.Response

	startTime := time.Now()

	if response, err = http.Get(target); err != nil {
		panic("Error")
	}

	defer response.Body.Close()
	var body []byte
	if body, err = ioutil.ReadAll(response.Body); err != nil {
		panic("Error")
	}
	fmt.Println(string(body))

	duration := time.Since(startTime)
	fmt.Println(duration)
}
