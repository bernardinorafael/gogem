package crypto

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bernardinorafael/gogem/uid"
	"github.com/golang-jwt/jwt/v5"
)

const (
	minSecretKeyLength = 32
	tokenIssuer        = "token-service"
)

func GenerateToken(secretKey, userID, sessionID string, orgID *string, duration time.Duration) (string, *TokenClaims, error) {
	if len(secretKey) != minSecretKeyLength {
		return "", nil, errors.New("invalid secret key length")
	}

	claims := newTokenClaims(userID, sessionID, orgID, duration)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", claims, fmt.Errorf("token signing failed: %w", err)
	}

	return token, claims, nil
}

func VerifyToken(secretKey string, v string) (*TokenClaims, error) {
	if strings.TrimSpace(v) == "" {
		return nil, errors.New("token is empty")
	}

	keyFunc := func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secretKey), nil
	}

	token, err := jwt.ParseWithClaims(v, &TokenClaims{}, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("token parsing failed: %w", err)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

type TokenClaims struct {
	UserID    string  `json:"userId"`
	OrgID     *string `json:"orgId"`
	SessionID string  `json:"sessionId"`
	jwt.RegisteredClaims
}

func newTokenClaims(userID, sessionID string, orgID *string, duration time.Duration) *TokenClaims {
	return &TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uid.New("tok"),
			Issuer:    tokenIssuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
		UserID:    userID,
		OrgID:     orgID,
		SessionID: sessionID,
	}
}

func (a *TokenClaims) Valid() error {
	if time.Now().After(a.ExpiresAt.Time) {
		return errors.New("token expired")
	}
	return nil
}
