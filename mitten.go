package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/elazarl/goproxy"
)

func run() error {
	port := "31773"
	addr := "localhost:" + port

	// Start HTTP proxy
	proxy := goproxy.NewProxyHttpServer()
	go func() {
		if err := http.ListenAndServe(addr, proxy); err != nil {
			log.Fatalf("listen and serve: %v", err)
		}
	}()

	cmdline := []string{
		"-t",                             // Force pty allocation
		"-o", "ExitOnForwardFailure=yes", // Exit on forwarding failure
		fmt.Sprintf("-R %s:%s", port, addr), // Forward the proxy port
	}
	cmdline = append(cmdline, os.Args[1:]...) // Add all that user specified

	// Export the environment variables and start a shell
	cmdline = append(cmdline,
		fmt.Sprintf("export http_proxy=\"http://%s\"; export https_proxy=\"http://%s\"; exec $SHELL -l", addr, addr),
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
