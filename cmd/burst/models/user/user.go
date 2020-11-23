package user

import (
	"context"
	"database/sql"
)

type User struct {
	ID     int64  `json:"id,omitempty"`
	Online bool   `json:"online,omitempty"`
	Seen   string `json:"seen,omitempty"`
}

type dbIterator interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func Create(ctx context.Context, db dbIterator, u *User) error {
	stmt := `
		INSERT
		INTO public."users" (
			id,
			online,
			seen
		)
		VALUES (
			$1,
			$2,
			$3
		)
		ON CONFLICT (id) 
		DO
			UPDATE SET online = $2, seen = $3
	`

	_, err := db.ExecContext(ctx, stmt, u.ID, u.Online, u.Seen)
	return err
}

func RemoveOlderThan(ctx context.Context, db dbIterator, date string) error {
	stmt := `
		DELETE
		FROM public."users"
		WHERE seen < $1
	`

	_, err := db.ExecContext(ctx, stmt, date)
	return err
}
