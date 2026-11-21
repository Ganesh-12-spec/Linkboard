package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrInvalidUserID    = errors.New("invalid user ID")
	ErrInvalidJWTConfig = errors.New("invalid JWT configuration")
	ErrUnexpectedMethod = errors.New("unexpected signing method")
)

type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTManager(secret string, ttl time.Duration) (*JWTManager, error) {
	if len(secret) < 32 {
		return nil, fmt.Errorf("%w: secret must be at least 32 bytes", ErrInvalidJWTConfig)
	}

	if ttl <= 0 {
		return nil, fmt.Errorf("%w: token expiry must be positive", ErrInvalidJWTConfig)
	}

	return &JWTManager{
		secret: []byte(secret),
		ttl:    ttl,
	}, nil
}

func (m *JWTManager) GenerateToken(userID int64) (string, error) {
	if userID <= 0 {
		return "", ErrInvalidUserID
	}

	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		Issuer:    "linkboard",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("sign JWT: %w", err)
	}

	return signedToken, nil
}

func (m *JWTManager) ValidateToken(tokenString string) (int64, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, ErrUnexpectedMethod
			}

			return m.secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("%w: verification failed", ErrInvalidToken)
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil || userID <= 0 {
		return 0, ErrInvalidUserID
	}

	return userID, nil
}
