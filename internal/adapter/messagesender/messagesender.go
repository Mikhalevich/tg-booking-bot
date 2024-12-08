package messagesender

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/infra/logger"
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

func (m *messageSender) ReplyText(
	ctx context.Context,
	chatID int64,
	replyToMsgID int,
	text string,
	buttons ...port.Button,
) {
	if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		ReplyParameters: &models.ReplyParameters{
			MessageID: replyToMsgID,
		},
		Text:        text,
		ReplyMarkup: makeButtonsMarkup(buttons...),
	}); err != nil {
		logger.FromContext(ctx).
			WithError(err).
			WithField("text_plain", text).
			Error("failed to reply text")
	}
}

func makeButtonsMarkup(buttons ...port.Button) models.ReplyMarkup {
	if len(buttons) == 0 {
		return nil
	}

	buttonRow := make([]models.InlineKeyboardButton, 0, len(buttons))
	for _, b := range buttons {
		buttonRow = append(buttonRow, models.InlineKeyboardButton{
			Text:         b.Text,
			CallbackData: b.Data,
		})
	}

	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			buttonRow,
		},
	}
}

func (m *messageSender) ReplyTextMarkdown(ctx context.Context, chatID int64, replyToMsgID int, text string) {
	if _, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		ReplyParameters: &models.ReplyParameters{
			MessageID: replyToMsgID,
		},
		ParseMode: models.ParseModeMarkdown,
		Text:      text,
	}); err != nil {
		logger.FromContext(ctx).
			WithError(err).
			WithField("text_markdown", text).
			Error("failed to reply text markdown")
	}
}

func (m *messageSender) EscapeMarkdown(s string) string {
	return bot.EscapeMarkdown(s)
}
