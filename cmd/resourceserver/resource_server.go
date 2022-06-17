package main

import (
	"github.com/initialcapacity/oauth-starter/pkg/resource"
	"github.com/initialcapacity/oauth-starter/pkg/websupport"
)

func main() {
	websupport.Start(resource.NewApp("0.0.0.0:8888"))
}
