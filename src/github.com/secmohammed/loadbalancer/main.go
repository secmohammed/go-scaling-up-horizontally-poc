package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

type webRequest struct {
    r *http.Request
    w http.ResponseWriter
    doneCh chan struct{}
}
var (
    requestCh = make(chan *webRequest)
    // Bypass the development cert and allow using https. Useful for development, but causes an issue for production, as you will bypass the cert.
    // The usage of transporting to https in dev mode is that since go 1.6.0, it switches automatically to http/2 which applies header compression automatically.
    transport = http.Transport{
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: true,
        },
    }

)
func init() {
    http.DefaultClient = &http.Client{Transport: &transport}
}
func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        doneCh := make(chan struct{})
        requestCh <- &webRequest{r: r, w: w, doneCh: doneCh}
        // wait for doneCh. indicates that response is received.
        <-doneCh
    })
    go processRequests()
    go http.ListenAndServeTLS(":2000", "cert.pem", "key.pem", nil)
    log.Println("Server Started!, press <ENTER> to exit")
    fmt.Scanln()
}
