package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/panjf2000/ants/v2"
	"github.com/valyala/fasthttp"
)

var (
	client = &fasthttp.Client{
		NoDefaultUserAgentHeader: true,
	}
)

func doRequest(i interface{}) {
	url := i.(string)
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	req.SetRequestURI(url)
	client.Do(req, resp)
	statusCode := resp.StatusCode()
	fmt.Println(statusCode)
}
func brute(request string, concurrent int) {
	defer ants.Release()
	fmt.Println(request)
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(concurrent, func(i interface{}) {
		doRequest(i)
		wg.Done()
	})
	defer p.Release()
	for i := 1; i < 2000; i++ {
		wg.Add(1)
		_ = p.Invoke(string(request))
	}
	wg.Wait()

}

func main() {
	var target string
	flag.StringVar(&target, "t", "https://example.com", "Target to brute force")

	var concurrent int
	flag.IntVar(&concurrent, "c", 50, "Concurrent")

	var requestPath string
	flag.StringVar(&requestPath, "r", "/path/to/request/file", "Path to request file")

	var mode int
	flag.IntVar(&mode, "m", 1, "Mode to scan, 1 for stdin, 2 for request file")

	flag.Parse()

	brute(target, concurrent)
}
