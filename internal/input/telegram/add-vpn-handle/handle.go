package addvpnhandle

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	openvpn "vpn-bot-api/internal/vpn/open-vpn"
)

type Handle struct {
	vpn *openvpn.Vpn
}

type VpnService interface {
	AddVPNClient(ctx context.Context, name string) (string, error)
}

func NewHandle(vpn *openvpn.Vpn) *Handle {
	return &Handle{vpn: vpn}
}

func (h *Handle) Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := update.Message.Text
	name := strings.TrimSpace(strings.TrimPrefix(text, "/add"))
	if name == "" {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "⚠️Укажи имя: /add user1",
		})
		if err != nil {
			log.Println(err)
		}

		return
	}

	cfg, err := h.vpn.AddVPNClient(ctx, name)
	if err != nil {
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("❌ Ошибка: %v", err),
		})

		return
	}

	doc := &models.InputFileUpload{
		Filename: fmt.Sprintf("%s.ovpn", name),
		Data:     bytes.NewReader([]byte(cfg)),
	}

	_, err = b.SendDocument(ctx, &bot.SendDocumentParams{
		ChatID:    update.Message.Chat.ID,
		Document:  doc,
		Caption:   fmt.Sprintf("✅ Профиль *%s* сгенерирован", name),
		ParseMode: models.ParseModeHTML,
	})

	if err != nil {
		log.Println(err)
	}
}
