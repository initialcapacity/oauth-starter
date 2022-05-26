package smoke_test

import (
  "github.com/initialcapacity/oauth-starter/pkg/authorization"
  "github.com/initialcapacity/oauth-starter/pkg/resource"
  "github.com/initialcapacity/oauth-starter/pkg/testsupport"
  "github.com/initialcapacity/oauth-starter/pkg/web"
  "github.com/initialcapacity/oauth-starter/pkg/websupport"
  "github.com/stretchr/testify/assert"
  "io"
  "net/http"
  "net/url"
  "os"
  "regexp"
  "strings"
  "testing"
)

func TestAuthentication(t *testing.T) {
  _ = os.Setenv("CREDENTIALS_FILE", `{
  "client_id":"101010",
  "client_secret":"super_private",
  "auth_uri":"http://localhost:18877/auth",
  "token_uri":"http://localhost:18877/token",
  "callback_url":"http://localhost:18879/callback"
}`)

  authorizationServer, authorizationListener := authorization.NewApp("localhost:18877")
  resourceServer, resourcesListener := resource.NewApp("localhost:18876")
  webserver, webListener := web.NewApp("localhost:18879", &http.Client{})

  go websupport.Start(authorizationServer, authorizationListener)
  go websupport.Start(resourceServer, resourcesListener)
  go websupport.Start(webserver, webListener)

  testsupport.WaitForHealthy(authorizationServer, "/health")
  testsupport.WaitForHealthy(resourceServer, "/health")
  testsupport.WaitForHealthy(webserver, "/health")

  getHomepage, _ := http.Get("http://localhost:18879")
  homepage, _ := io.ReadAll(getHomepage.Body)
  assert.Contains(t, string(homepage), "OAuth starter")
  assert.Contains(t, string(homepage), "Sign in with the Authorization Server")
  codeChallenge := regexp.MustCompile("&code_challenge=(.*?)").FindStringSubmatch(string(homepage))[1]

  // click sign in...

  values := url.Values{
    "username":       []string{"jerry.cantrell@gmail.com"},
    "password":       []string{"boogydepot"},
    "client_id":      []string{"aClientId"},
    "redirect_url":   []string{"http://localhost:18879/callback"},
    "code_challenge": []string{codeChallenge},
  }
  var cookie *http.Cookie
  client := &http.Client{
    CheckRedirect: func(req *http.Request, via []*http.Request) error {
      if cookies := req.Response.Cookies(); len(cookies) > 0 {
        cookie = cookies[0]
      }
      return nil
    },
  }
  signIn, _ := client.Post("http://localhost:18877/signin", "application/x-www-form-urlencoded", strings.NewReader(values.Encode()))
  assert.Equal(t, http.StatusOK, signIn.StatusCode)

  homepageRequest, _ := http.NewRequest("GET", "http://localhost:18879", nil)
  homepageRequest.AddCookie(cookie)
  getHomepage, _ = client.Do(homepageRequest)

  homepage, _ = io.ReadAll(getHomepage.Body)
  assert.Contains(t, string(homepage), "OAuth starter")
  assert.NotContains(t, string(homepage), "Sign in with the Authorization Server")

  websupport.Stop(webserver)
  websupport.Stop(resourceServer)
  websupport.Stop(authorizationServer)
}
