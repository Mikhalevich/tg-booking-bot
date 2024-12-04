package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/role"
)

var (
	_ port.ScheduleRepository = (*Postgres)(nil)
	_ port.EmployeeRepository = (*Postgres)(nil)
)

type Postgres struct {
	db sqlx.ExtContext
}

func New(db *sql.DB, driverName string) *Postgres {
	return &Postgres{
		db: sqlx.NewDb(db, driverName),
	}
}

func (p *Postgres) roleIDByName(ctx context.Context, r role.Role, db sqlx.QueryerContext) (int, error) {
	var roleID int
	if err := sqlx.GetContext(
		ctx,
		db,
		&roleID,
		`
			SELECT id
			FROM role
			WHERE name = $1
		`,
		r.String(),
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("role name %s doesn't exist", r.String())
		}

		return 0, fmt.Errorf("get role_id by name: %w", err)
	}

	return roleID, nil
}
