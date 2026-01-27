// // ------------------------------------------------------------
// // 📁 File: internal/api/http/handlers/cart/clear_cart_handler.go
// // 🧠 Handles DELETE /api/cart/clear
// //     - Extracts customer user_id from context
// //     - Calls the service layer to clear all active cart items
// //     - Returns cart cleared or already empty status

// package cart

// import (
// 	"net/http"

// 	service "tanmore_backend/internal/services/cart"
// 	"tanmore_backend/pkg/errors"
// 	"tanmore_backend/pkg/response"
// 	"tanmore_backend/pkg/token"

// 	"github.com/google/uuid"
// )

// // 📦 Handler struct
// type ClearCartHandler struct {
// 	Service *service.ClearCartService
// }

// // 🏗️ Constructor
// func NewClearCartHandler(service *service.ClearCartService) *ClearCartHandler {
// 	return &ClearCartHandler{Service: service}
// }

// // 🔁 DELETE /api/cart/clear
// func (h *ClearCartHandler) Handle(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	// 1️⃣ Get customer user ID from context
// 	rawUserID := ctx.Value(token.CtxUserIDKey)
// 	userID, err := uuid.Parse(rawUserID.(string))
// 	if err != nil {
// 		response.Unauthorized(w, errors.NewAuthError("invalid access token"))
// 		return
// 	}

// 	// 2️⃣ Build service input
// 	input := service.ClearCartInput{
// 		UserID: userID,
// 	}

// 	// 3️⃣ Call service
// 	result, err := h.Service.Start(ctx, input)
// 	if err != nil {
// 		response.ServerError(w, err)
// 		return
// 	}

// 	// 4️⃣ Return response
// 	response.OK(w, "Cart status", map[string]string{
// 		"status": result.Status,
// 	})
// }

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
type ClearCartHandler struct {
	Service *service.ClearCartService
}

// 🏗️ Constructor
func NewClearCartHandler(service *service.ClearCartService) *ClearCartHandler {
	return &ClearCartHandler{Service: service}
}

// 🔁 DELETE /api/cart/clear
func (h *ClearCartHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Extract user ID from access token (if any)
	var userID *uuid.UUID
	rawUserID := ctx.Value(token.CtxUserIDKey)
	if rawUserID != nil {
		parsed, err := uuid.Parse(rawUserID.(string))
		if err != nil {
			fmt.Println("❌ Failed to parse user ID:", err)
			response.Unauthorized(w, err)
			return
		}
		userID = &parsed
		fmt.Println("✅ Parsed authenticated user ID:", *userID)
	}

	// 2️⃣ Extract guest_user_id from header (if any)
	var guestUserID *uuid.UUID
	guestHeader := r.Header.Get("X-Tanmore-Guest-Id")
	if guestHeader != "" {
		parsed, err := uuid.Parse(guestHeader)
		if err != nil {
			fmt.Println("❌ Invalid guest_user_id UUID:", err)
			response.BadRequest(w, errors.NewValidationError("guest_user_id", "invalid UUID format in header"))
			return
		}
		guestUserID = &parsed
		fmt.Println("✅ Parsed guest user ID from header:", *guestUserID)
	}

	// 3️⃣ Enforce: exactly one of user_id or guest_user_id must be present
	if (userID == nil && guestUserID == nil) || (userID != nil && guestUserID != nil) {
		fmt.Println("❌ Must provide either user_id (token) or guest_user_id (header), but not both")
		response.BadRequest(w, errors.NewValidationError("authorization", "either user_id or guest_user_id must be provided, not both"))
		return
	}

	// 4️⃣ Build service input
	input := service.ClearCartInput{
		UserID:      userID,
		GuestUserID: guestUserID,
	}

	// 5️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 6️⃣ Return success
	response.OK(w, "Cart status", map[string]string{
		"status": result.Status,
	})
}
