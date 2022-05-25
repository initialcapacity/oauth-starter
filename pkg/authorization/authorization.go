package authorization

import (
  "embed"
  _ "embed"
  "encoding/json"
  "fmt"
  "github.com/gorilla/mux"
  "github.com/initialcapacity/oauth-starter/pkg/healthsupport"
  "github.com/initialcapacity/oauth-starter/pkg/metricssupport"
  "github.com/initialcapacity/oauth-starter/pkg/static"
  "github.com/initialcapacity/oauth-starter/pkg/websupport"
  "io/fs"
  "net"
  "net/http"
  "os"
)

var (
  //go:embed resources
  Resources embed.FS
)

type App struct {
}

func (a *App) LoadHandlers() func(x *mux.Router) {
  return func(router *mux.Router) {
    router.HandleFunc("/", a.dashboard).Methods("GET")
    router.HandleFunc("/auth", a.authenticate).Methods("GET")
    router.HandleFunc("/signin", a.signIn).Methods("POST")
    router.HandleFunc("/token", a.token).Methods("POST")

    router.HandleFunc("/health", healthsupport.HandlerFunction)
    router.HandleFunc("/metrics", metricssupport.HandlerFunction)
    router.Use(metricssupport.Middleware)

    s, _ := fs.Sub(static.Resources, "shared/static")
    fileServer := http.FileServer(http.FS(s))
    router.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))
  }
}

func (a *App) dashboard(writer http.ResponseWriter, _ *http.Request) {
  _ = websupport.ModelAndView(writer, &Resources, "index", websupport.Model{Map: map[string]any{}})
}

func (a *App) authenticate(writer http.ResponseWriter, request *http.Request) {
  data := map[string]any{
    "client_id":    request.URL.Query().Get("client_id"),
    "redirect_url": request.URL.Query().Get("redirect_url"),
  }
  _ = websupport.ModelAndView(writer, &Resources, "grant_access", websupport.Model{Map: data})
}

func (a *App) signIn(writer http.ResponseWriter, request *http.Request) {
  _ = request.ParseForm()
  redirectUrl := request.Form.Get("redirect_url")
  // todo - confirm client_id and redirect_url
  // todo - confirm username and password

  accessCode := "42"

  http.Redirect(writer, request, redirectUrl+"?code="+accessCode, http.StatusFound)
}

func (a *App) token(writer http.ResponseWriter, request *http.Request) {
  // todo - authenticate the client - client_id, client_secret
  // todo - ensure that the authorization code was issued to the authenticated confidential client
  // todo - verify that the authorization code is valid - code, client_id, redirect_url
  // todo - ensure that the "redirect_uri" parameter is present

  data, _ := json.Marshal(struct {
    AccessToken string `json:"access_token"`
  }{"anAccessToken"})

  writer.WriteHeader(http.StatusCreated)
  _, _ = writer.Write(data)
}

func NewApp(addr string) (*http.Server, net.Listener) {
  if found := os.Getenv("PORT"); found != "" {
    host, _, _ := net.SplitHostPort(addr)
    addr = fmt.Sprintf("%v:%v", host, found)
  }

  listener, _ := net.Listen("tcp", addr)
  app := App{}
  return websupport.Create(listener.Addr().String(), app.LoadHandlers()), listener
}
