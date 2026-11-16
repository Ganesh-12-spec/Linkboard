package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Ganesh-12-spec/Linkboard/internal/database"
	"github.com/Ganesh-12-spec/Linkboard/internal/repository"
)

func TestLoginIntegration(t *testing.T) {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		t.Fatal("DATABASE_URL is required for integration testing")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	pool, err := database.NewPool(ctx, connString)
	if err != nil {
		t.Fatalf("database connection failed: %v", err)
	}
	t.Cleanup(func() {
		pool.Close()
	})

	repo := repository.NewUserRepository(pool)
	userService := NewUserService(repo)

	email := fmt.Sprintf("test-%d@example.com", time.Now().UnixNano())
	password := "test-password-123"

	// Create a temporary user.
	createdUser, err := userService.CreateUser(
		ctx,
		"Test User",
		email,
		password,
	)
	if err != nil {
		t.Fatalf("user creation failed: %v", err)
	}

	// Delete the temporary user after the test.
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cleanupCancel()

		_, err := pool.Exec(
			cleanupCtx,
			"DELETE FROM users WHERE email = $1",
			email,
		)
		if err != nil {
			t.Errorf("temporary user cleanup failed: %v", err)
		}
	})

	// Test 1: Correct password.
	loggedInUser, err := userService.Login(ctx, email, password)
	if err != nil {
		t.Fatalf("login with correct password failed: %v", err)
	}
	if loggedInUser.ID != createdUser.ID {
		t.Fatal("login returned an unexpected user")
	}

	// Test 2: Incorrect password.
	_, err = userService.Login(ctx, email, "wrong-password")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got: %v", err)
	}

	// Test 3: Email does not exist.
	_, err = userService.Login(
		ctx,
		"missing@example.com",
		password,
	)
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got: %v", err)
	}
}
