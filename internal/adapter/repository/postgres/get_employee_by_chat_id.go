package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/empl"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/msginfo"
)

func (p *Postgres) GetEmployeeByChatID(ctx context.Context, chatID msginfo.ChatID) (*empl.Employee, error) {
	var emp model.Employee
	if err := sqlx.GetContext(
		ctx,
		p.db,
		&emp,
		`SELECT
			employee.id,
			employee.first_name,
			employee.last_name,
			role.name AS role_name,
			employee.chat_id,
			employee.state,
			employee.verification_code,
			employee.created_at,
			employee.updated_at
		FROM
			employee INNER JOIN role on employee.role_id = role.id
		WHERE
			employee.chat_id = $1
	`, chatID.Int64()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotFound
		}

		return nil, fmt.Errorf("select employee: %w", err)
	}

	return model.ToPortEmployee(emp), nil
}
