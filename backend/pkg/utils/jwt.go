package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret          []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewJWTManager(secret string, accessTokenTTL, refreshTokenTTL time.Duration) *JWTManager {
	return &JWTManager{secret: []byte(secret), accessTokenTTL: accessTokenTTL, refreshTokenTTL: refreshTokenTTL}
}

func (m *JWTManager) GenerateTokenPair(userID uuid.UUID, issuer string) (TokenPair, error) {
	access, err := m.generate(userID, issuer, m.accessTokenTTL)
	if err != nil {
		return TokenPair{}, err
	}
	refresh, err := m.generate(userID, issuer, m.refreshTokenTTL)
	if err != nil {
		return TokenPair{}, err
	}
	return TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (m *JWTManager) generate(userID uuid.UUID, issuer string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTManager) Verify(tokenString string, issuer string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return m.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}), jwt.WithIssuer(issuer))
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
