package tgbot

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/infra/logger"
	"github.com/Mikhalevich/tg-booking-bot/internal/infra/tracing"
)

type tgbot struct {
	bot         *bot.Bot
	logger      logger.Logger
	middlewares []Middleware
}

type Middleware func(next port.Handler) port.Handler

type Register interface {
	AddExactTextRoute(pattern string, handler port.Handler)
	AddDefaultTextHandler(handler port.Handler, middlewares ...Middleware)
	AddMiddleware(middleware Middleware)
}

type RouteRegisterFunc func(register Register)

func Start(
	ctx context.Context,
	b *bot.Bot,
	logger logger.Logger,
	routesRegisterFn RouteRegisterFunc,
) error {
	tbot := &tgbot{
		bot:    b,
		logger: logger,
	}

	routesRegisterFn(tbot)

	b.Start(ctx)

	return nil
}

func (t *tgbot) AddExactTextRoute(pattern string, handler port.Handler) {
	t.bot.RegisterHandler(
		bot.HandlerTypeMessageText,
		pattern,
		bot.MatchTypeExact,
		t.wrapTextHandler(pattern, handler, t.middlewares...),
	)
}

func (t *tgbot) AddDefaultTextHandler(h port.Handler, middlewares ...Middleware) {
	t.bot.RegisterHandler(
		bot.HandlerTypeMessageText,
		"",
		bot.MatchTypePrefix,
		t.wrapTextHandler("default_handler", h, middlewares...),
	)
}

func (t *tgbot) AddMiddleware(m Middleware) {
	t.middlewares = append(t.middlewares, m)
}

func (t *tgbot) wrapTextHandler(pattern string, h port.Handler, middlewares ...Middleware) bot.HandlerFunc {
	h = applyMiddleware(h, middlewares...)

	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		ctx, span := tracing.StartSpanName(ctx, pattern)
		defer span.End()

		if err := h(
			ctx,
			port.MessageInfo{
				MessageID: update.Message.ID,
				ChatID:    update.Message.Chat.ID,
				Text:      update.Message.Text,
			},
		); err != nil {
			t.logger.WithError(err).
				WithField("endpoint", pattern).
				Error("error while processing message")
		}
	}
}

func applyMiddleware(h port.Handler, middlewares ...Middleware) port.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}

	return h
}
