package ctxdata

import (
	"context"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

type ctxdatakey int

const (
	employeeKey ctxdatakey = 1
)

func WithEmployee(ctx context.Context, employee port.Employee) context.Context {
	return context.WithValue(ctx, employeeKey, employee)
}

func Employee(ctx context.Context) (port.Employee, bool) {
	val := ctx.Value(employeeKey)
	if val == nil {
		return port.Employee{}, false
	}

	employee, ok := val.(port.Employee)
	if !ok {
		return port.Employee{}, false
	}

	return employee, true
}
