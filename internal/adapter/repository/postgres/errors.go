package postgres

import (
	"errors"
)

var (
	errNotFound = errors.New("not found")
)

func (p *Postgres) IsNotFoundError(err error) bool {
	return errors.Is(err, errNotFound)
}
