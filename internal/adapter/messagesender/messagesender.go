package messagesender

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/paginator"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

var _ port.MessageSender = (*messageSender)(nil)

type messageSender struct {
	bot *bot.Bot
}

func New(bot *bot.Bot) *messageSender {
	return &messageSender{
		bot: bot,
	}
}

func (m *messageSender) ReplyText(ctx context.Context, chatID int64, replyToMsgID int, text string) error {
	if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		ReplyParameters: &models.ReplyParameters{
			MessageID: replyToMsgID,
		},
		Text: text,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (m *messageSender) ReplyTextMarkdown(ctx context.Context, chatID int64, replyToMsgID int, text string) error {
	if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		ReplyParameters: &models.ReplyParameters{
			MessageID: replyToMsgID,
		},
		ParseMode: models.ParseModeMarkdown,
		Text:      text,
	}); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (m *messageSender) EscapeMarkdown(s string) string {
	return bot.EscapeMarkdown(s)
}

func (m *messageSender) SendPaginator(ctx context.Context, chatID int64, data []string) error {
	opts := []paginator.Option{
		paginator.PerPage(1),
		paginator.WithCloseButton("Close"),
	}

	p := paginator.New(m.bot, data, opts...)

	if _, err := p.Show(ctx, m.bot, strconv.Itoa(int(chatID))); err != nil {
		return fmt.Errorf("show paginator: %w", err)
	}

	return nil
}
