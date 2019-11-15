package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type statusInfo struct {
	url     string
	content string
	status  int
}

func readFile(filePath string, urlQueue chan string, counter *int) error {
	urls := make(map[string]int) //进行url去重的临时变量
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	s := bufio.NewScanner(f)

	for s.Scan() {
		if strings.Contains(s.Text(), ":0:") || strings.Contains(s.Text(), "::") { //过滤大部分IPv6地址
			continue
		}

		//assembleHTTP := "http://" + s.Text() + "/solr/#/"
		//assembleHTTPs := "https://" + s.Text() + "/solr/#/"

		assembleHTTP := "http://" + s.Text() + ":8081/#/overview"
		assembleHTTPs := "https://" + s.Text() + ":8081/#/overview"

		urls[assembleHTTP] = 1
		urls[assembleHTTPs] = 1

		err = s.Err()
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	for k, _ := range urls {
		urlQueue <- k
		*counter += 1
	}
	fmt.Printf("Total urls to be processed:%d\n", lineCounter)
	close(urlQueue)
	return nil

}
func randomUA() string {
	rand.Seed(time.Now().Unix())
	ua := [] string{"Mozilla/5.0 (Linux; U; Android 4.1.2; zh-cn; GT-I9300 Build/JZO54K) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30 MicroMessenger/5.2.380",
		"MicroMessenger Client",
		"Mozilla/5.0 (Linux; Android 8.1.0; NX606J Build/OPM1.171019.026; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/61.0.3163.98 Mobile Safari/537.36 PARS/5041",
		"Mozilla/5.0 (Linux; Android 8.1.0; PBEM00 Build/OPM1.171019.026; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/62.0.3202.84 Mobile Safari/537.36 PARS/5041",
		"Mozilla/5.0 (Linux; Android 8.0.0; BND-AL10 Build/HONORBND-AL10; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/70.0.3538.110 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 5.1; GN8001L Build/LMY47D; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/47.0.2526.100 Mobile Safari/537.36 PARS/4242 PARSParam (35a7d47d0d845d334e29d3546a863ea79; 02; 4242; 682724741265797120; )",
		"Mozilla/5.0 (Linux; Android 8.1.0; PACT00 Build/O11019; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/62.0.3202.84 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 6.0.1; vivo X9i Build/MMB29M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/51.0.2704.81 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 6.0; DIG-AL00 Build/HUAWEIDIG-AL00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 12_1_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/16D57 PARS/5040;carplugin/0.1.0",
		,}
	return ua[rand.Intn(len(ua))]
}

func randomIP() string {
	return strconv.Itoa(rand.Intn(224)) + "." + strconv.Itoa(rand.Intn(255)) + "." + strconv.Itoa(
		rand.Intn(255)) + "." + strconv.Itoa(rand.Intn(255))
}

func worker1(wg *sync.WaitGroup) {
	client := &http.Client{
		//Timeout: 30 * time.Second,
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(15 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*13)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)

				return c, nil
			},
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	header := map[string]string{
		//"Host":                      "movie.douban.com",
		//"Connection":                "close",
		"Cache-Control":             "max-age=0",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                randomUA(),
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		//"Referer":                   "https://www.baidu.com",
		"X-Forwarded-For":           randomIP(),
	}

	for url := range urls {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		for key, value := range header {
			req.Header.Add(key, value)
		}
		resp, err := client.Do(req)
		if err != nil {
			//fmt.Printf("Connection error in %s\n \tError Details: %s\n", url,err.Error())
			continue
		}
		
		//fmt.Printf("URL: %s | status:%d\n",url,resp.StatusCode
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("IO error: ", err.Error(), "url: ", url)
		}
		content := string(body)
		
		//响应包的判断条件:
		switch strings.Contains(content, "Running Job List") { //Apache Flink
		//switch strings.Contains(content,"groups") { //Apache solr RCE
		//switch strings.Contains(content,"Nginx http upstream check status") {

		case true:
			info := statusInfo{url, content, resp.StatusCode}
			status <- info
		case false:
			info := statusInfo{url, "null", resp.StatusCode}
			status <- info
		}
		resp.Body.Close()

	}
	wg.Done()
}

func workerPool(noOfWorkers int) {
	var wg sync.WaitGroup
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go worker1(&wg)
	}
	wg.Wait()
	close(status)
}

func getResult(done chan bool) {
	for code := range status {
		if code.content != "null" {
			fmt.Printf("URL: %s | Status: %d | Content: null\n", code.url, code.status)
		}
	}
	done <- true
}

var urls = make(chan string, 5000)
var status = make(chan statusInfo, 5000)
var lineCounter int
var targetHitCounter int

func main() {
	startTime := time.Now()
	noDns := "K:\\Desktop\\multitest\\result\\Target_ip_20191114_unique_C.txt"
	go readFile(noDns, urls, &lineCounter)
	noOfWorkers := 3000 //3000并发连接
	done := make(chan bool)
	go getResult(done)
	workerPool(noOfWorkers)
	<-done
	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Printf("Total time taken %.2f seconds. Capability:%d url/s", diff.Seconds(), lineCounter/int(diff.Seconds()))
}
