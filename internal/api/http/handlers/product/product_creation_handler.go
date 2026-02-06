// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/product/create_product_handler.go
// 🧠 Handles POST /api/seller/products
//     - Parses JSON body
//     - Validates required fields
//     - Extracts seller user_id from context
//     - Calls service layer
//     - Returns full product object on success

package product

import (
	"encoding/json"
	"net/http"

	service "tanmore_backend/internal/services/product"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/google/uuid"
)

// 📦 CreateProductHandler wires handler to service
type CreateProductHandler struct {
	Service *service.CreateProductService
}

// 🚀 Constructor
func NewCreateProductHandler(service *service.CreateProductService) *CreateProductHandler {
	return &CreateProductHandler{Service: service}
}

// 📥 Request Body Structs

type createProductVariantRequest struct {
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

type createProductRequest struct {
	CategoryID    string                        `json:"category_id"`
	Title         string                        `json:"title"`
	Description   string                        `json:"description"`
	ImageURLs     []string                      `json:"image_urls"`
	PromoVideoURL *string                       `json:"promo_video_url,omitempty"`
	Variants      []createProductVariantRequest `json:"variants"`
}

// 🔁 POST /api/seller/products
func (h *CreateProductHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Decode request body
	var body createProductRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, err)
		return
	}

	// Step 1️⃣: Extract user ID from access token context
	rawUserID := ctx.Value(token.CtxUserIDKey)
	if rawUserID == nil {
		response.Unauthorized(w, errors.ErrAuthMissingToken())
		return
	}

	userID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		response.Unauthorized(w, errors.ErrAuthInvalidUserID())
		return
	}
	// 3️⃣ Parse category ID
	categoryID, err := uuid.Parse(body.CategoryID)
	if err != nil {
		response.BadRequest(w, err)
		return
	}

	// 4️⃣ Validate basic fields
	if body.Title == "" || body.Description == "" || len(body.ImageURLs) == 0 || len(body.Variants) == 0 {
		response.BadRequest(
			w,
			errors.NewValidationError(
				"request_body",
				"missing required fields: title, description, image_urls, variants",
			),
		)
		return
	}

	// 5️⃣ Build variant inputs
	var variants []service.CreateProductVariantInput
	for _, v := range body.Variants {
		variants = append(variants, service.CreateProductVariantInput{
			Color:                 v.Color,
			Size:                  v.Size,
			RetailPrice:           v.RetailPrice,
			InStock:               v.InStock,
			StockQuantity:         v.StockQuantity,
			RetailDiscount:        v.RetailDiscount,
			RetailDiscountType:    v.RetailDiscountType,
			WholesalePrice:        v.WholesalePrice,
			MinQtyWholesale:       v.MinQtyWholesale,
			WholesaleDiscount:     v.WholesaleDiscount,
			WholesaleDiscountType: v.WholesaleDiscountType,
			WeightGrams:           v.WeightGrams,
		})
	}

	// 6️⃣ Build final service input
	input := service.CreateProductInput{
		UserID:        userID,
		CategoryID:    categoryID,
		Title:         body.Title,
		Description:   body.Description,
		ImageURLs:     body.ImageURLs,
		PromoVideoURL: body.PromoVideoURL,
		Variants:      variants,
	}

	// 7️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// ✅ 8️⃣ Return full product in response (Updated)
	response.Created(w, "Product created successfully", result.Product)
}
