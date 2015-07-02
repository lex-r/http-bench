package main

import (
	"io/ioutil"
	"log"
	"fmt"
	"strings"
	"net/http"
	"math/rand"
	"time"
	"runtime"
	"net/url"
	"flag"
)

var concurrency = flag.Int("c", 100, "Concurrency")
var linksFile = flag.String("f", "./links.txt", "File with urls")

var links []string

func init() {
	flag.Parse()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(int64(time.Now().Second()))
	file, err := ioutil.ReadFile(*linksFile)
	if err != nil {
		log.Fatal("Error while reading file. ", err)
	}

	links = strings.Split(string(file), "\n")
	fmt.Printf("Links: \n%v\n", len(links))

	channels := make([]chan bool, *concurrency)

	for i := 0; i < len(channels); i++ {
		channels[i] = make(chan bool)
		go run(10000, channels[i])
	}

	for i := 0; i < len(channels); i++ {
		<- channels[i]
	}
}

func run(requests int, c chan bool) {
	maxFails := 100
	fails := 0
	maxLinks := len(links)

	client := &http.Client{}
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		log.Printf("Error while creating request. %v", err)
		c <- false
	}

	for i := 0; i < requests; i++ {
		log.Printf("Iter: %v", i)
		link := links[rand.Intn(maxLinks)]

		l, _ := url.Parse(link)
		req.URL = l
		req.Host = l.Host
		if err != nil {
			log.Printf("Error while creating request. %v", err)
		}
		start := time.Now()
		resp, err := client.Do(req)
		t := time.Since(start)

		if err != nil {
			fails++
			log.Printf("Error response.", err)
			if fails > maxFails {
				c <- true
				break
			}

			continue
		}

		log.Printf("Time: %v", t.Seconds())
		log.Printf("Response: %v", resp.Status)
	}

	c <- true
}
