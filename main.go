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
	NumRequests int `short:"r" long:"num-requests" description:"Number of requests to make" default:"1"`
	Concurrent  int `short:"c" long:"concurrent" description:"Number of concurrent connections to make" default:"1"`
}

type result struct {
	duration   time.Duration
	statusCode int
	bytesRead  int
	err        error
}

type Summary struct {
	numRequests          int
	totalRequestDuration time.Duration
	avgRequestDuration   time.Duration
	duration             time.Duration
}

var requestChan chan *http.Request
var resultChan chan *result
var summaryChan chan *Summary
var client *http.Client

func doRequests() {
	for request := range requestChan {
		startTime := time.Now()
		response, err := client.Do(request)
		if err != nil {
			resultChan <- &result{err: err}
			continue

		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			resultChan <- &result{err: err}
			continue
		}

		resultChan <- &result{duration: time.Since(startTime), statusCode: response.StatusCode, bytesRead: len(body)}
	}
}

func generateRequests(target string, numRequests int) {
	request, err := http.NewRequest("GET", target, nil)
	if err != nil {
		panic("Bad target")
	}
	for i := 0; i < numRequests; i++ {
		requestChan <- request
	}
	close(requestChan)
}

func summarizeResults(numRequests int, startTime time.Time) {
	summary := new(Summary)

	for i := 0; i < numRequests; i++ {
		result := <-resultChan
		summary.numRequests++
		summary.totalRequestDuration += result.duration
	}

	summary.duration = time.Since(startTime)
	summary.avgRequestDuration = time.Duration(int64(summary.totalRequestDuration) / int64(summary.numRequests))
	summaryChan <- summary
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

	requestChan = make(chan *http.Request)
	resultChan = make(chan *result)
	summaryChan = make(chan *Summary)
	client = &http.Client{}

	startTime := time.Now()

	for i := 0; i < opts.Concurrent; i++ {
		go doRequests()
	}
	go generateRequests(target, opts.NumRequests)
	go summarizeResults(opts.NumRequests, startTime)

	summary := <-summaryChan

	fmt.Printf("# Requests: %v\n", summary.numRequests)
	fmt.Printf("Duration: %v\n", summary.duration)
	fmt.Printf("Avergage Request Duration: %v\n", summary.avgRequestDuration)
}
