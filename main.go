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

func main() {
	os.Exit(realMain())
}

func realMain() int {
	requestUrl := flag.String("request-url", "", "URL to request")
	requestInterval := flag.Int("request-interval", 1, "interval in seconds to send requests")
	maxCrashDuration := flag.Int("crash", 0, "maximum duration to wait before crashing")
	caPath := flag.String("ca-path", "", "path to the CA certificate")
	flag.Parse()

	if *maxCrashDuration != 0 {
		setupCrashRoutine(*maxCrashDuration)
	}
	client := getClient(*caPath)

	for true {
		makeRequest(client, *requestUrl)
		time.Sleep(time.Duration(*requestInterval) * time.Second)
	}

	return 0
}

func makeRequest(client http.Client, requestUrl string) {
	resp, err := client.Get(requestUrl)
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

func getClient(pathCa string) http.Client {
	var client http.Client
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
		client = http.Client{Transport: transport}
	} else {
		client = *http.DefaultClient
	}
	return client
}

func setupCrashRoutine(maxCrashDuration int) {
	rand.Seed(time.Now().Unix())
	crashDuration := rand.Intn(maxCrashDuration)

	log.Printf("Crashing in [%v] seconds", crashDuration)
	go func() {
		time.Sleep(time.Duration(crashDuration) * time.Second)
		log.Fatal("Crashing...")
	}()
}
