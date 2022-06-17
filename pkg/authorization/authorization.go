package authorization

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/initialcapacity/oauth-starter/pkg/healthsupport"
	"github.com/initialcapacity/oauth-starter/pkg/metricssupport"
	"github.com/initialcapacity/oauth-starter/pkg/pkcesupport"
	"github.com/initialcapacity/oauth-starter/pkg/static"
	"github.com/initialcapacity/oauth-starter/pkg/websupport"
	"io/fs"
	"net"
	"net/http"
	"os"
	"strings"
)

var (
	//go:embed resources
	Resources embed.FS
)

type App struct {
	codeChallenge string
}

func (a *App) LoadHandlers() func(x *mux.Router) {
	return func(router *mux.Router) {
		router.HandleFunc("/", a.dashboard).Methods("GET")
		router.HandleFunc("/auth", a.authenticate).Methods("GET")
		router.HandleFunc("/signin", a.signIn).Methods("POST")
		router.HandleFunc("/token", a.token).Methods("POST", "OPTIONS")

		router.HandleFunc("/health", healthsupport.HandlerFunction)
		router.HandleFunc("/metrics", metricssupport.HandlerFunction)
		router.Use(metricssupport.Middleware)
		router.Use(mux.CORSMethodMiddleware(router))

		s, _ := fs.Sub(static.Resources, "shared/static")
		fileServer := http.FileServer(http.FS(s))
		router.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))
	}
}

func (a *App) dashboard(writer http.ResponseWriter, _ *http.Request) {
	_ = websupport.ModelAndView(writer, &Resources, "index", websupport.Model{Map: map[string]any{}})
}

func (a *App) authenticate(writer http.ResponseWriter, request *http.Request) {
	scope := request.URL.Query().Get("scope")
	data := map[string]any{
		"client_id":      request.URL.Query().Get("client_id"),
		"redirect_url":   request.URL.Query().Get("redirect_url"),
		"code_challenge": request.URL.Query().Get("code_challenge"),
		"scope":          scope,
	}
	// todo - confirm params were received
	// todo - verify client_id and redirect_url
	// todo - verify that a scope parameter is present and contains the openid scope value (openid 3.1.2.2)
	if scope == "" || !strings.Contains(scope, "openid") {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	_ = websupport.ModelAndView(writer, &Resources, "grant_access", websupport.Model{Map: data})
}

func (a *App) signIn(writer http.ResponseWriter, request *http.Request) {
	_ = request.ParseForm()
	redirectUrl := request.Form.Get("redirect_url")
	codeChallenge := request.Form.Get("code_challenge")
	// todo - confirm params were received; client_id, code_challenge and redirect_url
	// todo - verify client_id and redirect_url
	// todo - authenticate; username and password
	// todo - associate the "code_challenge" and "code_challenge_method" values with the authorization code (pkce ietf 4.4)

	a.codeChallenge = codeChallenge
	authorizationCode := "42"

	http.Redirect(writer, request, redirectUrl+"?code="+authorizationCode, http.StatusFound) // (oauth ietf 4.1.2)
}

func (a *App) token(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	if request.Method == http.MethodOptions {
		writer.WriteHeader(http.StatusOK)
		return
	}

	// todo - confirm params were received; authorization_code, client_id, code_verifier and redirect_url
	// todo - ensure that the authorization code was issued to the authenticated confidential client
	// todo - verify that the authorization code is valid - code, client_id, code_verifier, and redirect_url
	// todo - ensure that the "redirect_uri" parameter is present

	_ = request.ParseForm()
	verifier := request.Form.Get("code_verifier") // prevents from CSRF and authorization code attacks
	if a.codeChallenge != pkcesupport.CodeChallenge(verifier) {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// todo - move to json web token
	data, err := json.Marshal(struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		IdToken      string `json:"id_token"`
	}{
		"anAccessToken",
		"Bearer",
		3600,
		"aRefreshToken",
		"anIdToken",
	})
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated) // todo - check this
	_, err = writer.Write(data)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
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
