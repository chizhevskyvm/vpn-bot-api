package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-telegram/bot"
	"gopkg.in/yaml.v3"
	"vpn-bot-api/config"
	"vpn-bot-api/internal/input/telegram"
	"vpn-bot-api/internal/ssh"
	openvpn "vpn-bot-api/internal/vpn/open-vpn"
)

const configPath = "./conf.yaml"

func RunBot(ctx context.Context) error {
	cfg, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	sshClient, err := ssh.NewSsh(cfg)
	if err != nil {
		return fmt.Errorf("failed to init ssh client: %w", err)
	}

	openVpn, err := openvpn.NewVpn(cfg, sshClient)
	if err != nil {
		return fmt.Errorf("failed to init openvpn: %w", err)
	}

	b, err := bot.New(cfg.Telegram.Token)
	if err != nil {
		return fmt.Errorf("failed to create telegram bot: %w", err)
	}

	if err := telegram.RegisterHandlers(ctx, b, openVpn); err != nil {
		return fmt.Errorf("failed to register telegram handlers: %w", err)
	}

	log.Println("🤖 Bot is starting...")
	b.Start(ctx)
	log.Println("✅ Bot stopped")

	return nil
}

func LoadConfig() (*config.Config, error) {
	data, _ := os.ReadFile(configPath)

	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)

		return nil, err
	}

	return &cfg, nil
}
