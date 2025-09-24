package ssh

import (
	"context"
	"fmt"
	"time"
	"vpn-bot-api/config"

	"golang.org/x/crypto/ssh"
)

type Ssh struct {
	config  *ssh.ClientConfig
	address string
	port    string
	network string
}

func NewSsh(config *config.Config) (*Ssh, error) {
	clientConfig := &ssh.ClientConfig{
		User: config.Server.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Server.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	return &Ssh{config: clientConfig, address: config.Server.Address, port: config.Server.Port, network: config.Server.Network}, nil
}

func (s *Ssh) Execute(_ context.Context, cmds []string) ([]string, error) {
	address := fmt.Sprintf("%s:%s", s.address, s.port)
	client, err := ssh.Dial(s.network, address, s.config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}
	defer client.Close()

	results := make([]string, 0, len(cmds))

	for _, cmd := range cmds {
		session, err := client.NewSession()
		if err != nil {
			return nil, fmt.Errorf("failed to create session: %w", err)
		}

		output, err := session.CombinedOutput(cmd)
		session.Close()

		if err != nil {
			return nil, fmt.Errorf("failed to run command %q: %w", cmd, err)
		}

		results = append(results, string(output))
	}

	return results, nil
}
