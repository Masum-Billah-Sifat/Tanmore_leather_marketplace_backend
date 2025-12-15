// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/user/switch_mode_handler.go
// 🧠 Handles the POST /api/user/switch-mode endpoint. Parses mode from JWT,
//     reads body, and calls service to update current mode and issue new token.

package user

import (
	"encoding/json"
	"net/http"

	"tanmore_backend/internal/services/user_mode_switch"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/google/uuid"
)

// 📦 SwitchModeHandler wires the handler to service
type SwitchModeHandler struct {
	Service *user_mode_switch.SwitchModeService
}

// 🚀 Constructor
func NewSwitchModeHandler(service *user_mode_switch.SwitchModeService) *SwitchModeHandler {
	return &SwitchModeHandler{Service: service}
}

// 📥 Request JSON
type switchModeRequest struct {
	ToMode string `json:"to_mode"`
}

// 🔁 POST /api/user/switch-mode
func (h *SwitchModeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Decode request body
	var body switchModeRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, err)
		return
	}

	// 2️⃣ Extract from context
	rawUserID := ctx.Value(token.CtxUserIDKey)
	rawSessionID := ctx.Value(token.CtxSessionIDKey)
	rawMode := ctx.Value(token.CtxModeKey)

	// 3️⃣ Basic type assertions (no validation yet)
	userID, _ := uuid.Parse(rawUserID.(string))
	sessionID, _ := uuid.Parse(rawSessionID.(string))
	fromMode := rawMode.(string)

	// 4️⃣ Create service input
	input := user_mode_switch.SwitchModeInput{
		UserID:    userID,
		SessionID: sessionID,
		FromMode:  fromMode,
		ToMode:    body.ToMode,
	}

	// 5️⃣ Invoke service
	output, err := h.Service.Handle(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 6️⃣ Send success
	response.OK(w, "Mode switched successfully", map[string]interface{}{
		"access_token": output.AccessToken,
		"expires_in":   output.ExpiresIn,
		"mode":         output.Mode,
	})
}
