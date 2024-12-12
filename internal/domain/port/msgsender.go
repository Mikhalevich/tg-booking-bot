package port

import (
	"context"
)

type Button struct {
	Text string
	Data string
}

type MessageSender interface {
	ReplyText(ctx context.Context, chatID ChatID, replyToMsgID MessageID, text string, buttons ...Button)
	ReplyTextMarkdown(ctx context.Context, chatID ChatID, replyToMsgID MessageID, text string)
	EscapeMarkdown(s string) string
}
