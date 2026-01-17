package exercises

import (
	"fmt"
	"net/http"
)

//we must giver our error paths the same attention we give our algorithms.

// In this function we are only printing if there are an error from the response
func CheckStatus() {

	checkStatus := func(
		done <-chan struct{},
		urls ...string,
	) <-chan *http.Response {

		responses := make(chan *http.Response)

		go func() {
			defer close(responses)
			for _, url := range urls {
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println(err)
					continue
				}
				select {
				case <-done:
					return
				case responses <- resp:
				}

			}

		}()
		return responses
	}

	done := make(chan struct{})
	defer close(done)

	urls := []string{"https://www.google.com", "https://badhost"}

	for response := range checkStatus(done, urls...) {
		fmt.Printf("response: %v\n", response.Status)
	}

}

type Result struct {
	Error    error
	Response *http.Response
}

func CheckStatusHandleError() {

	checkStatus := func(
		done <-chan struct{},
		urls ...string,
	) <-chan Result {

		responses := make(chan Result)

		go func() {
			defer close(responses)
			for _, url := range urls {
				resp, err := http.Get(url)
				result := Result{err, resp}
				select {
				case <-done:
					return
				case responses <- result:
				}

			}

		}()
		return responses
	}

	done := make(chan struct{})
	defer close(done)

	urls := []string{"https://www.google.com", "https://badhost"}
	//here we separate the concerns of error handling from our producer goroutine
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v\n", result.Error)
			continue
		}

		fmt.Printf("response: %v\n", result.Response.Status)
	}

}
