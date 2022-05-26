package web

import (
  "embed"
  _ "embed"
  "encoding/json"
  "fmt"
  "github.com/gorilla/mux"
  "github.com/gorilla/sessions"
  "github.com/initialcapacity/oauth-starter/pkg/healthsupport"
  "github.com/initialcapacity/oauth-starter/pkg/metricssupport"
  "github.com/initialcapacity/oauth-starter/pkg/pkcesupport"
  "github.com/initialcapacity/oauth-starter/pkg/static"
  "github.com/initialcapacity/oauth-starter/pkg/websupport"
  "io"
  "io/fs"
  "net"
  "net/http"
  "net/url"
  "os"
  "strings"
)

var (
  //go:embed resources
  Resources embed.FS
)

type HTTPClient interface {
  Post(url, contentType string, body io.Reader) (*http.Response, error)
}

type App struct {
  cookieStore  *sessions.CookieStore
  client       HTTPClient
  clientId     string
  clientSecret string
  authUri      string
  tokenUri     string
  callbackUrl  string
  verifier     string
}

func (a *App) LoadHandlers() func(x *mux.Router) {
  return func(router *mux.Router) {
    router.HandleFunc("/", a.dashboard).Methods("GET")
    router.HandleFunc("/login", a.login).Methods("GET")
    router.HandleFunc("/callback", a.callback).Methods("GET")

    router.HandleFunc("/health", healthsupport.HandlerFunction)
    router.HandleFunc("/metrics", metricssupport.HandlerFunction)
    router.Use(metricssupport.Middleware)

    s, _ := fs.Sub(static.Resources, "shared/static")
    fileServer := http.FileServer(http.FS(s))
    router.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))
  }
}

func (a *App) dashboard(writer http.ResponseWriter, request *http.Request) {
  session, _ := a.cookieStore.Get(request, "session")
  accessToken := session.Values["access_token"]
  if accessToken == nil {
    http.Redirect(writer, request, "/login", http.StatusSeeOther)
    return
  }
  _ = websupport.ModelAndView(writer, &Resources, "index", websupport.Model{Map: map[string]any{}})
}

func (a *App) login(writer http.ResponseWriter, _ *http.Request) {
  a.verifier = pkcesupport.CodeVerifier()
  challenge := pkcesupport.CodeChallenge(a.verifier)
  data := map[string]any{
    "client_id":             a.clientId,
    "authentication_url":    a.authUri,
    "response_type":         "code", // (oauth ietf 4.1.1.)
    "redirect_url":          a.callbackUrl,
    "code_challenge":        challenge,
    "code_challenge_method": "S256",
    "scope":                 "openid email",
  }
  _ = websupport.ModelAndView(writer, &Resources, "login", websupport.Model{Map: data})
}

type accessTokenResponse struct {
  AccessToken string `json:"access_token"`
}

func (a *App) callback(writer http.ResponseWriter, request *http.Request) {
  data := url.Values{
    "grant_type":    []string{"authorization_code"},
    "code":          []string{request.URL.Query().Get("code")},
    "client_id":     []string{a.clientId},
    "client_secret": []string{a.clientSecret}, // todo - not needed for pkce
    "redirect_url":  []string{a.callbackUrl},
    "code_verifier": []string{a.verifier}, // todo - code_verifier replaces the client_secret
  }

  // todo - move code to authorization header
  // todo - check status code
  post, _ := a.client.Post(a.tokenUri, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

  all, _ := io.ReadAll(post.Body)
  var tok accessTokenResponse
  _ = json.Unmarshal(all, &tok)

  session, _ := a.cookieStore.Get(request, "session")
  session.Options.MaxAge = 60
  session.Values["access_token"] = tok.AccessToken
  _ = session.Save(request, writer)

  // todo - check for errors

  http.Redirect(writer, request, "/", http.StatusSeeOther)
}

type credentials struct {
  ClientId     string `json:"client_id"`
  ClientSecret string `json:"client_secret"` // todo - not needed for pkce
  AuthUri      string `json:"auth_uri"`
  TokenUri     string `json:"token_uri"`
  CallbackUrl  string `json:"callback_url"`
}

func NewApp(addr string, client HTTPClient) (*http.Server, net.Listener) {
  if found := os.Getenv("PORT"); found != "" {
    host, _, _ := net.SplitHostPort(addr)
    addr = fmt.Sprintf("%v:%v", host, found)
  }
  key := "super_private"
  if found := os.Getenv("SESSION_KEY"); found != "" {
    key = found
  }
  credentialsJson := os.Getenv("CREDENTIALS_FILE")
  if credentialsJson == "" {
    panic("missing credentials file.")
  }

  var c credentials
  _ = json.NewDecoder(strings.NewReader(credentialsJson)).Decode(&c)
  cookieStore := sessions.NewCookieStore([]byte(os.Getenv(key)))
  listener, _ := net.Listen("tcp", addr)
  app := App{cookieStore: cookieStore, client: client,
    clientId:     c.ClientId,
    clientSecret: c.ClientSecret,
    authUri:      c.AuthUri,
    tokenUri:     c.TokenUri,
    callbackUrl:  c.CallbackUrl,
  }
  return websupport.Create(listener.Addr().String(), app.LoadHandlers()), listener
}
