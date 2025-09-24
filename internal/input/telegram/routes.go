package telegram

import (
	"context"
	
	"github.com/go-telegram/bot"
	addvpnhandle "vpn-bot-api/internal/input/telegram/add-vpn-handle"
	openvpn "vpn-bot-api/internal/vpn/open-vpn"
)

const (
	MTAddVpn = "/add"
)

func RegisterHandlers(_ context.Context, b *bot.Bot, openVpn *openvpn.Vpn) error {
	addVpnHandle := addvpnhandle.NewHandle(openVpn)

	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		MTAddVpn,
		bot.MatchTypePrefix,
		addVpnHandle.Handler,
	)

	return nil
}
