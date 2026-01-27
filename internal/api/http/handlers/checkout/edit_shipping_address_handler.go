// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/checkout/edit_shipping_address_handler.go
// 🧑‍💼 Handles PUT /api/checkout/{checkout_session_id}/shipping-address/{shipping_address_id}/edit
//     - Extracts authenticated user_id from context
//     - Parses checkout_session_id and shipping_address_id from route params
//     - Parses optional JSON body fields
//     - Validates conditional rules (lat/lng, payment_method)
//     - Calls service layer
//     - Returns checkout_session_id and shipping_address_id

package checkout

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	service "tanmore_backend/internal/services/checkout"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/google/uuid"
)

// 📦 Handler struct
type EditShippingAddressHandler struct {
	Service *service.EditShippingAddressService
}

// 🏗️ Constructor
func NewEditShippingAddressHandler(service *service.EditShippingAddressService) *EditShippingAddressHandler {
	return &EditShippingAddressHandler{Service: service}
}

func (h *EditShippingAddressHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Extract authenticated user_id
	rawUserID := ctx.Value(token.CtxUserIDKey)
	if rawUserID == nil {
		response.Unauthorized(w, errors.NewAuthError("unauthenticated access"))
		return
	}

	userID, err := uuid.Parse(rawUserID.(string))
	if err != nil {
		response.Unauthorized(w, errors.NewAuthError("invalid user id"))
		return
	}

	// 2️⃣ Extract IDs from URL path
	// Expected:
	// /api/checkout/{checkout_session_id}/shipping-address/{shipping_address_id}/edit
	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 7 {
		response.BadRequest(w, errors.NewValidationError("path", "invalid URL format"))
		return
	}

	checkoutSessionID, err := uuid.Parse(segments[3])
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("checkout_session_id", "invalid UUID"))
		return
	}

	shippingAddressID, err := uuid.Parse(segments[5])
	if err != nil {
		response.BadRequest(w, errors.NewValidationError("shipping_address_id", "invalid UUID"))
		return
	}

	// 3️⃣ Decode request body (ALL OPTIONAL)
	var req struct {
		RecipientName  *string  `json:"recipient_name"`
		RecipientPhone *string  `json:"recipient_phone"`
		RecipientEmail *string  `json:"recipient_email"`
		AddressLine    *string  `json:"address_line"`
		DeliveryNote   *string  `json:"delivery_note"`
		CityID         *int32   `json:"city_id"`
		ZoneID         *int32   `json:"zone_id"`
		AreaID         *int32   `json:"area_id"`
		Latitude       *float64 `json:"latitude"`
		Longitude      *float64 `json:"longitude"`
		PaymentMethod  *string  `json:"payment_method"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
		return
	}

	// 4️⃣ Conditional validations

	// lat/lng must come together
	if (req.Latitude != nil && req.Longitude == nil) ||
		(req.Latitude == nil && req.Longitude != nil) {
		response.BadRequest(w, errors.NewValidationError(
			"coordinates",
			"latitude and longitude must be provided together",
		))
		return
	}

	// payment_method validation
	if req.PaymentMethod != nil {
		if *req.PaymentMethod != "cod" && *req.PaymentMethod != "prepaid" {
			response.BadRequest(w, errors.NewValidationError(
				"payment_method",
				"must be 'cod' or 'prepaid'",
			))
			return
		}
	}

	// 5️⃣ Call service layer
	result, err := h.Service.Start(ctx, service.EditShippingAddressInput{
		UserID:            userID,
		CheckoutSessionID: checkoutSessionID,
		ShippingAddressID: shippingAddressID,

		RecipientName:  req.RecipientName,
		RecipientPhone: req.RecipientPhone,
		RecipientEmail: req.RecipientEmail,
		AddressLine:    req.AddressLine,
		DeliveryNote:   req.DeliveryNote,

		CityID: req.CityID,
		ZoneID: req.ZoneID,
		AreaID: req.AreaID,

		Latitude:  req.Latitude,
		Longitude: req.Longitude,

		PaymentMethod: req.PaymentMethod,
	})
	if err != nil {
		fmt.Println("❌ Service error:", err)
		response.ServerError(w, err)
		return
	}

	// 6️⃣ Success response
	response.OK(w, "Shipping address updated", map[string]interface{}{
		"checkout_session_id": result.CheckoutSessionID,
		"shipping_address_id": result.ShippingAddressID,
		"status":              result.Status,
	})
}
