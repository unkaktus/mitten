package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"

	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"
)

func GetFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

func GenerateToken() string {
	size := 16
	rb := make([]byte, size)
	_, err := rand.Read(rb)

	if err != nil {
		panic(err)
	}

	rs := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(rb)
	return rs
}

func run() error {
	if len(os.Args) == 1 {
		log.Fatalf("Specify the host")
	}

	token := GenerateToken()

	port, err := GetFreePort()
	if err != nil {
		log.Fatalf("get free port to listen: %v", err)
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

	cmdline := []string{
		"-t",                             // Force pty allocation
		"-o", "ExitOnForwardFailure=yes", // Exit on forwarding failure
		fmt.Sprintf("-R %v:%s", port, addr), // Forward the proxy port
	}
	cmdline = append(cmdline, os.Args[1:]...) // Add all that user specified

	// Export the environment variables and start a shell
	cmdline = append(cmdline,
		fmt.Sprintf("export http_proxy=\"http://mitten:%s@%s\"; export https_proxy=\"http://mitten:%s@%s\"; exec $SHELL -l", token, addr, token, addr),
	)

	cmd := exec.Command("ssh", cmdline...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("execute ssh: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}

}
