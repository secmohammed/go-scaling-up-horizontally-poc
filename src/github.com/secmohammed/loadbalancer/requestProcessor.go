package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
    appservers = []string{}
    currentIndex = 0
    client = http.Client{Transport: &transport}
)

// The idea behind this basic loadbalancer is to distribute the requests equally on the servers registered.
func processRequests() {
    for {
        select {
        case <-heartbeatCh:
            println("hearbeat")
            server := appservers[:]
            go func(servers []string) {
                for _, server := range servers {
                    resp, err := http.Get("https://" + server + "/ping")
                    if err != nil || resp.StatusCode != 200 {
                        unregisterCh <- server
                    }
                }
            }(servers)
        case host := <-unregisterCh:
            fmt.Println("unregister " + host)
            for i := len(appservers) - 1; i >= 0; i-- {
                if appservers[i] == host {
                    appservers = append(appservers[:i], appservers[i+1:]...)
                }
            }
        case host := <-registerCh:
            fmt.Println("register " + host)
            isFound := false
            for _, h := range appservers {
                if h == host {
                    isFound = true
                    break
                }
            }
            if !isFound {
                appservers = append(appservers, host)
            }


        case request := <-requestCh:
            fmt.Println("request")
            if len(appservers) == 0 {
                request.w.WriteHeader(http.StatusInternalServerError)
                request.w.Write([]byte("No app servers found"))
                request.doneCh <- struct{}{}
                continue
            }
            currentIndex++
            if currentIndex == len(appservers) {
                currentIndex = 0
            }
            host := appservers[currentIndex]
            go processRequest(host, request)

        }
    }
}

func processRequest(host string, request *webRequest) {
    hostURL, _ := url.Parse(request.r.URL.String())
    hostURL.Scheme = "https"
    hostURL.Host = host
    fmt.Println(host)
    fmt.Println(hostURL.String())
    req, _ := http.NewRequest(request.r.Method, hostURL.String(), request.r.Body)

    for k, v := range request.r.Header {
        values := ""
        for _, headerValue := range v {
            values += headerValue + " "
        }
        req.Header.Add(k, values)
    }
    resp, err := client.Do(req)
    if err != nil {
        request.w.WriteHeader(http.StatusInternalServerError)
        request.doneCh <- struct{}{}
        return

    }
    // This will loop through all of the returned headers of the response from the server which recieved the request
    // The headers will expose which server served the http request, and other headers.
    // TODO: we will have to go through the response headers and filter the ones that exposes the server info.
    for k, v := range resp.Header {
        values := ""
        for _, headerValue := range v {
            values += headerValue + " "

        }
        request.w.Header().Add(k, values)
    }
    // copy the response from the server which we made through the http request we fired internally at the function.
    // to the request response writer struct which have the original request and response.
    // which means in essence, we are proxying the request/repsonse through the loadbalancer
    // and copy the response we made internally to the original response
    // note that request is actually a pointer since it's a channel.
    io.Copy(request.w, resp.Body)
    request.doneCh <- struct{}{}

}
