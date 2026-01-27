// ------------------------------------------------------------
// 📁 File: internal/services/checkout/add_shipping_address_service.go
// 🧠 Handles adding a shipping address to an existing checkout session.
//     - Verifies customer validity
//     - Validates session ownership and status
//     - Calls Pathao API for delivery charge
//     - Inserts address and updates checkout in single step

package checkout

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/checkout/add_shipping_address"
	"tanmore_backend/pkg/errors"

	// "tanmore_backend/pkg/httpclient"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"

	"github.com/google/uuid"
)

// 📥 Input
type AddShippingAddressInput struct {
	UserID            uuid.UUID
	CheckoutSessionID uuid.UUID
	RecipientName     string
	RecipientPhone    string
	RecipientEmail    *string
	AddressLine       string
	DeliveryNote      *string
	CityID            int32
	ZoneID            int32
	AreaID            int32
	Latitude          *float64
	Longitude         *float64
	PaymentMethod     string // "cod" or "prepaid"
}

// 📤 Output
type AddShippingAddressResult struct {
	Status            string
	CheckoutSessionID uuid.UUID
	ShippingAddressID uuid.UUID
}

// ⚙️ Dependencies
type AddShippingAddressService struct {
	Repo        repo.AddShippingAddressRepoInterface
	PathaoToken string
	StoreID     int // From Pathao merchant dashboard
}

// 🚀 Constructor
func NewAddShippingAddressService(repo repo.AddShippingAddressRepoInterface, pathaoToken string, storeID int) *AddShippingAddressService {
	return &AddShippingAddressService{
		Repo:        repo,
		PathaoToken: pathaoToken,
		StoreID:     storeID,
	}
}

func (s *AddShippingAddressService) Start(ctx context.Context, input AddShippingAddressInput) (*AddShippingAddressResult, error) {
	now := timeutil.NowUTC()
	var result *AddShippingAddressResult

	err := s.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// Step 1️⃣: Validate user
		user, err := q.GetUserByID(ctx, input.UserID)
		if err != nil {
			return errors.NewNotFoundError("user")
		}
		if user.IsArchived || user.IsBanned {
			return errors.NewAuthError("user is not allowed to place orders")
		}

		// Step 2️⃣: Validate session
		session, err := q.GetCheckoutSessionByID(ctx, input.CheckoutSessionID)
		if err != nil {
			return errors.NewNotFoundError("checkout_session")
		}
		if session.UserID != input.UserID {
			return errors.NewAuthError("checkout session ownership mismatch")
		}
		if session.Status != "awaiting_for_shipping" {
			return errors.NewValidationError("checkout_session", "invalid session status")
		}

		// Step 3️⃣: Count items
		itemCount, err := q.CountCheckoutItemsBySessionID(ctx, input.CheckoutSessionID)
		if err != nil {
			return errors.NewTableError("checkout_items", err.Error())
		}
		if itemCount == 0 {
			return errors.NewServerError("no_items_to_deliver")
		}

		// Step 4️⃣: Call Pathao API BEFORE inserting address
		deliveryCharge := 80 // default fallback
		pathaoReq := map[string]interface{}{
			"store_id":       s.StoreID,
			"item_type":      2,
			"delivery_type":  48,
			"item_weight":    max(0.5, float64(itemCount)*0.5),
			"recipient_city": input.CityID,
			"recipient_zone": input.ZoneID,
		}

		body, _ := json.Marshal(pathaoReq)
		req, _ := http.NewRequest("POST", "https://sandbox.pathao.com/aladdin/api/v1/merchant/price-plan", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.PathaoToken)

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)

		if err == nil && resp.StatusCode == 200 {
			var parsed struct {
				Data struct {
					FinalPrice float64 `json:"final_price"`
				} `json:"data"`
			}
			_ = json.NewDecoder(resp.Body).Decode(&parsed)
			if parsed.Data.FinalPrice > 0 {
				deliveryCharge = int(parsed.Data.FinalPrice)
			}
		}

		subtotal, _ := strconv.Atoi(session.Subtotal)
		discount := 0
		if session.PlatformDiscountAmountApplied.Valid {
			discount, _ = strconv.Atoi(session.PlatformDiscountAmountApplied.String)
		}
		totalPayable := subtotal - discount + deliveryCharge

		shippingID := uuid.New()
		err = q.UpdateCheckoutSessionWithShipping(ctx, sqlc.UpdateCheckoutSessionWithShippingParams{
			ShippingAddressID: sqlnull.UUID(shippingID),
			PaymentMethod:     input.PaymentMethod,
			DeliveryCharge:    sqlnull.String(strconv.Itoa(deliveryCharge)),
			TotalPayable:      strconv.Itoa(totalPayable),

			ID: input.CheckoutSessionID,
		})
		if err != nil {
			return errors.NewTableError("checkout_sessions.update_shipping", err.Error())
		}

		_, err = q.InsertShippingAddress(ctx, sqlc.InsertShippingAddressParams{
			ID:                shippingID,
			CheckoutSessionID: input.CheckoutSessionID,
			RecipientName:     input.RecipientName,
			RecipientPhone:    input.RecipientPhone,
			RecipientEmail:    sqlnull.StringPtr(input.RecipientEmail),
			AddressLine:       input.AddressLine,
			DeliveryNote:      sqlnull.StringPtr(input.DeliveryNote),
			CityID:            input.CityID,
			ZoneID:            input.ZoneID,
			AreaID:            input.AreaID,

			Latitude:  sqlnull.StringPtr(floatPtrToString(input.Latitude)),
			Longitude: sqlnull.StringPtr(floatPtrToString(input.Longitude)),

			CreatedAt: sqlnull.Time(now),
		})
		if err != nil {
			return errors.NewTableError("shipping_addresses.insert", err.Error())
		}

		result = &AddShippingAddressResult{
			Status:            "shipping_address_added",
			CheckoutSessionID: input.CheckoutSessionID,
			ShippingAddressID: shippingID,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

// func max(a, b float64) float64 {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }

// func floatPtrToString(f *float64) *string {
// 	if f == nil {
// 		return nil
// 	}
// 	s := fmt.Sprintf("%.7f", *f)
// 	return &s
// }
