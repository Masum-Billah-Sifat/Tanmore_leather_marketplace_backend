// ------------------------------------------------------------
// 📁 File: pkg/token/token.go
// 🧠 This file provides reusable helpers for generating and validating
//     access and refresh tokens used throughout the authentication system.

package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// 🔑 JWT Secret (should come from env/config in production)
var jwtSecret = []byte("supersecretkey") // Replace with env

// 🧱 Access token payload
type AccessTokenClaims struct {
	Sub  string `json:"sub"`
	SID  string `json:"sid"`
	Mode string `json:"mode"`
	jwt.RegisteredClaims
}

// 🚀 GenerateAccessToken creates a signed JWT token for access
func GenerateAccessToken(userID uuid.UUID, sessionID uuid.UUID, mode string, expiryMinutes int) (string, error) {
	now := time.Now()
	claims := AccessTokenClaims{
		Sub:  userID.String(),
		SID:  sessionID.String(),
		Mode: mode,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expiryMinutes) * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// 🔍 ParseAccessToken decodes and validates a JWT access token
func ParseAccessToken(tokenStr string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid access token")
	}

	return claims, nil
}

// ✅ Context keys
type ctxKey string

const (
	CtxUserIDKey    ctxKey = "user_id"
	CtxSessionIDKey ctxKey = "session_id"
	CtxModeKey      ctxKey = "current_mode"
)

// ✅ AttachAccessToken middleware — parses token and sets values into context
func AttachAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			next.ServeHTTP(w, r) // skip silently if no token
			return
		}

		rawToken := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := ParseAccessToken(rawToken)
		if err != nil {
			// Skip adding to context if parsing fails
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), CtxUserIDKey, claims.Sub)
		ctx = context.WithValue(ctx, CtxSessionIDKey, claims.SID)
		ctx = context.WithValue(ctx, CtxModeKey, claims.Mode)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// 🔐 GenerateRefreshToken creates a secure random token string (64 chars)
func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 48)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// 🧪 HashRefreshToken hashes a refresh token (to store in DB securely)
func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}

// 🔁 CompareRefreshTokenHash checks if token matches stored hash
func CompareRefreshTokenHash(storedHash, providedToken string) bool {
	return storedHash == HashRefreshToken(providedToken)
}
