package controller

import (
    "bytes"
    "flag"
    "io"
    "net/http"
    "strconv"
)

var cacheServiceURL = flag.String("cacheservice", "https://172.18.0.13:5000", "Address of caching service")

func getFromCache(key string) (io.ReadCloser, bool) {
    resp, err := http.Get(*cacheServiceURL + "/?key=" + key)
    if err != nil || resp.StatusCode != http.StatusOK {
        println("get fail")
        return nil, false
    }
    return resp.Body, true
}

func saveToCache(key string, duration int64, data []byte) {
    req, _ := http.NewRequest(http.MethodPost, *cacheServiceURL+"/?key="+key, bytes.NewBuffer(data))
    req.Header.Add("cacche-control", "maxage="+strconv.FormatInt(duration, 10))
    http.DefaultClient.Do(req)
}

func invalidateCacheEntry(key string) {
    http.Get(*cacheServiceURL + "/invalidate?key=" + key)
}
