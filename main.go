package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Client struct {
	requestUrl       string
	requestInterval  int
	maxCrashDuration int
	caPath           string

	httpClient http.Client
}

func main() {
	os.Exit(realMain())
}

func realMain() int {
	client := &Client{}
	client.parseArgs()

	client.setupCrashRoutine(client.maxCrashDuration)
	client.setupClient(client.caPath)

	for true {
		client.makeRequest(client.requestUrl)
		time.Sleep(time.Duration(client.requestInterval) * time.Second)
	}

	return 0
}

func (c *Client) makeRequest(requestUrl string) {
	resp, err := c.httpClient.Get(requestUrl)
	if err != nil {
		log.Printf("Error requesting [%v]: %v", requestUrl, err)
		return
	}
	if resp == nil || resp.StatusCode != 200 {
		log.Printf("Request to [%v] returned [%v]", requestUrl, resp.StatusCode)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response from [%v]: %v", requestUrl, err)
	}

	log.Printf("Request to [%v] returned [%v]: %v", requestUrl, resp.StatusCode, string(data))
}

func (c *Client) setupClient(pathCa string) {
	if pathCa != "" {
		fileCa, err := ioutil.ReadFile(pathCa)
		if err != nil {
			log.Fatal(err)
		}

		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(fileCa)

		tlsConfig := &tls.Config{RootCAs: certPool}
		tlsConfig.BuildNameToCertificate()
		transport := &http.Transport{TLSClientConfig: tlsConfig}
		c.httpClient = http.Client{Transport: transport}
	} else {
		c.httpClient = *http.DefaultClient
	}
}

func (c *Client) parseArgs() {
	flag.StringVar(&c.requestUrl, "request-url", "", "URL to request")
	flag.IntVar(&c.requestInterval, "request-interval", 1, "interval in seconds to send requests")
	flag.IntVar(&c.maxCrashDuration, "crash", 0, "maximum duration to wait before crashing (default \"0\" - e.g. don't crash on purpose)")
	flag.StringVar(&c.caPath, "ca-path", "", "path to the CA certificate to verify the server's certificate (default \"\")")
	flag.Parse()
}

func (c *Client) setupCrashRoutine(maxCrashDuration int) {
	if c.maxCrashDuration == 0 {
		return
	}

	rand.Seed(time.Now().Unix())
	crashDuration := rand.Intn(c.maxCrashDuration)

	log.Printf("Crashing in [%v] seconds...", crashDuration)
	go func() {
		time.Sleep(time.Duration(crashDuration) * time.Second)
		log.Fatal("Crashing NOW...")
	}()
}
