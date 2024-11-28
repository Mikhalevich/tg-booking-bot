package port

import (
	"context"
)

type MessageInfo struct {
	MessageID int
	ChatID    int64
	UserID    int64
	Text      string
}

type Handler func(ctx context.Context, info MessageInfo) error
