package main

import (
	"fmt"
	"io"
	"log"

	"github.com/gliderlabs/ssh"
	"github.com/pkg/sftp"
)

func sftpHandler(sess ssh.Session) {
	debugStream := io.Discard
	serverOptions := []sftp.ServerOption{
		sftp.WithDebug(debugStream),
	}
	server, err := sftp.NewServer(
		sess,
		serverOptions...,
	)
	if err != nil {
		log.Printf("sftp server init error: %s\n", err)
		return
	}
	if err := server.Serve(); err == io.EOF {
		server.Close()
	} else if err != nil {
		fmt.Println("sftp server completed with error:", err)
	}
}

func runSFTPServer(addr, token string) error {
	sshServer := ssh.Server{
		Addr: addr,
		Handler: func(s ssh.Session) {
			fmt.Fprintf(s, "SCP is not supported, use SFTP instead.\n")
			s.Exit(1)
		},
		PublicKeyHandler: func(ctx ssh.Context, key ssh.PublicKey) bool {
			user := ctx.User()
			if user != token {
				panic(fmt.Sprintf("user:'%s', token:'%s'", user, token))
			}
			return user == token
		},
		SubsystemHandlers: map[string]ssh.SubsystemHandler{
			"sftp": sftpHandler,
		},
	}
	return sshServer.ListenAndServe()
}

func setupSFTP() (*Tunnel, error) {
	token := GenerateToken()

	port, err := GetFreePort()
	if err != nil {
		return nil, fmt.Errorf("get free port to listen: %w", err)
	}
	addr := fmt.Sprintf("localhost:%d", port)

	go func() {
		if err := runSFTPServer(addr, token); err != nil {
			log.Fatalf("run SFTP server: %v", err)
		}
	}()

	// Generate a key to enable authentication
	keygenCommand := "ssh-keygen -t ed25519 -f $HOME/.ssh/mitten_key -N '' -q  <<<y >/dev/null 2>&1"
	// Invoke sftp
	sftpCommand := fmt.Sprintf("sftp -o Port=%d -o User=%s -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -i $HOME/.ssh/mitten_key -q localhost", port, token)

	tunnel := &Tunnel{
		Command:     fmt.Sprintf("mittenfs(){ %s; %s; };", keygenCommand, sftpCommand),
		ForwardSpec: fmt.Sprintf("-R %v:%s", port, addr),
	}
	return tunnel, nil
}
