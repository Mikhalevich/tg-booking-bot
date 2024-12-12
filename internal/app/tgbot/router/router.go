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
	AddDefaultCallbackQueryHander(h port.Handler)
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
		r.wrapHandler(pattern, handler),
	)
}

func (r *Router) AddDefaultTextHandler(h port.Handler) {
	r.bot.RegisterHandler(
		bot.HandlerTypeMessageText,
		"",
		bot.MatchTypePrefix,
		r.wrapHandler("default_text_handler", h),
	)
}

func (r *Router) AddDefaultCallbackQueryHander(h port.Handler) {
	r.bot.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"",
		bot.MatchTypePrefix,
		r.wrapHandler("default_callback_query", h),
	)
}

func (r *Router) AddMiddleware(m Middleware) {
	r.middlewares = append(r.middlewares, m)
}

func (r *Router) wrapHandler(pattern string, h port.Handler) bot.HandlerFunc {
	h = r.applyMiddleware(h)

	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		ctx, span := tracing.StartSpanName(ctx, pattern)
		defer span.End()

		var (
			msgInfo = makeMsgInfoFromUpdate(update)
			log     = r.logger.WithContext(ctx).WithField("endpoint", pattern)
			ctxLog  = logger.WithLogger(ctx, log)
		)

		if err := h(ctxLog, msgInfo); err != nil {
			log.WithError(err).Error("error while processing message")

			if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: msgInfo.ChatID,
				ReplyParameters: &models.ReplyParameters{
					MessageID: msgInfo.MessageID.Int(),
				},
				Text: "internal error",
			}); err != nil {
				log.WithError(err).Error("send internal error message")
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

func makeMsgInfoFromUpdate(u *models.Update) port.MessageInfo {
	if u.Message != nil {
		return port.MessageInfo{
			MessageID: port.MessageIDFromInt(u.Message.ID),
			ChatID:    port.ChatIDFromInt(u.Message.Chat.ID),
			Text:      u.Message.Text,
		}
	}

	if u.CallbackQuery != nil {
		if u.CallbackQuery.Message.Message != nil {
			return port.MessageInfo{
				MessageID: port.MessageIDFromInt(u.CallbackQuery.Message.Message.ID),
				ChatID:    port.ChatIDFromInt(u.CallbackQuery.Message.Message.Chat.ID),
				Data:      u.CallbackQuery.Data,
			}
		}

		if u.CallbackQuery.Message.InaccessibleMessage != nil {
			return port.MessageInfo{
				MessageID: port.MessageIDFromInt(u.CallbackQuery.Message.InaccessibleMessage.MessageID),
				ChatID:    port.ChatIDFromInt(u.CallbackQuery.Message.InaccessibleMessage.Chat.ID),
				Data:      u.CallbackQuery.Data,
			}
		}
	}

	return port.MessageInfo{}
}
