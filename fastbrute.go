package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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
func brute(request string, concurrent int, path string) {
	defer ants.Release()
	fmt.Println(request)
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(concurrent, func(i interface{}) {
		doRequest(i)
		wg.Done()
	})
	defer p.Release()

	//Read wordlist
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Can't open wordlist")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wg.Add(1)
		url := request + "/" + scanner.Text()
		fmt.Println(url)
		_ = p.Invoke(string(request + ""))
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

	var wordlist string
	flag.StringVar(&wordlist, "w", "/path/to/wordlist/", "Path to wordlist file")

	var mode int
	flag.IntVar(&mode, "m", 1, "Mode to scan, 1 for stdin, 2 for request file")

	flag.Parse()

	brute(target, concurrent, wordlist)
}
