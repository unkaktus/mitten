package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/term"
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

const banner string = `
       ▗▟▀▀▙
      ▗▛   ▐▌
    ▗▟▘   ▗▛
▗▄▄▟▀     ▀▀▀▀▀▀▀▜▄
█  █              ▝▜▖
█  █      mitten    ▙
█  █      magic!   ▗▌
▜▄▄█▄            ▗▟▀
     ▀▀▀▀▄▄▄▄▄▄▄▀▀      
`

var bannerHeight int = strings.Count(banner, "\n")

type Tunnel struct {
	Command     string
	ForwardSpec string
}

func run() error {
	if len(os.Args) == 1 {
		return fmt.Errorf("no host specified")
	}

	httpProxyTunnel, err := setupHTTPProxy()
	if err != nil {
		return fmt.Errorf("create HTTP proxy: %w", err)
	}

	sftpTunnel, err := setupSFTP()
	if err != nil {
		return fmt.Errorf("create HTTP proxy: %w", err)
	}

	cmdline := []string{
		"-t",                             // Force pty allocation
		"-o", "ExitOnForwardFailure=yes", // Exit on forwarding failure
		httpProxyTunnel.ForwardSpec, // Forward the proxy port
		sftpTunnel.ForwardSpec,      // Forward the SFTP port
	}
	cmdline = append(cmdline, os.Args[1:]...) // Add all that user specified

	cmd := exec.Command("ssh", cmdline...)

	// Start the command with a pty.
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("execute ssh: %w", err)
	}
	// Make sure to close the pty at the end.
	defer func() { ptmx.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH                        // Initial resize.
	defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Export the environment variables
	mittenCommand := fmt.Sprintf(` %s%secho -e '%s';#mitten_marker#`, httpProxyTunnel.Command, sftpTunnel.Command, banner)

	shellFinder := NewShellFindReader(ptmx)

	// Copy stdin to the pty and the pty to stdout.
	go func() {
		_, _ = io.Copy(ptmx, os.Stdin)
	}()

	go func() {
		<-shellFinder.Found
		_, err := io.Copy(ptmx, strings.NewReader(mittenCommand))
		if err != nil {
			log.Fatalf("write mitten command to the remote: %v", err)
		}
		<-shellFinder.SkippedEcho
		_, err = ptmx.Write([]byte("\n"))
		if err != nil {
			log.Fatalf("write mitten command to the remote: %v", err)
		}
	}()
	_, _ = io.Copy(os.Stdout, shellFinder)

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "mitten: %v\n", err)
		os.Exit(1)
	}
}
