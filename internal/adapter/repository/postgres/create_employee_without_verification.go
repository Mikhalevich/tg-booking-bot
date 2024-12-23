package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/empl"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/role"
)

func (p *Postgres) CreateEmployeeWithoutVerification(
	ctx context.Context,
	r role.Role,
	chatID msginfo.ChatID,
) (empl.EmployeeID, error) {
	roleID, err := roleIDByName(ctx, r, p.db)
	if err != nil {
		return 0, fmt.Errorf("get role by name: %w", err)
	}

	query, args, err := sqlx.Named(
		`INSERT INTO employee(
				role_id,
				state,
				chat_id
			) VALUES (
				:role_id,
				:state,
				:chat_id
			)
			RETURNING id
			`, map[string]any{
			"role_id": roleID,
			"state":   empl.EmployeeStateRegistered,
			"chat_id": chatID.Int64(),
		})

	if err != nil {
		return 0, fmt.Errorf("create named query: %w", err)
	}

	var employeeID int
	if err := sqlx.GetContext(ctx, p.db, &employeeID, p.db.Rebind(query), args...); err != nil {
		return 0, fmt.Errorf("insert employee: %w", err)
	}

	return empl.EmployeeIDFromInt(employeeID), nil
}
