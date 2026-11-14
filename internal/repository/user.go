package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID           int64
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return UserRepository{db: db}
}
func (r UserRepository) CreateUser(
	ctx context.Context,
	name, email, passwordHash string,
) (User, error) {
	var user User

	err := r.db.QueryRow(
		ctx,
		`INSERT INTO users (name, email, password_hash)
         VALUES ($1, $2, $3)
         RETURNING id, name, email, created_at`,
		name,
		email,
		passwordHash,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}
