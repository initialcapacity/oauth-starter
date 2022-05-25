package authorization_test

import (
  "fmt"
  "github.com/initialcapacity/oauth-starter/pkg/authorization"
  "github.com/initialcapacity/oauth-starter/pkg/testsupport"
  "github.com/initialcapacity/oauth-starter/pkg/websupport"
  "github.com/stretchr/testify/assert"
  "io"
  "net/http"
  "net/url"
  "os"
  "strings"
  "testing"
)

func TestAuthentication(t *testing.T) {
  _ = os.Setenv("PORT", "0")
  server, listener := authorization.NewApp("localhost:0")
  go websupport.Start(server, listener)
  testsupport.WaitForHealthy(server, "/health")

  get, _ := http.Get(fmt.Sprintf("http://%s/", server.Addr))
  getBody, _ := io.ReadAll(get.Body)
  assert.Contains(t, string(getBody), "Authorization server")

  postForAuth, _ := http.Get(fmt.Sprintf("http://%s/auth?response_type=code&client_id=%s&redirect_url=%s",
    server.Addr, "aClientId", "http://localhost:8889/callback"))
  body, _ := io.ReadAll(postForAuth.Body)
  assert.Contains(t, string(body), "<p>Sign in to continue to the client application.</p>")

  signInData := url.Values{
    "client_id":    []string{"aClientId"},
    "redirect_url": []string{"http://localhost:8889/callback"},
  }
  client := &http.Client{
    CheckRedirect: func(req *http.Request, via []*http.Request) error {
      return http.ErrUseLastResponse
    },
  }
  postForSignIn, _ := client.Post(fmt.Sprintf("http://%s/signin", server.Addr), "application/x-www-form-urlencoded", strings.NewReader(signInData.Encode()))
  assert.Equal(t, http.StatusFound, postForSignIn.StatusCode)
  assert.Equal(t, "http://localhost:8889/callback?code=42", postForSignIn.Header.Get("Location"))

  data := url.Values{
    "grant_type":    []string{"authorization_code"},
    "code":          []string{"42"},
    "client_id":     []string{"aClientId"},
    "client_secret": []string{"aClientSecret"},
    "redirect_url":  []string{"http://localhost:8889/callback"},
  }

  postForToken, _ := client.Post(fmt.Sprintf("http://%s/token", server.Addr), "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
  assert.Equal(t, http.StatusCreated, postForToken.StatusCode)

  websupport.Stop(server)
}
