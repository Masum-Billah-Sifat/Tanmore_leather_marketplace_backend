// // ------------------------------------------------------------
// // 📁 File: internal/api/http/handlers/cart/get_all_cart_items_handler.go
// // 🧠 Handles GET /api/cart/items
// //     - Extracts customer user_id from context
// //     - Calls service to fetch enriched + grouped cart data
// //     - Returns grouped valid_items and flat invalid_items
// //     - Handles moderation + stock validation in service

// package cart

// import (
// 	"net/http"

// 	service "tanmore_backend/internal/services/cart"
// 	"tanmore_backend/pkg/response"
// 	"tanmore_backend/pkg/token"

// 	"github.com/google/uuid"
// )

// // 📦 Handler struct
// type GetAllCartItemsHandler struct {
// 	Service *service.GetAllCartItemsService
// }

// // 🏗️ Constructor
// func NewGetAllCartItemsHandler(service *service.GetAllCartItemsService) *GetAllCartItemsHandler {
// 	return &GetAllCartItemsHandler{Service: service}
// }

// // 🔁 GET /api/cart/items
// func (h *GetAllCartItemsHandler) Handle(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	// 1️⃣ Extract customer user ID from context
// 	rawUserID := ctx.Value(token.CtxUserIDKey)
// 	userID, err := uuid.Parse(rawUserID.(string))
// 	if err != nil {
// 		response.Unauthorized(w, err)
// 		return
// 	}

// 	// 2️⃣ Call service layer
// 	result, err := h.Service.Start(ctx, service.GetAllCartItemsInput{
// 		UserID: userID,
// 	})
// 	if err != nil {
// 		response.ServerError(w, err)
// 		return
// 	}

// 	// 3️⃣ Return grouped valid_items and flat invalid_items
// 	response.OK(w, "Cart items fetched successfully", map[string]interface{}{
// 		"valid_items":   result.ValidItems,
// 		"invalid_items": result.InvalidItems,
// 	})
// }

// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/cart/get_all_cart_items_handler.go
// 🧠 Handles GET /api/cart/items
//     - Extracts user_id from token or guest_user_id from header
//     - Calls service to fetch grouped + enriched cart items
//     - Returns grouped valid_items and flat invalid_items

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
type GetAllCartItemsHandler struct {
	Service *service.GetAllCartItemsService
}

// 🏗️ Constructor
func NewGetAllCartItemsHandler(service *service.GetAllCartItemsService) *GetAllCartItemsHandler {
	return &GetAllCartItemsHandler{Service: service}
}

// 🔁 GET /api/cart/items
func (h *GetAllCartItemsHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	// 4️⃣ Call service layer
	result, err := h.Service.Start(ctx, service.GetAllCartItemsInput{
		UserID:      userID,
		GuestUserID: guestUserID,
	})
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 5️⃣ Return grouped valid_items and flat invalid_items
	response.OK(w, "Cart items fetched successfully", map[string]interface{}{
		"valid_items":   result.ValidItems,
		"invalid_items": result.InvalidItems,
	})
}
