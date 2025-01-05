package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"
)

func setupHTTPProxy(disableAuth bool) (*Tunnel, error) {
	port, err := GetFreePort()
	if err != nil {
		return nil, fmt.Errorf("get free port to listen: %w", err)
	}
	addr := fmt.Sprintf("localhost:%v", port)
	proxyURL := fmt.Sprintf("http://%s", addr)

	// Start HTTP proxy
	proxy := goproxy.NewProxyHttpServer()
	if !disableAuth {
		token := GenerateToken()
		auth.ProxyBasic(proxy, "mitten", func(user, password string) bool {
			return user == "mitten" && password == token
		})
		proxyURL = fmt.Sprintf("http://mitten:%s@%s", token, addr)
	}
	go func() {
		if err := http.ListenAndServe(addr, proxy); err != nil {
			log.Fatalf("listen and serve: %v", err)
		}
	}()

	tunnel := &Tunnel{
		Command:     fmt.Sprintf(`export http_proxy="%s";export https_proxy=$http_proxy;`, proxyURL),
		ForwardSpec: fmt.Sprintf("-R %v:%s", port, addr),
	}
	return tunnel, nil
}
