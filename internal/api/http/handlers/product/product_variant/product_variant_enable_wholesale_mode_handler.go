// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product_variant/enable_variant_wholesale_mode_handler.go
// 🧠 Handles POST /api/seller/products/:product_id/variants/:variant_id/wholesale-mode
//     - Extracts seller user_id from context
//     - Extracts product_id and variant_id from path
//     - Parses wholesale pricing and discount fields
//     - Calls service layer
//     - Returns variant_id and wholesale enabled status

package product_variant

import (
	"encoding/json"
	"net/http"

	service "tanmore_backend/internal/services/product/product_variant"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// 📦 Handler struct
type EnableVariantWholesaleModeHandler struct {
	Service *service.EnableWholesaleModeService
}

// 🛠️ Constructor
func NewEnableVariantWholesaleModeHandler(service *service.EnableWholesaleModeService) *EnableVariantWholesaleModeHandler {
	return &EnableVariantWholesaleModeHandler{Service: service}
}

// 🔁 POST /api/seller/products/:product_id/variants/:variant_id/wholesale-mode
func (h *EnableVariantWholesaleModeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Get seller user ID from context
	rawUserID := ctx.Value(token.CtxUserIDKey)
	userID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		response.Unauthorized(w, err)
		return
	}

	// 2️⃣ Get product_id from path
	productIDParam := chi.URLParam(r, "product_id")
	productID, err := uuid.Parse(productIDParam)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("product_id", "invalid product ID"))
		return
	}

	// 3️⃣ Get variant_id from path
	variantIDParam := chi.URLParam(r, "variant_id")
	variantID, err := uuid.Parse(variantIDParam)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("variant_id", "invalid variant ID"))
		return
	}

	// 4️⃣ Parse request JSON body
	var req struct {
		WholesalePrice        int64   `json:"wholesale_price"`
		MinQtyWholesale       int64   `json:"min_qty_wholesale"`
		WholesaleDiscount     *int64  `json:"wholesale_discount,omitempty"`      // optional
		WholesaleDiscountType *string `json:"wholesale_discount_type,omitempty"` // optional: "flat" or "percentage"
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
		return
	}

	// 5️⃣ Validate required fields
	if req.WholesalePrice <= 0 {
		response.BadRequest(w, errors.NewValidationError("wholesale_price", "must be a positive number"))
		return
	}
	if req.MinQtyWholesale <= 0 {
		response.BadRequest(w, errors.NewValidationError("min_qty_wholesale", "must be a positive number"))
		return
	}

	// 6️⃣ Validate optional discount type (if provided)
	if req.WholesaleDiscountType != nil {
		if *req.WholesaleDiscountType != "flat" && *req.WholesaleDiscountType != "percentage" {
			response.BadRequest(w, errors.NewValidationError("wholesale_discount_type", "must be 'flat' or 'percentage'"))
			return
		}
	}

	// 7️⃣ Build service input
	input := service.EnableWholesaleModeInput{
		UserID:                userID,
		ProductID:             productID,
		VariantID:             variantID,
		WholesalePrice:        req.WholesalePrice,
		MinQtyWholesale:       req.MinQtyWholesale,
		WholesaleDiscount:     req.WholesaleDiscount,
		WholesaleDiscountType: req.WholesaleDiscountType,
	}

	// 8️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 9️⃣ Return success
	response.OK(w, "Wholesale mode enabled successfully", map[string]interface{}{
		"variant_id":        result.VariantID,
		"wholesale_enabled": result.WholesaleEnabled,
		"status":            result.Status,
	})
}
