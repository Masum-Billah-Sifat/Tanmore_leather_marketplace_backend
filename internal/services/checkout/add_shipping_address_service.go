// ------------------------------------------------------------
// 📁 File: internal/services/checkout/add_shipping_address_service.go
// 🧠 Handles adding a shipping address to an existing checkout session.
//     - Verifies customer validity
//     - Validates session ownership and status
//     - Calls Pathao API for delivery charge
//     - Inserts address and updates checkout in single transaction
//     - ✅ Marks checkout session as ready_to_order when done

package checkout

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	stdErrors "errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/checkout/add_shipping_address"
	"tanmore_backend/pkg/errors"
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
			return errors.ErrAuthUserNotFound()
		}
		if user.IsArchived {
			return errors.ErrAuthArchivedUser()
		}
		if user.IsBanned {
			return errors.ErrAuthBannedUser()
		}

		// Optional guard (handler likely already validates)
		if input.PaymentMethod != "cod" && input.PaymentMethod != "prepaid" {
			return errors.NewValidationError("payment_method", "must be 'cod' or 'prepaid'")
		}

		// Step 2️⃣: Validate session
		session, err := q.GetCheckoutSessionByID(ctx, input.CheckoutSessionID)
		if err != nil {
			return errors.NewNotFoundError("checkout_session")
		}
		if session.UserID != input.UserID {
			return errors.NewAuthError("checkout session ownership mismatch")
		}
		if session.Status != "awaiting_shipping_info" {
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
		req, _ := http.NewRequest(
			"POST",
			"https://sandbox.pathao.com/aladdin/api/v1/merchant/price-plan",
			bytes.NewReader(body),
		)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.PathaoToken)

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
		}

		if err == nil && resp != nil && resp.StatusCode == 200 {
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

		shippingID := uuid.New()

		// Money math in cents
		subtotalCents, err := moneyToCentsTwo(session.Subtotal)
		if err != nil {
			return errors.NewServerError("invalid subtotal format in checkout session")
		}

		discountCents := int64(0)
		if session.PlatformDiscountAmountApplied.Valid {
			d, err := moneyToCentsTwo(session.PlatformDiscountAmountApplied.String)
			if err != nil {
				return errors.NewServerError("invalid discount format in checkout session")
			}
			discountCents = d
		}

		deliveryChargeCents := int64(deliveryCharge) * 100

		totalPayableCents := subtotalCents - discountCents + deliveryChargeCents
		if totalPayableCents < 0 {
			totalPayableCents = 0
		}

		// Update checkout totals + attach shipping_id
		err = q.UpdateCheckoutSessionWithShipping(ctx, sqlc.UpdateCheckoutSessionWithShippingParams{
			ShippingAddressID: sqlnull.UUID(shippingID),
			PaymentMethod:     input.PaymentMethod,
			DeliveryCharge:    sqlnull.String(centsToMoneyTwo(deliveryChargeCents)),
			TotalPayable:      centsToMoneyTwo(totalPayableCents),
			ID:                input.CheckoutSessionID,
		})
		if err != nil {
			return errors.NewTableError("checkout_sessions.update_shipping", err.Error())
		}

		// Insert shipping address row
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
			Latitude:          sqlnull.StringPtr(floatPtrToString(input.Latitude)),
			Longitude:         sqlnull.StringPtr(floatPtrToString(input.Longitude)),
			CreatedAt:         sqlnull.Time(now),
		})
		if err != nil {
			return errors.NewTableError("shipping_addresses.insert", err.Error())
		}

		// ✅ NEW: status -> ready_to_order (guards status transition)
		_, err = q.MarkCheckoutSessionReadyToOrder(ctx, sqlc.MarkCheckoutSessionReadyToOrderParams{
			ID:     input.CheckoutSessionID,
			UserID: input.UserID,
		})
		if err != nil {
			// common case when WHERE status='awaiting_shipping_info' didn't match
			if stdErrors.Is(err, sql.ErrNoRows) {
				return errors.NewValidationError("checkout_session", "invalid session status")
			}
			return errors.NewTableError("checkout_sessions.status", err.Error())
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

func moneyToCentsTwo(s string) (int64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return int64(math.Round(f * 100)), nil
}

func centsToMoneyTwo(cents int64) string {
	return fmt.Sprintf("%.2f", float64(cents)/100)
}
