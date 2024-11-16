package port

import (
	"context"
)

type MessageSender interface {
	Reply(ctx context.Context, chatID int64, replyToMsgID int, text string) error
}
