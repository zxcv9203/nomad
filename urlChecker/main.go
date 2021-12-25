package main

import (
	"errors"
	"fmt"
	"net/http"
)

type result struct {
	url    string
	status string
}

var (
	errRequestFailed = errors.New("Request Failed")
)

func main() {
	c := make(chan result)
	urls := []string{
		"https://www.naver.com",
		"https://www.google.com",
		"https://www.reddit.com",
		"https://www.daum.com",
	}
	for _, url := range urls {
		go hitURL(url, c)
	}
	for range urls {
		ret := <-c
		fmt.Println(ret.url, ret.status)
	}
}

func hitURL(url string, c chan<- result) {
	resp, err := http.Get(url)
	status := "OK"
	if err != nil || resp.StatusCode >= 400 {
		status = "FAILED"
	}
	c <- result{url: url, status: status}
}
