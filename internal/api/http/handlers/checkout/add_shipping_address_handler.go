// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/checkout/add_shipping_address_handler.go
// 🧑‍💼 Handles POST /api/checkout/{checkout_session_id}/add-shipping-address
//     - Extracts authenticated user_id from context
//     - Parses checkout_session_id from route param
//     - Parses JSON body with shipping address and payment method
//     - Calls service layer
//     - Returns checkout_session_id and shipping_address_id

package checkout

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"tanmore_backend/internal/services/checkout"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/google/uuid"
)

// 🛆 Handler struct
type AddShippingAddressHandler struct {
	Service *checkout.AddShippingAddressService
}

// 📇 Constructor
func NewAddShippingAddressHandler(service *checkout.AddShippingAddressService) *AddShippingAddressHandler {
	return &AddShippingAddressHandler{Service: service}
}

func (h *AddShippingAddressHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	// 2️⃣ Extract checkout_session_id from URL path
	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 4 {
		response.BadRequest(w, errors.NewValidationError("checkout_session_id", "missing in URL path"))
		return
	}
	sessionIDStr := segments[3] // assuming path like /api/checkout/{id}/add-shipping-address
	checkoutSessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("checkout_session_id", "invalid UUID format"))
		return
	}

	// 3️⃣ Decode request body
	var req struct {
		RecipientName  string   `json:"recipient_name"`
		RecipientPhone string   `json:"recipient_phone"`
		RecipientEmail *string  `json:"recipient_email"`
		AddressLine    string   `json:"address_line"`
		DeliveryNote   *string  `json:"delivery_note"`
		CityID         int32    `json:"city_id"`
		ZoneID         int32    `json:"zone_id"`
		AreaID         int32    `json:"area_id"`
		Latitude       *float64 `json:"latitude"`
		Longitude      *float64 `json:"longitude"`
		PaymentMethod  string   `json:"payment_method"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
		return
	}

	// 4️⃣ Basic validation
	if req.RecipientName == "" || req.RecipientPhone == "" || req.AddressLine == "" || req.CityID == 0 || req.ZoneID == 0 || req.AreaID == 0 {
		response.BadRequest(w, errors.NewValidationError("fields", "missing required address fields"))
		return
	}
	if req.PaymentMethod != "cod" && req.PaymentMethod != "prepaid" {
		response.BadRequest(w, errors.NewValidationError("payment_method", "must be 'cod' or 'prepaid'"))
		return
	}
	if (req.Latitude != nil && req.Longitude == nil) || (req.Latitude == nil && req.Longitude != nil) {
		response.BadRequest(w, errors.NewValidationError("coordinates", "both latitude and longitude must be provided together"))
		return
	}

	// 5️⃣ Call service layer
	result, err := h.Service.Start(ctx, checkout.AddShippingAddressInput{
		UserID:            userID,
		CheckoutSessionID: checkoutSessionID,
		RecipientName:     req.RecipientName,
		RecipientPhone:    req.RecipientPhone,
		RecipientEmail:    req.RecipientEmail,
		AddressLine:       req.AddressLine,
		DeliveryNote:      req.DeliveryNote,
		CityID:            req.CityID,
		ZoneID:            req.ZoneID,
		AreaID:            req.AreaID,
		Latitude:          req.Latitude,
		Longitude:         req.Longitude,
		PaymentMethod:     req.PaymentMethod,
	})
	if err != nil {
		fmt.Println("❌ Service error:", err)
		response.ServerError(w, err)
		return
	}

	// 6️⃣ Respond success
	response.Created(w, "Shipping address added", map[string]interface{}{
		"checkout_session_id": result.CheckoutSessionID,
		"shipping_address_id": result.ShippingAddressID,
		"status":              result.Status,
	})
}
