// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product_variant/add_product_variant_handler.go
// 🧠 Handles POST /api/seller/products/:product_id/variants
//     - Parses JSON body
//     - Extracts seller user_id from context
//     - Extracts product_id from path
//     - Validates required fields
//     - Calls service layer
//     - Returns product_id and variant_id

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
type AddProductVariantHandler struct {
	Service *service.AddProductVariantService
}

// 🏗️ Constructor
func NewAddProductVariantHandler(service *service.AddProductVariantService) *AddProductVariantHandler {
	return &AddProductVariantHandler{Service: service}
}

// 📥 Request body
type addVariantRequest struct {
	Color                 string  `json:"color"`
	Size                  string  `json:"size"`
	RetailPrice           int64   `json:"retail_price"`
	InStock               bool    `json:"in_stock"`
	StockQuantity         int64   `json:"stock_quantity"`
	RetailDiscount        *int64  `json:"retail_discount,omitempty"`
	RetailDiscountType    *string `json:"retail_discount_type,omitempty"`
	WholesalePrice        *int64  `json:"wholesale_price,omitempty"`
	MinQtyWholesale       *int64  `json:"min_qty_wholesale,omitempty"`
	WholesaleDiscount     *int64  `json:"wholesale_discount,omitempty"`
	WholesaleDiscountType *string `json:"wholesale_discount_type,omitempty"`
	WeightGrams           int64   `json:"weight_grams"`
}

// 🔁 POST /api/seller/products/:product_id/variants
func (h *AddProductVariantHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Decode JSON
	var body addVariantRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, err)
		return
	}

	// 2️⃣ Get seller user ID from context
	rawUserID := ctx.Value(token.CtxUserIDKey)
	userID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		response.Unauthorized(w, err)
		return
	}

	// 3️⃣ Get product_id from path
	productIDParam := chi.URLParam(r, "product_id")
	productID, err := uuid.Parse(productIDParam)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("product_id", "invalid product ID"))
		return
	}

	// 4️⃣ Validate required fields
	if body.Color == "" || body.Size == "" || body.RetailPrice <= 0 || body.StockQuantity <= 0 || body.WeightGrams <= 0 {
		response.BadRequest(
			w,
			errors.NewValidationError(
				"request_body",
				"missing required fields: color, size, retail_price, stock_quantity, weight_grams",
			),
		)
		return
	}

	// 5️⃣ Build service input
	input := service.AddProductVariantInput{
		UserID:                userID,
		ProductID:             productID,
		Color:                 body.Color,
		Size:                  body.Size,
		RetailPrice:           body.RetailPrice,
		InStock:               body.InStock,
		StockQuantity:         body.StockQuantity,
		RetailDiscount:        body.RetailDiscount,
		RetailDiscountType:    body.RetailDiscountType,
		WholesalePrice:        body.WholesalePrice,
		MinQtyWholesale:       body.MinQtyWholesale,
		WholesaleDiscount:     body.WholesaleDiscount,
		WholesaleDiscountType: body.WholesaleDiscountType,
		WeightGrams:           body.WeightGrams,
	}

	// 6️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 7️⃣ Return success
	response.Created(w, "Variant added successfully", map[string]interface{}{
		"product_id": result.ProductID,
		"variant_id": result.VariantID,
		"status":     "variant_created",
	})
}
