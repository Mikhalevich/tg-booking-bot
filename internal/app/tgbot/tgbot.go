package tgbot

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type tgbot struct {
}

func Start(ctx context.Context, token string) error {
	tbot := &tgbot{}

	opts := []bot.Option{
		bot.WithDefaultHandler(tbot.handler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		return fmt.Errorf("creating bot with options: %w", err)
	}

	b.Start(ctx)

	return nil
}

func (t *tgbot) handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	//nolint:errcheck
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
