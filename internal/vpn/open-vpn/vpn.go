package openvpn

import (
	"bytes"
	"context"
	"fmt"
	"text/template"
	"vpn-bot-api/config"
)

type Vpn struct {
	config   *config.Config
	executor Executor
}

type OVPNData struct {
	RemoteIP   string
	RemotePort int
	Auth       string
	Cipher     string
	CA         string
	Cert       string
	Key        string
	TLSKey     string
}

type Executor interface {
	Execute(ctx context.Context, cmds []string) ([]string, error)
}

func NewVpn(config *config.Config, executor Executor) (*Vpn, error) {
	if !config.Vpn.OpenVpn.Enabled {
		return nil, nil
	}

	return &Vpn{
		config:   config,
		executor: executor,
	}, nil
}

func (s *Vpn) AddVPNClient(ctx context.Context, name string) (string, error) {
	_, err := s.executor.Execute(ctx, []string{
		fmt.Sprintf("cd /etc/openvpn/server/easy-rsa && EASYRSA_BATCH=1 ./easyrsa build-client-full %s nopass", name),
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate client cert: %w", err)
	}

	files := map[string]string{
		"ca":     "/etc/openvpn/server/easy-rsa/pki/ca.crt",
		"cert":   "/etc/openvpn/server/easy-rsa/pki/issued/" + name + ".crt",
		"key":    "/etc/openvpn/server/easy-rsa/pki/private/" + name + ".key",
		"tlskey": "/etc/openvpn/server/tc.key",
	}

	content := make(map[string]string)
	for k, path := range files {
		out, err := s.executor.Execute(ctx, []string{"cat " + path})
		if err != nil {
			return "", fmt.Errorf("failed to read %s: %w", path, err)
		}
		if len(out) == 0 {
			return "", fmt.Errorf("empty file: %s", path)
		}
		content[k] = out[0]
	}

	data := OVPNData{
		RemoteIP:   s.config.Server.Address,
		RemotePort: s.config.Vpn.OpenVpn.RemotePort,
		Auth:       s.config.Vpn.OpenVpn.Auth,
		Cipher:     s.config.Vpn.OpenVpn.Cipher,
		CA:         content["ca"],
		Cert:       content["cert"],
		Key:        content["key"],
		TLSKey:     content["tlskey"],
	}

	tmpl, err := template.New("ovpn").Parse(s.config.Vpn.OpenVpn.Template)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return buf.String(), nil
}
