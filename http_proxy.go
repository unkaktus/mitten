package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"
)

func setupHTTPProxy() (*Tunnel, error) {
	token := GenerateToken()

	port, err := GetFreePort()
	if err != nil {
		return nil, fmt.Errorf("get free port to listen: %w", err)
	}
	addr := fmt.Sprintf("localhost:%v", port)

	// Start HTTP proxy
	proxy := goproxy.NewProxyHttpServer()
	auth.ProxyBasic(proxy, "mitten", func(user, password string) bool {
		return user == "mitten" && password == token
	})
	go func() {
		if err := http.ListenAndServe(addr, proxy); err != nil {
			log.Fatalf("listen and serve: %v", err)
		}
	}()

	tunnel := &Tunnel{
		Command:     fmt.Sprintf(`export http_proxy="http://mitten:%s@%s";export https_proxy=$http_proxy;`, token, addr),
		ForwardSpec: fmt.Sprintf("-R %v:%s", port, addr),
	}
	return tunnel, nil
}
