package web_test

import (
  "bytes"
  "encoding/json"
  "fmt"
  "github.com/initialcapacity/oauth-starter/pkg/testsupport"
  "github.com/initialcapacity/oauth-starter/pkg/web"
  "github.com/initialcapacity/oauth-starter/pkg/websupport"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/mock"
  "io"
  "io/ioutil"
  "net/http"
  "os"
  "testing"
)

type MockClient struct {
  mock.Mock
}

func (m *MockClient) Post(_, _ string, _ io.Reader) (*http.Response, error) {
  resp, _ := json.Marshal(struct {
    AccessToken string `json:"access_token"`
  }{"anAccessToken"})
  return &http.Response{StatusCode: 201, Body: ioutil.NopCloser(bytes.NewReader(resp))}, nil
}

func TestMissingCredential(t *testing.T) {
  defer func() {
    if err := recover(); err == nil {
      t.Fail()
    }
  }()
  _, _ = web.NewApp("localhost:0", &MockClient{})
}

func TestApp(t *testing.T) {
  _ = os.Setenv("CREDENTIALS_FILE", `{
  "client_id":"aClientId",
  "client_secret":"aClientSecret",
  "auth_uri":"http://localhost:8877/auth",
  "token_uri":"http://localhost:8877/token",
  "callback_url":"http://localhost:8879/callback"
}`)

  _ = os.Setenv("PORT", "0")
  _ = os.Setenv("SESSION_KEY", "super_secret")
  server, listener := web.NewApp("localhost:0", &MockClient{})
  go websupport.Start(server, listener)
  testsupport.WaitForHealthy(server, "/health")

  get, _ := http.Get(fmt.Sprintf("http://%s/", server.Addr))
  body, _ := io.ReadAll(get.Body)
  assert.Contains(t, string(body), "Sign in with the Authorization Server")

  getLogin, _ := http.Get(fmt.Sprintf("http://%s/login", server.Addr))
  bodyLogin, _ := io.ReadAll(getLogin.Body)
  assert.Contains(t, string(bodyLogin),
    "<a href=\"http://localhost:8877/auth?response_type=code&client_id=aClientId&redirect_url=http%3a%2f%2flocalhost%3a8879%2fcallback&scope=openid%20email&code_challenge_method=S256&code_challenge=")

  client := &http.Client{
    CheckRedirect: func(req *http.Request, via []*http.Request) error {
      return http.ErrUseLastResponse
    },
  }

  getCallback, _ := client.Get(fmt.Sprintf("http://%s/callback?access_token=42", server.Addr))
  request, _ := http.NewRequest("GET", fmt.Sprintf("http://%s", server.Addr), nil)
  request.AddCookie(getCallback.Cookies()[0])
  response, _ := client.Do(request)

  responseBody, _ := io.ReadAll(response.Body)
  assert.Contains(t, string(responseBody), "OAuth starter")
  assert.NotContains(t, string(responseBody), "Sign in with the Authorization Server")

  websupport.Stop(server)
}
