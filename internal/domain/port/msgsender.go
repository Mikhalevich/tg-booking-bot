package port

import (
	"context"
)

type Button struct {
	Text string
	Data string
}

type MessageSender interface {
	ReplyText(ctx context.Context, chatID int64, replyToMsgID int, text string, buttons ...Button) error
	ReplyTextMarkdown(ctx context.Context, chatID int64, replyToMsgID int, text string) error
	EscapeMarkdown(s string) string
}
