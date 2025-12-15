// ------------------------------------------------------------
// 📁 File: internal/api/http/handler/google_auth/refresh_token_handler.go
// 🧠 Handles the POST /api/auth/refresh endpoint. Parses input and headers,
//     securely rotates refresh token via service, and returns new token pair.

package google_auth

import (
	"encoding/json"
	"net/http"

	"tanmore_backend/internal/services/google_auth"
	"tanmore_backend/pkg/response"
)

// 📦 RefreshTokenHandler handles refresh logic
type RefreshTokenHandler struct {
	Service *google_auth.RefreshTokenService
}

// 🚀 Constructor
func NewRefreshTokenHandler(service *google_auth.RefreshTokenService) *RefreshTokenHandler {
	return &RefreshTokenHandler{Service: service}
}

// 📥 Request payload structure
type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// 🔁 POST /api/auth/refresh
func (h *RefreshTokenHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Decode request body
	var body refreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, err)
		return
	}

	// 2️⃣ Extract headers
	userAgent := r.Header.Get("User-Agent")
	platform := r.Header.Get("X-Platform")
	deviceFingerprint := r.Header.Get("X-Device-Fingerprint")
	ipAddress := r.RemoteAddr

	// 3️⃣ Create input for service
	input := google_auth.RefreshTokenInput{
		RawToken:          body.RefreshToken,
		UserAgent:         userAgent,
		Platform:          platform,
		DeviceFingerprint: deviceFingerprint,
		IPAddress:         ipAddress,
	}

	// 4️⃣ Invoke service layer
	result, err := h.Service.HandleRefreshTokenRotation(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 5️⃣ Return response
	response.OK(w, "Refresh successful", map[string]interface{}{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"expires_in":    result.ExpiresIn,
	})
}
