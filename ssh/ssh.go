package ssh

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
)

type Remote struct {
	addr      string
	sshConfig *ssh.ClientConfig
}

func NewRemote(addr, user, password string) (*Remote, error) {
	return &Remote{
		addr: addr,
		sshConfig: &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
	}, nil
}

type SessionOpts interface {
	apply(*ssh.Session) error
}

type WithEnv map[string]string

func (w WithEnv) apply(s *ssh.Session) error {
	for k, v := range w {
		if err := s.Setenv(k, v); err != nil {
			return err
		}
	}

	return nil
}

func (r *Remote) Run(cmd string, options ...SessionOpts) error {
	sshClient, err := ssh.Dial("tcp", r.addr, r.sshConfig)
	if err != nil {
		return err
	}

	session, err := sshClient.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	for _, option := range options {
		if err := option.apply(session); err != nil {
			return err
		}
	}

	out, err := session.CombinedOutput(cmd)
	if err != nil {
		log.Printf("error running command: %s", out)

		return err
	}

	return nil
}

func (r *Remote) Copy(ctx context.Context, src, dest string) error {
	sshClient, err := ssh.Dial("tcp", r.addr, r.sshConfig)
	if err != nil {
		return err
	}

	scpClient, err := scp.NewClientBySSH(sshClient)
	if err != nil {
		return err
	}

	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	if err := scpClient.CopyFromFile(ctx, *f, dest, "0655"); err != nil {
		return fmt.Errorf("copying file: %w", err)
	}

	return nil
}

func (r *Remote) Read(ctx context.Context, src string) ([]byte, error) {
	sshClient, err := ssh.Dial("tcp", r.addr, r.sshConfig)
	if err != nil {
		return nil, err
	}

	scpClient, err := scp.NewClientBySSH(sshClient)
	if err != nil {
		return nil, err
	}

	tmpFile, err := os.CreateTemp("", "uberfx")
	if err != nil {
		return nil, fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if err := scpClient.CopyFromRemote(ctx, tmpFile, src); err != nil {
		return nil, fmt.Errorf("copying file: %w", err)
	}

	return os.ReadFile(tmpFile.Name())
}

func (r *Remote) Write(ctx context.Context, dest, content string) error {
	sshClient, err := ssh.Dial("tcp", r.addr, r.sshConfig)
	if err != nil {
		return err
	}

	scpClient, err := scp.NewClientBySSH(sshClient)
	if err != nil {
		return err
	}

	src := strings.NewReader(content)

	if err := scpClient.CopyFile(ctx, src, dest, "0655"); err != nil {
		return fmt.Errorf("copying file: %w", err)
	}

	return nil
}
