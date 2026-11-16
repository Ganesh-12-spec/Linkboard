package service

import (
	"context"
	"errors"

	"github.com/Ganesh-12-spec/Linkboard/internal/repository"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return UserService{repo: repo}
}

func (s UserService) CreateUser(
	ctx context.Context,
	name, email, password string,
) (repository.User, error) {
	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return repository.User{}, err
	}

	// Save the user with the hashed password
	user, err := s.repo.CreateUser(
		ctx,
		name,
		email,
		string(hashedPassword),
	)
	if err != nil {
		return repository.User{}, err
	}

	return user, nil
}

var ErrInvalidCredentials = errors.New("invalid email or password")

func (s UserService) Login(
	ctx context.Context,
	email, password string,
) (repository.User, error) {

	// Find the user using email
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.User{}, ErrInvalidCredentials
		}

		return repository.User{}, err
	}

	// Verify the entered password against the stored hash
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return repository.User{}, ErrInvalidCredentials
		}

		return repository.User{}, err
	}

	return user, nil
}
