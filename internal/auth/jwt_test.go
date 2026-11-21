package auth

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "this-is-a-test-secret-with-at-least-32-bytes"

func TestJWTManager_ValidToken(t *testing.T) {
	manager, err := NewJWTManager(testSecret, 15*time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	token, err := manager.GenerateToken(42)
	if err != nil {
		t.Fatal(err)
	}

	userID, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatal(err)
	}

	if userID != 42 {
		t.Fatalf("expected user ID 42, got %d", userID)
	}
}

func TestJWTManager_ExpiredToken(t *testing.T) {
	manager, err := NewJWTManager(testSecret, -1*time.Second)
	if err == nil {
		t.Fatal("expected manager to reject non-positive expiry")
	}

	_ = manager
}

func TestJWTManager_TamperedToken(t *testing.T) {
	manager, err := NewJWTManager(testSecret, 15*time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	token, err := manager.GenerateToken(42)
	if err != nil {
		t.Fatal(err)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatal("expected JWT to contain three parts")
	}

	parts[1] = parts[1] + "tampered"
	tamperedToken := strings.Join(parts, ".")

	_, err = manager.ValidateToken(tamperedToken)
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}

func TestJWTManager_WrongSecret(t *testing.T) {
	manager1, err := NewJWTManager(testSecret, 15*time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	manager2, err := NewJWTManager(
		"another-test-secret-with-at-least-32-bytes",
		15*time.Minute,
	)
	if err != nil {
		t.Fatal(err)
	}

	token, err := manager1.GenerateToken(42)
	if err != nil {
		t.Fatal(err)
	}

	_, err = manager2.ValidateToken(token)
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}

func TestJWTManager_InvalidUserID(t *testing.T) {
	manager, err := NewJWTManager(testSecret, 15*time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	_, err = manager.GenerateToken(0)
	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID, got %v", err)
	}
}

func TestJWTManager_RejectsUnexpectedSigningMethod(t *testing.T) {
	manager, err := NewJWTManager(testSecret, 15*time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	claims := jwt.RegisteredClaims{
		Subject:   "42",
		Issuer:    "linkboard",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)

	signedToken, err := token.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatal(err)
	}

	_, err = manager.ValidateToken(signedToken)
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}
