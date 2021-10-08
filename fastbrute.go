package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
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

func FasthttpHTTPDialer(proxy string) fasthttp.DialFunc {
	return func(addr string) (net.Conn, error) {
		var auth string

		if strings.Contains(proxy, "@") {
			split := strings.Split(proxy, "@")
			auth = base64.StdEncoding.EncodeToString([]byte(split[0]))
			proxy = split[1]

		}

		conn, err := fasthttp.Dial(proxy)
		if err != nil {
			return nil, err
		}

		req := "CONNECT " + addr + " HTTP/1.1\r\n"
		if auth != "" {
			req += "Proxy-Authorization: Basic " + auth + "\r\n"
		}
		req += "\r\n"

		if _, err := conn.Write([]byte(req)); err != nil {
			return nil, err
		}

		res := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(res)

		res.SkipBody = true

		if err := res.Read(bufio.NewReader(conn)); err != nil {
			conn.Close()
			return nil, err
		}
		if res.Header.StatusCode() != 200 {
			conn.Close()
			return nil, fmt.Errorf("could not connect to proxy")
		}
		return conn, nil
	}
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
	fmt.Printf("[%d] %s \n", statusCode, url)
}
func brute(target string, concurrent int, path string) {
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

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wg.Add(1)
		url := target + "/" + scanner.Text()
		_ = p.Invoke(string(url))
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

	var proxy string
	flag.StringVar(&proxy, "x", "http://localhost:8080", "HTTP proxy setting")

	flag.Parse()

	brute(target, concurrent, wordlist)
}
