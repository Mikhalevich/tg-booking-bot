package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

func (p *Postgres) CodeVerification(ctx context.Context, code string, chatID port.ChatID) (*port.Employee, error) {
	query, args, err := sqlx.Named(
		`
			UPDATE employee SET
				chat_id = :chat_id,
				state = :state_registered
			WHERE
				verification_code = :verification_code AND
				state = :state_verification_required
			RETURNING *
		`, map[string]any{
			"chat_id":                     chatID.Int64(),
			"verification_code":           code,
			"state_registered":            port.EmployeeStateRegistered.String(),
			"state_verification_required": port.EmployeeStateVerificationRequired.String(),
		})
	if err != nil {
		return nil, fmt.Errorf("named: %w", err)
	}

	var empl employee
	if err := sqlx.GetContext(ctx, p.db, &empl, p.db.Rebind(query), args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotFound
		}

		return nil, fmt.Errorf("get employee: %w", err)
	}

	portEmpl := convertToEmployee(empl)

	return &portEmpl, nil
}
