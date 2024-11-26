package postgres

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

var (
	_ port.ScheduleRepository = (*Postgres)(nil)
	_ port.EmployeeRepository = (*Postgres)(nil)
)

type Postgres struct {
	db *sqlx.DB
}

func New(db *sql.DB, driverName string) *Postgres {
	return &Postgres{
		db: sqlx.NewDb(db, driverName),
	}
}
