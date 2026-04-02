// jwt configs code

package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret      []byte
	exipryHours int
}

func NewManager(secret string, expiryHours int) *Manager {
	return &Manager{
		secret:      []byte(secret),
		exipryHours: expiryHours,
	}
}

// function to generate token

func (m *Manager) GenerateToken(userID, email, role string) (string, error) {

	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(m.exipryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "finance-dashboard",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)

	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signed, nil

}

// function to parse the JWT

func (m *Manager) ParseToken(tokenStr string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])

		}

		return m.secret, nil

	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)

	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil

}
