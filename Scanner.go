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
		"Mozilla/5.0 (Linux; Android 8.1.0; PADM00 Build/O11019; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/62.0.3202.84 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 8.1.0; vivo X20Plus A Build/OPM1.171019.011; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/62.0.3202.84 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 8.0.0; DUK-AL20 Build/HUAWEIDUK-AL20; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/68.0.3440.91 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 12_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 PARS/5000",
		"Mozilla/5.0 (Linux; Android 6.0; vivo Y67 Build/MRA58K; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/51.0.2704.81 Mobile Safari/537.36 PARS/5020",
		"Mozilla/5.0 (Linux; Android 9; HMA-AL00 Build/HUAWEIHMA-AL00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/71.0.3578.99 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 7.0; 20170829D Build/NRD90M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/62.0.3202.84 Mobile Safari/537.36 PARS/4242 PARSParam (3c31ffd5d7e4b03d0c481b456fd61ebd2; 02; 4242; 1049239852270952448; )",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 11_1 like Mac OS X) AppleWebKit/604.3.5 (KHTML, like Gecko) Mobile/15B93 PARS/5040",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 11_2 like Mac OS X) AppleWebKit/604.4.7 (KHTML, like Gecko) Mobile/15C114 PARS/5040 hybridwebview cordovashouxian pajkslaAppVersion/3058",
		"Mozilla/5.0 (Linux; Android 8.0.0; HWI-AL00 Build/HUAWEIHWI-AL00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/69.0.3497.100 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 8.1.0; PAAM00 Build/OPM1.171019.011; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/62.0.3202.84 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 6.0.1; OPPO R9s Build/MMB29M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36 PARS/5041",
		"Mozilla/5.0 (Linux; Android 7.1.1; OPPO R11s Build/NMF26X; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 8.1.0; Redmi 5 Plus Build/OPM1.171019.019; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/70.0.3538.110 Mobile Safari/537.36 PARS/5023",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Mobile/13B143 PARS/5040",
		"Mozilla/5.0 (Linux; Android 8.1.0; NX606J Build/OPM1.171019.026; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/61.0.3163.98 Mobile Safari/537.36 PARS/5041",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 11_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15G77 PARS/5040",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 11_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15F79 PARS/5040",
		"Mozilla/5.0 (Linux; Android 6.0; vivo Y67 Build/MRA58K; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/51.0.2704.81 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 8.1.0; JKM-TL00 Build/HUAWEIJKM-TL00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/70.0.3538.110 Mobile Safari/537.36 PARS/5030",
		"Mozilla/5.0 (Linux; Android 8.1.0; PBEM00 Build/OPM1.171019.026; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/62.0.3202.84 Mobile Safari/537.36 PARS/5041",
		"Mozilla/5.0 (Linux; Android 6.0.1; NX549J Build/MMB29M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 8.0.0; BND-AL10 Build/HONORBND-AL10; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/70.0.3538.110 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 12_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 PARS/5040 hybridwebview cordovashouxian pajkslaAppVersion/3058",
		"Mozilla/5.0 (Linux; Android 7.1.1; OPPO R11s Plus Build/NMF26X; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 5.1; GN8001L Build/LMY47D; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/47.0.2526.100 Mobile Safari/537.36 PARS/4242 PARSParam (35a7d47d0d845d334e29d3546a863ea79; 02; 4242; 682724741265797120; )",
		"Mozilla/5.0 (Linux; Android 8.1.0; PACT00 Build/O11019; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/62.0.3202.84 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 6.0.1; vivo X9i Build/MMB29M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/51.0.2704.81 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 8.1.0; OE106 Build/OPM1.171019.026; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/62.0.3202.84 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 8.1.0; PBEM00 Build/OPM1.171019.026; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/62.0.3202.84 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 8.1.0; V1809A Build/OPM1.171019.026; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/62.0.3202.84 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 7.1.2; vivo X9 Build/N2G47H; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36 PARS/5023",
		"Mozilla/5.0 (Linux; Android 5.1.1; OPPO R9 Plusm A Build/LMY47V; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 9; CLT-AL01 Build/HUAWEICLT-AL01; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/71.0.3578.99 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 6.0; DIG-AL00 Build/HUAWEIDIG-AL00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 8.1.0; vivo X21A Build/OPM1.171019.011; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/62.0.3202.84 Mobile Safari/537.36 PARS/5041",
		"Mozilla/5.0 (Linux; Android 9; MI 8 Build/PKQ1.180729.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/71.0.3578.99 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 7.0; MHA-AL00 Build/HUAWEIMHA-AL00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 12_1_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/16D57 PARS/5040;carplugin/0.1.0",
		"Mozilla/5.0 (Linux; Android 8.1.0; PADM00 Build/O11019; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/62.0.3202.84 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 5.1; vivo X6Plus D Build/LMY47I) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/39.0.0.0 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 8.1.0; vivo X20Plus A Build/OPM1.171019.011; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/62.0.3202.84 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 11_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E302 PARS/5040 hybridwebview cordovashouxian pajkslaAppVersion/3058",
		"Mozilla/5.0 (Linux; Android 8.0.0; DUK-AL20 Build/HUAWEIDUK-AL20; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/68.0.3440.91 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 12_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 PARS/5010",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 12_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 PARS/5000",
		"Mozilla/5.0 (Linux; Android 6.0.1; OPPO R9s Build/MMB29M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 5.1; OPPO R9tm Build/LMY47I; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/43.0.2357.121 Mobile Safari/537.36 PARS/5041",
		"Mozilla/5.0 (Linux; Android 6.0; vivo Y67 Build/MRA58K; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/51.0.2704.81 Mobile Safari/537.36 PARS/5020",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 12_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 PARS/5010 hybridwebview cordovashouxian pajkslaAppVersion/3057",
		"Mozilla/5.0 (Linux; Android 9; HMA-AL00 Build/HUAWEIHMA-AL00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/71.0.3578.99 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 7.0; 20170829D Build/NRD90M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/62.0.3202.84 Mobile Safari/537.36 PARS/4242",
		"Mozilla/5.0 (Linux; Android 8.0.0; HWI-AL00 Build/HUAWEIHWI-AL00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/63.0.3239.111 Mobile Safari/537.36 PARS/5023",
		"Mozilla/5.0 (Linux; Android 8.1.0; Mi Note 3 Build/OPM1.171019.019; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/68.0.3440.91 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 11_1 like Mac OS X) AppleWebKit/604.3.5 (KHTML, like Gecko) Mobile/15B93 PARS/5040",
		"Mozilla/5.0 (Linux; Android 6.0.1; OPPO R9s Plus Build/MMB29M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36 PARS/5040",
		"Mozilla/5.0 (Linux; Android 9; MI 8 Lite Build/PKQ1.181007.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/71.0.3578.99 Mobile Safari/537.36 PARS/5041",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 11_2 like Mac OS X) AppleWebKit/604.4.7 (KHTML, like Gecko) Mobile/15C114 PARS/5040 hybridwebview cordovashouxian pajkslaAppVersion/3058",}
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
			//DialContext: (&net.Dialer{
			//	Timeout:   30 * time.Second,
			//	Deadline:time.Now().Add(10 * time.Second),
			//}).DialContext, //会导致提前结束

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

		//Transport: &http.Transport{
		//	DialContext: (&net.Dialer{
		//		Timeout:   8 * time.Second,
		//		Deadline:time.Now().Add(10 * time.Second),
		//	}).DialContext,
		//	//MaxIdleConns:          200,
		//	//IdleConnTimeout:       10 * time.Second,
		//	//TLSHandshakeTimeout:   10 * time.Second,
		//	//ExpectContinueTimeout: 10 * time.Second,
		//	//TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		//},

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
