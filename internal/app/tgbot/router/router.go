package router

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/infra/logger"
	"github.com/Mikhalevich/tg-booking-bot/internal/infra/tracing"
)

type Middleware func(next port.Handler) port.Handler

type Register interface {
	AddExactTextRoute(pattern string, handler port.Handler)
	AddDefaultTextHandler(handler port.Handler)
	AddMiddleware(middleware Middleware)
	MiddlewareGroup(fn func(r Register))
}

type Router struct {
	bot         *bot.Bot
	logger      logger.Logger
	middlewares []Middleware
}

func New(b *bot.Bot, logger logger.Logger) *Router {
	return &Router{
		bot:    b,
		logger: logger,
	}
}

func (r *Router) MiddlewareGroup(fn func(r Register)) {
	group := &Router{
		bot:         r.bot,
		logger:      r.logger,
		middlewares: r.middlewares[:len(r.middlewares):len(r.middlewares)],
	}

	fn(group)
}

func (r *Router) AddExactTextRoute(pattern string, handler port.Handler) {
	r.bot.RegisterHandler(
		bot.HandlerTypeMessageText,
		pattern,
		bot.MatchTypeExact,
		r.wrapTextHandler(pattern, handler),
	)
}

func (r *Router) AddDefaultTextHandler(h port.Handler) {
	r.bot.RegisterHandler(
		bot.HandlerTypeMessageText,
		"",
		bot.MatchTypePrefix,
		r.wrapTextHandler("default_handler", h),
	)
}

func (r *Router) AddMiddleware(m Middleware) {
	r.middlewares = append(r.middlewares, m)
}

func (r *Router) wrapTextHandler(pattern string, h port.Handler) bot.HandlerFunc {
	h = r.applyMiddleware(h)

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
			r.logger.WithContext(ctx).
				WithError(err).
				WithField("endpoint", pattern).
				Error("error while processing message")

			if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				ReplyParameters: &models.ReplyParameters{
					MessageID: update.Message.ID,
				},
				Text: "internal error",
			}); err != nil {
				r.logger.WithContext(ctx).
					WithError(err).
					WithField("endpoint", pattern).
					Error("send internal error message")
			}
		}
	}
}

func (r *Router) applyMiddleware(h port.Handler) port.Handler {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		h = r.middlewares[i](h)
	}

	return h
}
