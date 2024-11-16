package infra

import (
	"context"
	"flag"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/jinzhu/configor"

	"github.com/Mikhalevich/tg-booking-bot/internal/adapter/messagesender"
	"github.com/Mikhalevich/tg-booking-bot/internal/adapter/repository"
	"github.com/Mikhalevich/tg-booking-bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/schedule"
	"github.com/Mikhalevich/tg-booking-bot/internal/infra/logger"
)

func LoadConfig(cfg any) error {
	configFile := flag.String("config", "config/config.yaml", "consumer worker config file")
	flag.Parse()

	if err := configor.Load(cfg, *configFile); err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	return nil
}

func SetupLogger(lvl string) (logger.Logger, error) {
	log, err := logger.NewLogrusWithLevel(lvl)
	if err != nil {
		return nil, fmt.Errorf("creating new logger: %w", err)
	}

	logger.SetStdLogger(log)

	return log, nil
}

func MakeBotAPI(token string) (*bot.Bot, error) {
	b, err := bot.New(token)
	if err != nil {
		return nil, fmt.Errorf("creating bot: %w", err)
	}

	return b, nil
}

func MakeScheduler(b *bot.Bot) tgbot.Scheduler {
	return schedule.New(
		repository.NewNoop(),
		messagesender.New(b),
	)
}

func StartScheduleBot(
	ctx context.Context,
	b *bot.Bot,
	logger logger.Logger,
	scheduler tgbot.Scheduler,
) error {
	if err := tgbot.Start(
		ctx,
		b,
		logger,
		tgbot.ScheduleRoutes(scheduler),
	); err != nil {
		return fmt.Errorf("start bot: %w", err)
	}

	return nil
}
