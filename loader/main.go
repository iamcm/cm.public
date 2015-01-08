package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"time"
)

var (
	urls          []string
	numworkers    int
	numiterations int
)

type UrlResult struct {
	Response http.Response
	Duration int64
	Err      error
}

func loadurl(c chan UrlResult, theUrl string) {
	r := UrlResult{}
	start := time.Now().Unix()

	client := &http.Client{
		Timeout: time.Second * 3,
	}

	resp, err := client.Get(theUrl)

	end := time.Now().Unix()
	elapsed := end - start

	if err == nil {
		r.Response = *resp
	}
	r.Duration = elapsed
	r.Err = err

	c <- r
}

func login(c chan UrlResult, theUrl string) {
	r := UrlResult{}
	start := time.Now().Unix()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout: time.Second * 3,
	}
	client.Transport = tr
	cj, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client.Jar = cj

	resp, err := client.Get("")
	dump, _ := httputil.DumpResponse(resp, true)
	fmt.Println(string(dump))
	resp, err = client.PostForm("", url.Values{"_username": {""}, "_password": {""}})

	end := time.Now().Unix()
	elapsed := end - start

	if err == nil {
		r.Response = *resp
	}
	r.Duration = elapsed
	r.Err = err

	c <- r
}

func main() {
	urls := make([]string, 0)
	urls = append(urls, "http://surf.iamcm.co.uk")

	numworkers = 5
	numiterations = 3

	for i := 0; i < numiterations; i++ {
		c := make(chan UrlResult)

		for j := 0; j < numworkers; j++ {
			for _, url := range urls {
				go loadurl(c, url)
			}
		}

		var totalTime int64 = 0
		errors := 0
		for j := 0; j < numworkers*len(urls); j++ {
			res := <-c
			if res.Err != nil {
				errors += 1
				fmt.Println(res.Err)
			} else {
				fmt.Printf("%s: %d seconds", res.Response.Status, res.Duration)
				fmt.Println("")
				totalTime = totalTime + res.Duration

				/*dump, _ := httputil.DumpResponse(&res.Response, true)
				fmt.Println(string(dump))*/
			}
		}

		successes := (numworkers * len(urls)) - errors

		fmt.Printf("Average response time: %d", totalTime/int64(successes))
		fmt.Println("")
		fmt.Printf("Errors: %d", errors)
		fmt.Println("")

		time.Sleep(time.Second * 1)
	}
}
