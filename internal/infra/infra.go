package infra

import (
	"context"
	"flag"
	"fmt"

	"github.com/go-telegram/bot"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jinzhu/configor"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"

	"github.com/Mikhalevich/tg-booking-bot/internal/adapter/messagesender"
	"github.com/Mikhalevich/tg-booking-bot/internal/adapter/repository/postgres"
	"github.com/Mikhalevich/tg-booking-bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-booking-bot/internal/config"
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

func MakeScheduler(
	b *bot.Bot,
	postgresCfg config.Postgres,
) (tgbot.Scheduler, func(), error) {
	pg, cleanup, err := MakePostgres(postgresCfg)
	if err != nil {
		return nil, nil, fmt.Errorf("make postgres: %w", err)
	}

	return schedule.New(
		pg,
		messagesender.New(b),
	), cleanup, nil
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

func MakePostgres(cfg config.Postgres) (*postgres.Postgres, func(), error) {
	if cfg.Connection == "" {
		return nil, func() {}, nil
	}

	db, err := otelsql.Open("pgx", cfg.Connection)
	if err != nil {
		return nil, nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, nil, fmt.Errorf("ping: %w", err)
	}

	p := postgres.New(db, "pgx")

	return p, func() {
		db.Close()
	}, nil
}
