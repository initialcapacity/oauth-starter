package resource

import (
  "embed"
  _ "embed"
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

func (a App) LoadHandlers() func(x *mux.Router) {
  return func(router *mux.Router) {
    router.HandleFunc("/", a.dashboard).Methods("GET")
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

func NewApp(addr string) (*http.Server, net.Listener) {
  if found := os.Getenv("PORT"); found != "" {
    host, _, _ := net.SplitHostPort(addr)
    addr = fmt.Sprintf("%v:%v", host, found)
  }
  listener, _ := net.Listen("tcp", addr)
  return websupport.Create(listener.Addr().String(), App{}.LoadHandlers()), listener
}
