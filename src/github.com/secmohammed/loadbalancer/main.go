package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type webRequest struct {
    r *http.Request
    w http.ResponseWriter
    doneCh chan struct{}
}
var (
    requestCh = make(chan *webRequest)
    registerCh = make(chan string)
    unregisterCh = make(chan string)
    heartbeatCh = time.Tick(5 * time.Second)
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
    go http.ListenAndServeTLS(":2001", "cert.pem", "key.pem", new(appserverHandler))
    log.Println("Server Started!, press <ENTER> to exit")
    fmt.Scanln()
}
type appserverHandler struct {

}
func (h *appserverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ip := strings.Split(r.RemoteAddr, ":")[0]
    port := r.URL.Query().Get("port") // incoming port
    switch r.URL.Path {
    case "/register":
        registerCh <- ip + ":" + port
    case "/unregister":
        unregisterCh <- ip + ":" + port
    }


}
