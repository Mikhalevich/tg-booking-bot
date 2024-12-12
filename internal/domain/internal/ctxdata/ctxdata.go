package ctxdata

import (
	"context"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/empl"
)

type ctxdatakey int

const (
	employeeKey ctxdatakey = 1
)

func WithEmployee(ctx context.Context, employee empl.Employee) context.Context {
	return context.WithValue(ctx, employeeKey, employee)
}

func Employee(ctx context.Context) (empl.Employee, bool) {
	val := ctx.Value(employeeKey)
	if val == nil {
		return empl.Employee{}, false
	}

	employee, ok := val.(empl.Employee)
	if !ok {
		return empl.Employee{}, false
	}

	return employee, true
}
