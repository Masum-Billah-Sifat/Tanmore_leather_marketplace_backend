// ------------------------------------------------------------
// 📁 File: internal/api/http/handler/google_auth/google_auth_handler.go
// 🧠 Handles the POST /api/auth/google route. Extracts headers & JSON body,
//     calls the service to perform Google login or registration,
//     and returns tokens + user info.

package google_auth

import (
	"encoding/json"
	"net/http"

	googleauthsvc "tanmore_backend/internal/services/google_auth"
	"tanmore_backend/pkg/response"
)

// 📦 Handler struct holds dependencies
type Handler struct {
	Service *googleauthsvc.GoogleAuthService
}

// 🛠️ Constructor
func NewHandler(service *googleauthsvc.GoogleAuthService) *Handler {
	return &Handler{Service: service}
}

// 📦 Request body
type googleLoginRequest struct {
	IDToken string `json:"id_token"`
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body googleLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, err)
		return
	}

	userAgent := r.Header.Get("User-Agent")
	platform := r.Header.Get("X-Platform")
	deviceFP := r.Header.Get("X-Device-Fingerprint")
	ipAddress := r.RemoteAddr

	input := googleauthsvc.GoogleLoginInput{
		IDToken:           body.IDToken,
		UserAgent:         userAgent,
		Platform:          platform,
		DeviceFingerprint: deviceFP,
		IPAddress:         ipAddress,
	}

	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	response.OK(w, "Login successful", map[string]interface{}{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"expires_in":    result.ExpiresIn,
		"user": map[string]interface{}{
			"id":                         result.User.ID.String(),
			"name":                       result.User.Name,
			"email":                      result.User.Email,
			"image":                      result.User.Image,
			"is_seller_profile_approved": result.User.IsSellerProfileApproved, // ✅ added
		},
	})
}

// // 🚪 POST /api/auth/google
// func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	// 📨 Parse JSON body
// 	var body googleLoginRequest
// 	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
// 		response.BadRequest(w, err)
// 		return
// 	}

// 	// 🧩 Extract headers
// 	userAgent := r.Header.Get("User-Agent")
// 	platform := r.Header.Get("X-Platform")
// 	deviceFP := r.Header.Get("X-Device-Fingerprint")
// 	ipAddress := r.RemoteAddr

// 	// 📦 Map to service input
// 	input := googleauthsvc.GoogleLoginInput{
// 		IDToken:           body.IDToken,
// 		UserAgent:         userAgent,
// 		Platform:          platform,
// 		DeviceFingerprint: deviceFP,
// 		IPAddress:         ipAddress,
// 	}

// 	// 🚀 Call service
// 	result, err := h.Service.Start(ctx, input)
// 	if err != nil {
// 		response.ServerError(w, err)
// 		return
// 	}

// 	// ✅ Respond with tokens + user info
// 	response.OK(w, "Login successful", map[string]interface{}{
// 		"access_token":  result.AccessToken,
// 		"refresh_token": result.RefreshToken,
// 		"expires_in":    result.ExpiresIn,
// 		"user": map[string]interface{}{
// 			"id":    result.User.ID.String(),
// 			"name":  result.User.Name,
// 			"email": result.User.Email,
// 			"image": result.User.Image,
// 		},
// 	})
// }
