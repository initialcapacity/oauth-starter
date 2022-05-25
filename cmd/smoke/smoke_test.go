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
  "strings"
  "testing"
)

func TestAuthentication(t *testing.T) {
  authorizationServer, authorizationListener := authorization.NewApp("localhost:18877")
  go websupport.Start(authorizationServer, authorizationListener)
  testsupport.WaitForHealthy(authorizationServer, "/health")

  resourceServer, resourcesListener := resource.NewApp("localhost:18876")
  go websupport.Start(resourceServer, resourcesListener)
  testsupport.WaitForHealthy(resourceServer, "/health")

  _ = os.Setenv("CREDENTIALS_FILE", `{
  "client_id":"101010",
  "client_secret":"super_private",
  "auth_uri":"http://localhost:18877/auth",
  "token_uri":"http://localhost:18877/token",
  "callback_url":"http://localhost:18879/callback"
}`)

  webserver, webListener := web.NewApp("localhost:18879", &http.Client{})
  go websupport.Start(webserver, webListener)
  testsupport.WaitForHealthy(webserver, "/health")

  signInData := url.Values{
    "client_id":    []string{"aClientId"},
    "redirect_url": []string{"http://localhost:18879/callback"},
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
  resp, _ := client.Post("http://localhost:18877/signin", "application/x-www-form-urlencoded", strings.NewReader(signInData.Encode()))
  assert.Equal(t, http.StatusOK, resp.StatusCode)

  request, _ := http.NewRequest("GET", "http://localhost:18879", nil)
  request.AddCookie(cookie)
  response, _ := client.Do(request)

  responseBody, _ := io.ReadAll(response.Body)
  assert.Contains(t, string(responseBody), "Client application")
  assert.NotContains(t, string(responseBody), "Sign in with the Authorization Server")

  websupport.Stop(webserver)
  websupport.Stop(resourceServer)
  websupport.Stop(authorizationServer)
}
