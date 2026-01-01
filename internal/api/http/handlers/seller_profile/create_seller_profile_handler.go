// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/seller_profile/create_seller_profile_handler.go
// 🧠 Handles POST /api/seller/profile/metadata
//     - Extracts customer user_id from context
//     - Parses seller profile fields from JSON body
//     - Validates all required fields
//     - Calls service layer
//     - Returns seller_profile_id and status

package seller_profile

import (
	"encoding/json"
	"net/http"

	service "tanmore_backend/internal/services/seller_profile"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/response"
	"tanmore_backend/pkg/token"

	"github.com/google/uuid"
)

// 📦 Handler struct
type CreateSellerProfileHandler struct {
	Service *service.CreateSellerProfileService
}

// 🏗️ Constructor
func NewCreateSellerProfileHandler(service *service.CreateSellerProfileService) *CreateSellerProfileHandler {
	return &CreateSellerProfileHandler{Service: service}
}

// 🔁 POST /api/seller/profile/metadata
func (h *CreateSellerProfileHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1️⃣ Extract user_id from context
	rawUserID := ctx.Value(token.CtxUserIDKey)
	if rawUserID == nil {
		response.Unauthorized(w, errors.NewAuthError("missing user context"))
		return
	}

	userIDStr, ok := rawUserID.(string)
	if !ok {
		response.Unauthorized(w, errors.NewAuthError("invalid user context"))
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(w, errors.NewAuthError("invalid user id"))
		return
	}
	// if err != nil {
	// 	response.Unauthorized(w, errors.NewAuthError("invalid or missing user token"))
	// 	return
	// }

	// 2️⃣ Parse request JSON body
	var req struct {
		SellerStoreName         string  `json:"seller_store_name"`
		SellerContactNo         string  `json:"seller_contact_no"`
		SellerWhatsappContactNo string  `json:"seller_whatsapp_contact_no"`
		SellerWebsiteLink       *string `json:"seller_website_link,omitempty"`
		SellerFacebookPageName  *string `json:"seller_facebook_page_name,omitempty"`
		SellerEmail             *string `json:"seller_email,omitempty"`
		SellerPhysicalLocation  string  `json:"seller_physical_location"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, errors.NewValidationError("body", "invalid JSON body"))
		return
	}

	// 3️⃣ Validate required fields
	if req.SellerStoreName == "" {
		response.BadRequest(w, errors.NewValidationError("seller_store_name", "is required"))
		return
	}
	if req.SellerContactNo == "" {
		response.BadRequest(w, errors.NewValidationError("seller_contact_no", "is required"))
		return
	}
	if req.SellerWhatsappContactNo == "" {
		response.BadRequest(w, errors.NewValidationError("seller_whatsapp_contact_no", "is required"))
		return
	}
	if req.SellerPhysicalLocation == "" {
		response.BadRequest(w, errors.NewValidationError("seller_physical_location", "is required"))
		return
	}

	// 4️⃣ Build service input
	input := service.CreateSellerProfileInput{
		UserID:                  userID,
		SellerStoreName:         req.SellerStoreName,
		SellerContactNo:         req.SellerContactNo,
		SellerWhatsappContactNo: req.SellerWhatsappContactNo,
		SellerWebsiteLink:       req.SellerWebsiteLink,
		SellerFacebookPageName:  req.SellerFacebookPageName,
		SellerEmail:             req.SellerEmail,
		SellerPhysicalLocation:  req.SellerPhysicalLocation,
	}

	// 5️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 6️⃣ Return success
	response.Created(w, "Seller profile created successfully", map[string]interface{}{
		"seller_profile_id": result.SellerProfileID.String(),
		"status":            result.Status,
	})
}
