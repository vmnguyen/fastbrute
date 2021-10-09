package main

import (
	"bufio"
	"crypto/tls"
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
		//Dial:                     fasthttpproxy.FasthttpHTTPDialer("localhost:8080"),
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}
)
var results []string
var valid []string
var invalid []string
var invalid_tmp []string

func validateResponse(statusCode int, body string) bool {
	// Simple checker
	if statusCode == 404 {
		return false
	}
	return true
}

func doRequest(i interface{}) {
	url := i.(string)
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	req.SetRequestURI(url)
	client.Do(req, resp)
	statusCode := resp.StatusCode()
	body := resp.Body()
	if validateResponse(statusCode, string(body)) {
		valid = append(valid, url)
	} else {
		invalid_tmp = append(invalid_tmp, url)
	}

	results = append(results, fmt.Sprintf("[%d] %s", statusCode, url))
	fmt.Printf("[%d] %s \n", statusCode, url)
}
func scan(target string, concurrent int, path string, level int) {

}
func brute(target string, concurrent int, path string, level int) {
	defer ants.Release()

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

	// Reading wordlist
	var wordlist []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wordlist = append(wordlist, scanner.Text())
	}

	invalid = append(invalid, target)
	for j := 0; j < level; j++ {
		for line := 0; line < len(invalid); line++ {
			for i := 0; i < len(wordlist); i++ {
				wg.Add(1)
				url := invalid[line] + "/" + wordlist[i]
				_ = p.Invoke(string(url))
			}
		}
		wg.Wait()
		invalid = invalid_tmp
		invalid_tmp = nil

	}

	fmt.Println(valid)

}

func main() {
	var target string
	flag.StringVar(&target, "t", "https://example.com", "Target to brute force")

	var concurrent int
	flag.IntVar(&concurrent, "c", 50, "Concurrent")

	var requestPath string
	flag.StringVar(&requestPath, "r", "/path/to/request/file", "Path to request file")

	var wordlist string
	flag.StringVar(&wordlist, "w", "/wordlist/actions.txt", "Path to wordlist file")

	var mode string
	flag.StringVar(&mode, "m", "brute", "Mode to scan - [brute|scan]")

	var level int
	flag.IntVar(&level, "l", 3, "Recursive level")

	var proxy string
	flag.StringVar(&proxy, "x", "http://localhost:8080", "HTTP proxy setting")

	flag.Parse()

	if mode == "brute" {
		brute(target, concurrent, wordlist, level)
	} else {
		scan(target, concurrent, wordlist, level)
	}

}
