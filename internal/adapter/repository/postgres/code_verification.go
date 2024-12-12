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

func (p *Postgres) CodeVerification(ctx context.Context, code string, chatID msginfo.ChatID) (*empl.Employee, error) {
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
			"state_registered":            empl.EmployeeStateRegistered.String(),
			"state_verification_required": empl.EmployeeStateVerificationRequired.String(),
		})
	if err != nil {
		return nil, fmt.Errorf("named: %w", err)
	}

	var empl model.Employee
	if err := sqlx.GetContext(ctx, p.db, &empl, p.db.Rebind(query), args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotFound
		}

		return nil, fmt.Errorf("get employee: %w", err)
	}

	portEmpl := model.ToPortEmployee(empl)

	return &portEmpl, nil
}
