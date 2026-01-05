// ------------------------------------------------------------
// 📁 File: internal/api/http/handler/logout/logout_handler.go
// 🧠 Handles the POST /api/auth/logout route. Extracts user_id and session_id
//     from context, calls service layer to revoke session + deprecate tokens.

package logout

import (
	"net/http"

	logoutsvc "tanmore_backend/internal/services/logout"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/google/uuid"
)

// 📦 Handler struct holds dependencies
type Handler struct {
	Service *logoutsvc.LogoutService
}

// 🛠️ Constructor
func NewHandler(service *logoutsvc.LogoutService) *Handler {
	return &Handler{Service: service}
}

// 🚪 POST /api/auth/logout
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 🧩 Extract from context (middleware stores as string)
	userIDStr := ctx.Value(token.CtxUserIDKey)
	sessionIDStr := ctx.Value(token.CtxSessionIDKey)

	if userIDStr == nil || sessionIDStr == nil {
		response.Unauthorized(w, errors.NewAuthError("missing token claims in context"))
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		response.Unauthorized(w, errors.NewAuthError("invalid user_id in token"))
		return
	}

	sessionID, err := uuid.Parse(sessionIDStr.(string))
	if err != nil {
		response.Unauthorized(w, errors.NewAuthError("invalid session_id in token"))
		return
	}

	// 🚀 Call service
	result, err := h.Service.HandleLogout(ctx, logoutsvc.LogoutInput{
		UserID:    userID,
		SessionID: sessionID,
	})
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// ✅ Success response
	response.OK(w, result.Message, nil)
}
