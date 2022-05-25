package websupport_test

import (
  "github.com/gorilla/mux"
  "github.com/initialcapacity/oauth-starter/pkg/healthsupport"
  "github.com/initialcapacity/oauth-starter/pkg/testsupport"
  "github.com/initialcapacity/oauth-starter/pkg/websupport"
  "github.com/stretchr/testify/assert"
  "net"
  "testing"
)

func TestServer(t *testing.T) {
  listener, _ := net.Listen("tcp", "localhost:0")
  server := websupport.Create(listener.Addr().String(), func(router *mux.Router) {
    router.HandleFunc("/health", healthsupport.HandlerFunction).Methods("GET")
  })
  go websupport.Start(server, listener)
  assert.True(t, testsupport.WaitForHealthy(server, "health"))
  websupport.Stop(server)
}
