package main

import (
	"github.com/initialcapacity/oauth-starter/pkg/web"
	"github.com/initialcapacity/oauth-starter/pkg/websupport"
	"net/http"
)

func main() {
	websupport.Start(web.NewApp("0.0.0.0:8889", &http.Client{}))
}
