// ------------------------------------------------------------
// 📁 File: pkg/token/token.go
// 🧠 This file provides reusable helpers for generating and validating
//     access and refresh tokens used throughout the authentication system.

package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	// "errors"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"
	"time"

	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// 🔑 JWT Secret (should come from env/config in production)
var jwtSecret = []byte("supersecretkey") // Replace with env

// correct access token claims as of now
type AccessTokenClaims struct {
	Sub  uuid.UUID `json:"sub"`
	SID  uuid.UUID `json:"sid"`
	Mode string    `json:"mode"`
	jwt.RegisteredClaims
}

// fixed generate access token function
func GenerateAccessToken(
	userID uuid.UUID,
	sessionID uuid.UUID,
	mode string,
	expiryMinutes int,
) (string, error) {

	now := time.Now().UTC()

	claims := AccessTokenClaims{
		Sub:  userID,
		SID:  sessionID,
		Mode: mode,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expiryMinutes) * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

// fixed parsedaccesstokenfunction
func ParseAccessToken(tokenStr string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&AccessTokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)

	if err != nil {
		// covers expired, malformed, invalid signature, etc.
		return nil, errors.ErrAuthInvalidToken()
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.ErrAuthInvalidToken()
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

func AttachAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Unauthorized(w, errors.ErrAuthMissingToken())
			return
		}

		rawToken := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := ParseAccessToken(rawToken)
		if err != nil {
			response.Unauthorized(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), CtxUserIDKey, claims.Sub.String())
		ctx = context.WithValue(ctx, CtxSessionIDKey, claims.SID.String())
		ctx = context.WithValue(ctx, CtxModeKey, claims.Mode)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AttachAccessTokenForCustomer(next http.Handler) http.Handler {
	return AttachAccessToken(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mode, _ := r.Context().Value(CtxModeKey).(string)
		if mode != "customer" {
			response.Forbidden(w, errors.ErrAuthOnlyCustomer())
			return
		}
		next.ServeHTTP(w, r)
	}))
}

func AttachAccessTokenForSeller(next http.Handler) http.Handler {
	return AttachAccessToken(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mode, _ := r.Context().Value(CtxModeKey).(string)
		if mode != "seller" {
			response.Forbidden(w, errors.ErrAuthOnlySeller())
			return
		}
		next.ServeHTTP(w, r)
	}))
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

// --------------------------------------------------
// func AttachAccessToken(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		authHeader := r.Header.Get("Authorization")
// 		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
// 			fmt.Println("❌ Missing Authorization header")
// 			response.Unauthorized(w, errors.NewAuthError("missing access token"))
// 			return
// 		}

// 		rawToken := strings.TrimPrefix(authHeader, "Bearer ")

// 		claims, err := ParseAccessToken(rawToken)
// 		if err != nil {
// 			fmt.Println("❌ Access token parse error:", err)
// 			response.Unauthorized(w, err)
// 			return
// 		}

// 		fmt.Println("✅ Token parsed. Sub:", claims.Sub.String(), "SID:", claims.SID.String(), "Mode:", claims.Mode)

// 		// Add to context
// 		ctx := context.WithValue(r.Context(), CtxUserIDKey, claims.Sub.String())
// 		ctx = context.WithValue(ctx, CtxSessionIDKey, claims.SID.String())
// 		ctx = context.WithValue(ctx, CtxModeKey, claims.Mode)

// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

// added two new ones as of now
// func AttachAccessTokenForCustomer(next http.Handler) http.Handler {
// 	return AttachAccessToken(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		mode := r.Context().Value(CtxModeKey)
// 		if modeStr, ok := mode.(string); !ok || modeStr != "customer" {
// 			response.Forbidden(w, errors.NewAuthError("only accessible to customers"))
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	}))
// }

// func AttachAccessTokenForSeller(next http.Handler) http.Handler {
// 	return AttachAccessToken(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		mode := r.Context().Value(CtxModeKey)
// 		if modeStr, ok := mode.(string); !ok || modeStr != "seller" {
// 			response.Forbidden(w, errors.NewAuthError("only accessible to sellers"))
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	}))
// }
