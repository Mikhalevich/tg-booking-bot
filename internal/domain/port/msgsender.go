package port

import (
	"context"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/msginfo"
)

type Button struct {
	Text string
	Data string
}

type MessageSender interface {
	ReplyText(ctx context.Context, chatID msginfo.ChatID, replyToMsgID msginfo.MessageID, text string, buttons ...Button)
	ReplyTextMarkdown(ctx context.Context, chatID msginfo.ChatID, replyToMsgID msginfo.MessageID, text string)
	EscapeMarkdown(s string) string
}
