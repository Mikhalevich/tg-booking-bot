package port

import (
	"context"
)

type MessageInfo struct {
	MessageID int
	ChatID    int64
	// for text message
	Text string
	// for callback query
	Data string
}

type Handler func(ctx context.Context, info MessageInfo) error
