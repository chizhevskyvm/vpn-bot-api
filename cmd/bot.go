package cmd

import (
	"context"
	"github.com/go-telegram/bot"
	"log"
	"os"
	"vpn-bot-api/internal/input/telegram"

	"gopkg.in/yaml.v3"
	"vpn-bot-api/config"
	"vpn-bot-api/internal/ssh"
	openvpn "vpn-bot-api/internal/vpn/open-vpn"
)

const configPath = "./conf.yaml"

func RunBot(ctx context.Context) error {
	cfg := GetConfig()

	sshClient, err := ssh.NewSsh(cfg)
	if err != nil {
		return err
	}

	openVpn, err := openvpn.NewVpn(cfg, sshClient)
	if err != nil {
		return err
	}

	b, err := bot.New(cfg.Telegram.Token)
	if err != nil {
		log.Fatalf("Error create telegram bot: %v", err)
	}

	err = telegram.RegisterHandlers(ctx, b, openVpn)
	if err != nil {
		log.Fatalf("Error register telegram handler: %v", err)
	}

	go b.Start(ctx)

	return nil
}

func GetConfig() *config.Config {
	data, _ := os.ReadFile(configPath)

	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}

	return &cfg
}
