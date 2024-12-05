package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/role"
)

func (p *Postgres) IsEmployeeWithRoleExists(
	ctx context.Context,
	role role.Role,
) (bool, error) {
	var isExists bool
	if err := sqlx.GetContext(
		ctx,
		p.db,
		&isExists,
		`
			SELECT
				TRUE
			FROM
				employee INNER JOIN role ON employee.role_id = role.id
			WHERE
				role.name = $1
		`,
		role.String(),
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("get context: %w", err)
	}

	return true, nil
}
