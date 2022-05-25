package metricssupport

import (
  "encoding/json"
  "fmt"
  "github.com/gorilla/mux"
  "log"
  "net/http"
  "strconv"
  "sync/atomic"
)

var requests uint64

func HandlerFunction(w http.ResponseWriter, _ *http.Request) {
  w.Header().Set("content-type", "application/json")
  data, _ := json.Marshal(map[string]string{"requests": strconv.FormatUint(requests, 10)})
  _, _ = w.Write(data)
  w.WriteHeader(http.StatusOK)
}

func Middleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    path, _ := mux.CurrentRoute(r).GetPathTemplate()
    log.Println(fmt.Sprintf("collecting metrics for request %s", path))
    atomic.AddUint64(&requests, 1)
    next.ServeHTTP(w, r)
  })
}
