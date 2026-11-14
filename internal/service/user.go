package service

import (
	"context"

	"github.com/Ganesh-12-spec/Linkboard/internal/repository"
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
