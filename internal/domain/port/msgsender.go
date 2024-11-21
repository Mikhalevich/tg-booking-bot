package port

import (
	"context"
)

type MessageSender interface {
	ReplyText(ctx context.Context, chatID int64, replyToMsgID int, text string) error
	ReplyTextMarkdown(ctx context.Context, chatID int64, replyToMsgID int, text string) error
	EscapeMarkdown(s string) string
}
