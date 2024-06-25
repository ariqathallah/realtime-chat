package user

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	var lastInserId int
	query := "INSERT INTO users(username, email, password) VALUES ($1, $2, $3) returning id"
	if err := r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(&lastInserId); err != nil {
		return &User{}, err
	}

	user.ID = int64(lastInserId)
	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u := User{}

	query := "SELECT id, username, email, password FROM USERS WHERE email = $1"
	if err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Username, &u.Email, &u.Password); err != nil {
		return &User{}, err
	}

	return &u, nil
}
