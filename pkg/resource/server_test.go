package resource_test

import (
  "fmt"
  "github.com/initialcapacity/oauth-starter/pkg/resource"
  "github.com/initialcapacity/oauth-starter/pkg/testsupport"
  "github.com/initialcapacity/oauth-starter/pkg/websupport"
  "github.com/stretchr/testify/assert"
  "io"
  "net/http"
  "os"
  "testing"
)

func TestResource(t *testing.T) {
  _ = os.Setenv("PORT", "0")
  server, listener := resource.NewApp("localhost:0")
  go websupport.Start(server, listener)
  testsupport.WaitForHealthy(server, "/health")

  get, _ := http.Get(fmt.Sprintf("http://%s/", server.Addr))
  getBody, _ := io.ReadAll(get.Body)
  assert.Contains(t, string(getBody), "Resource server")
}
