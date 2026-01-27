// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/cart/merge_guest_cart_handler.go
// 🧠 Handles POST /api/cart/merge-guest
//     - Extracts authenticated user ID from access token
//     - Extracts guest_user_id from header
//     - Validates both presence and exclusivity
//     - Calls service layer to merge guest cart into user cart
//     - Returns merge status and count

package cart

import (
	"fmt"
	"net/http"

	service "tanmore_backend/internal/services/cart"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/google/uuid"
)

// 📦 Handler struct
type MergeGuestCartHandler struct {
	Service *service.MergeGuestCartService
}

// 🏗️ Constructor
func NewMergeGuestCartHandler(service *service.MergeGuestCartService) *MergeGuestCartHandler {
	return &MergeGuestCartHandler{Service: service}
}

// 🔁 POST /api/cart/merge-guest
func (h *MergeGuestCartHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Extract authenticated user ID from access token (middleware)
	rawUserID := ctx.Value(token.CtxUserIDKey)
	if rawUserID == nil {
		fmt.Println("❌ Missing user ID in context")
		response.Unauthorized(w, errors.NewAuthError("unauthorized"))
		return
	}
	userID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		fmt.Println("❌ Invalid user ID format:", err)
		response.Unauthorized(w, err)
		return
	}
	fmt.Println("✅ Parsed authenticated user ID:", userID)

	// 2️⃣ Extract guest_user_id from header
	guestHeader := r.Header.Get("X-Tanmore-Guest-Id")
	if guestHeader == "" {
		fmt.Println("❌ Missing guest_user_id header")
		response.BadRequest(w, errors.NewValidationError("guest_user_id", "header is required"))
		return
	}
	guestUserID, err := uuid.Parse(guestHeader)
	if err != nil {
		fmt.Println("❌ Invalid guest_user_id UUID:", err)
		response.BadRequest(w, errors.NewValidationError("guest_user_id", "invalid UUID format"))
		return
	}
	fmt.Println("✅ Parsed guest user ID from header:", guestUserID)

	// 3️⃣ Call service layer
	result, err := h.Service.Start(ctx, service.MergeGuestCartInput{
		UserID:      userID,
		GuestUserID: guestUserID,
	})
	if err != nil {
		fmt.Println("❌ Service returned error:", err)
		response.ServerError(w, err)
		return
	}

	// 4️⃣ Send response
	response.OK(w, "Guest cart merged successfully", result)
}
