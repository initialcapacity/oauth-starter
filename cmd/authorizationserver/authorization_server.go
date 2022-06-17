package main

import (
	"github.com/initialcapacity/oauth-starter/pkg/authorization"
	"github.com/initialcapacity/oauth-starter/pkg/websupport"
)

func main() {
	websupport.Start(authorization.NewApp("0.0.0.0:8887"))
}
