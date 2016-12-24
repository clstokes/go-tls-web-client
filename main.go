package main

import (
  "crypto/tls"
  "crypto/x509"
  "io/ioutil"
  "log"
  "net/http"
  "os"
)

func main() {
  serviceAddr := os.Getenv("SERVICE_ADDR")
  if serviceAddr == "" {
    log.Fatal("SERVICE_ADDR must be set and non-empty")
  }

  pathCa := os.Getenv("PATH_CA")
  if pathCa == "" {
    log.Fatal("PATH_CA must be set and non-empty")
  }

  fileCa, err := ioutil.ReadFile(pathCa)
  if err != nil {
    log.Fatal(err)
  }

  certPool := x509.NewCertPool()
  certPool.AppendCertsFromPEM(fileCa)

  tlsConfig := &tls.Config{RootCAs: certPool}
  tlsConfig.BuildNameToCertificate()

  transport := &http.Transport{TLSClientConfig: tlsConfig}
  client := &http.Client{Transport: transport}

  resp, err := client.Get(serviceAddr)
  if err != nil {
    log.Fatal(err)
  }
  defer resp.Body.Close()

  data, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Fatal(err)
  }

  log.Println(string(data))
}