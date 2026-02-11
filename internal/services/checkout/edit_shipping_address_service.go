// ------------------------------------------------------------
// 📁 File: internal/services/checkout/edit_shipping_address_service.go
// 🧠 Handles editing a shipping address for a checkout session.
//     - Verifies customer validity
//     - Validates session and ownership
//     - Validates target shipping address
//     - Optionally updates Pathao delivery fee if city/zone changed

package checkout

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	sqlc "tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/checkout/edit_shipping_address"
	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"

	// "tanmore_backend/pkg/timeutil"

	"bytes"
	"encoding/json"

	"github.com/google/uuid"
)

// 📅 Input
type EditShippingAddressInput struct {
	UserID            uuid.UUID
	CheckoutSessionID uuid.UUID
	ShippingAddressID uuid.UUID

	RecipientName  *string
	RecipientPhone *string
	RecipientEmail *string
	AddressLine    *string
	DeliveryNote   *string
	CityID         *int32
	ZoneID         *int32
	AreaID         *int32
	Latitude       *float64
	Longitude      *float64
	PaymentMethod  *string // optional
}

// 📥 Output
type EditShippingAddressResult struct {
	CheckoutSessionID uuid.UUID
	ShippingAddressID uuid.UUID
	Status            string
}

// ⚙️ Dependencies
type EditShippingAddressService struct {
	Repo        repo.EditShippingAddressRepoInterface
	PathaoToken string
	StoreID     int
}

// 🚀 Constructor
func NewEditShippingAddressService(repo repo.EditShippingAddressRepoInterface, pathaoToken string, storeID int) *EditShippingAddressService {
	return &EditShippingAddressService{
		Repo:        repo,
		PathaoToken: pathaoToken,
		StoreID:     storeID,
	}
}

func (s *EditShippingAddressService) Start(ctx context.Context, input EditShippingAddressInput) (*EditShippingAddressResult, error) {
	// now := timeutil.NowUTC()
	var result *EditShippingAddressResult

	err := s.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// Step 1: Validate customer
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

		// Step 2: Validate session
		session, err := q.GetCheckoutSessionByID(ctx, input.CheckoutSessionID)
		if err != nil || session.UserID != input.UserID {
			return errors.NewAuthError("unauthorized access to session")
		}

		// after session fetch + ownership check
		if session.Status != "ready_to_order" {
			return errors.NewValidationError("checkout_session", "invalid session status")
		}

		if !session.ShippingAddressID.Valid || session.ShippingAddressID.UUID != input.ShippingAddressID {
			return errors.NewValidationError("shipping_address_id", "does not match checkout session")
		}

		// Step 3: Validate shipping address belongs to session
		address, err := q.GetShippingAddressByIDAndCheckoutID(ctx, sqlc.GetShippingAddressByIDAndCheckoutIDParams{
			ID:                input.ShippingAddressID,
			CheckoutSessionID: input.CheckoutSessionID,
		})
		if err != nil {
			return errors.NewNotFoundError("shipping_address")
		}

		// Step 4: Conditionally update shipping address
		err = q.UpdateShippingAddressByID(ctx, sqlc.UpdateShippingAddressByIDParams{
			ID:             input.ShippingAddressID,
			RecipientName:  sqlnull.StringPtr(input.RecipientName),
			RecipientPhone: sqlnull.StringPtr(input.RecipientPhone),
			RecipientEmail: sqlnull.StringPtr(input.RecipientEmail),
			AddressLine:    sqlnull.StringPtr(input.AddressLine),
			DeliveryNote:   sqlnull.StringPtr(input.DeliveryNote),
			// CityID:         sqlnull.Int32Ptr(input.CityID),
			// ZoneID:         sqlnull.Int32Ptr(input.ZoneID),
			// AreaID:         sqlnull.Int32Ptr(input.AreaID),
			CityID: sqlnull.Int32Ptr(int32PtrToInt64Ptr(input.CityID)),
			ZoneID: sqlnull.Int32Ptr(int32PtrToInt64Ptr(input.ZoneID)),
			AreaID: sqlnull.Int32Ptr(int32PtrToInt64Ptr(input.AreaID)),

			Latitude:  sqlnull.StringPtr(floatPtrToString(input.Latitude)),
			Longitude: sqlnull.StringPtr(floatPtrToString(input.Longitude)),
		})
		if err != nil {
			return errors.NewTableError("shipping_addresses.update", err.Error())
		}

		// Step 5: Update payment method if provided
		if input.PaymentMethod != nil {
			err := q.UpdateCheckoutSessionPaymentMethod(ctx, sqlc.UpdateCheckoutSessionPaymentMethodParams{
				PaymentMethod: *input.PaymentMethod,
				ID:            input.CheckoutSessionID,
			})
			if err != nil {
				return errors.NewTableError("checkout_sessions.payment_method", err.Error())
			}
		}

		// Step 6: If city/zone changed → recalculate Pathao delivery fee
		cityChanged := input.CityID != nil && *input.CityID != address.CityID
		zoneChanged := input.ZoneID != nil && *input.ZoneID != address.ZoneID
		if cityChanged || zoneChanged {
			itemCount, err := q.CountCheckoutItemsBySessionID(ctx, input.CheckoutSessionID)
			if err != nil || itemCount == 0 {
				return errors.NewServerError("failed item count or zero items")
			}

			deliveryCharge := 80 // fallback
			pathaoReq := map[string]interface{}{
				"store_id":       s.StoreID,
				"item_type":      2,
				"delivery_type":  48,
				"item_weight":    max(0.5, float64(itemCount)*0.5),
				"recipient_city": *input.CityID,
				"recipient_zone": *input.ZoneID,
			}
			body, _ := json.Marshal(pathaoReq)
			req, _ := http.NewRequest("POST", "https://sandbox.pathao.com/aladdin/api/v1/merchant/price-plan", bytes.NewReader(body))
			req.Header.Set("Authorization", "Bearer "+s.PathaoToken)
			req.Header.Set("Content-Type", "application/json")

			resp, err := (&http.Client{Timeout: 5 * time.Second}).Do(req)
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

			err = q.UpdateCheckoutSessionDeliveryPricing(ctx, sqlc.UpdateCheckoutSessionDeliveryPricingParams{
				DeliveryCharge: sqlnull.String(strconv.Itoa(deliveryCharge)),
				TotalPayable:   strconv.Itoa(totalPayable),
				ID:             input.CheckoutSessionID,
			})
			if err != nil {
				return errors.NewTableError("checkout_sessions.update_delivery_charge", err.Error())
			}
		}

		result = &EditShippingAddressResult{
			CheckoutSessionID: input.CheckoutSessionID,
			ShippingAddressID: input.ShippingAddressID,
			Status:            "shipping_address_updated",
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func floatPtrToString(f *float64) *string {
	if f == nil {
		return nil
	}
	s := fmt.Sprintf("%.7f", *f)
	return &s
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func int32PtrToInt64Ptr(p *int32) *int64 {
	if p == nil {
		return nil
	}
	v := int64(*p)
	return &v
}
