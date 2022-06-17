package testsupport

import (
	"fmt"
	"net"
	"net/http"
)

func Server[T http.Handler](router T) *http.Server {
	listener, _ := net.Listen("tcp", "localhost:0")
	server := &http.Server{
		Addr:    listener.Addr().String(),
		Handler: router,
	}
	go func() {
		_ = server.Serve(listener)
	}()
	return server
}

func WaitForHealthy(server *http.Server, path string) bool {
	var isLive bool
	for !isLive {
		resp, err := http.Get(fmt.Sprintf("http://%s/%s", server.Addr, path))
		if err == nil && resp.StatusCode == http.StatusOK {
			isLive = true
		}
	}
	return isLive
}
