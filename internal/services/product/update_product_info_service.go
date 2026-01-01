// ------------------------------------------------------------
// 📁 File: internal/services/product/update_product_info_service.go
// 🧠 Handles updating the title and/or description of a product.
//     - Validates seller identity and moderation
//     - Confirms product ownership
//     - Updates optional fields using COALESCE
//     - Emits product.info.updated event
//     - Returns updated_fields and product ID

package product

import (
	"context"
	"database/sql"
	"encoding/json"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_update_info"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type UpdateProductInfoInput struct {
	UserID      uuid.UUID
	ProductID   uuid.UUID
	Title       *string // Optional
	Description *string // Optional
}

// ------------------------------------------------------------
// 📤 Result to return
type UpdateProductInfoResult struct {
	ProductID     uuid.UUID
	UpdatedFields map[string]string
	Updated       bool
}

// ------------------------------------------------------------
// 🧱 Dependencies
type UpdateProductInfoServiceDeps struct {
	Repo repo.ProductUpdateInfoRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service Definition
type UpdateProductInfoService struct {
	Deps UpdateProductInfoServiceDeps
}

// 🚀 Constructor
func NewUpdateProductInfoService(deps UpdateProductInfoServiceDeps) *UpdateProductInfoService {
	return &UpdateProductInfoService{Deps: deps}
}

// 🚀 Entrypoint
func (s *UpdateProductInfoService) Start(
	ctx context.Context,
	input UpdateProductInfoInput,
) (*UpdateProductInfoResult, error) {

	now := timeutil.NowUTC()

	var updatedFields map[string]string

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// ------------------------------------------------------------
		// Step 1: Validate seller user
		user, err := q.GetUserByID(ctx, input.UserID)
		if err != nil {
			return errors.NewNotFoundError("user")
		}
		if user.IsBanned {
			return errors.NewAuthError("user is banned")
		}
		if user.IsArchived {
			return errors.NewAuthError("user is archived")
		}
		if !user.IsSellerProfileApproved {
			return errors.NewValidationError("seller", "profile is not approved")
		}

		// ------------------------------------------------------------
		// Step 2: Confirm product ownership
		product, err := q.GetProductByIDAndSellerID(ctx, sqlc.GetProductByIDAndSellerIDParams{
			ID:       input.ProductID,
			SellerID: input.UserID,
		})
		if err != nil {
			return errors.NewNotFoundError("product")
		}
		if product.IsBanned || product.IsArchived {
			return errors.NewValidationError("product", "cannot update banned or archived product")
		}

		// ------------------------------------------------------------
		// Step 3: Update product fields (COALESCE-based)
		var titleNull, descNull sql.NullString
		updatedFields = make(map[string]string)

		if input.Title != nil {
			titleNull = sqlnull.String(*input.Title)
			updatedFields["title"] = *input.Title
		}
		if input.Description != nil {
			descNull = sqlnull.String(*input.Description)
			updatedFields["description"] = *input.Description
		}

		err = q.UpdateProductTitleDesc(ctx, sqlc.UpdateProductTitleDescParams{
			Title:       titleNull,
			Description: descNull,
			UpdatedAt:   now,
			ProductID:   input.ProductID,
			SellerID:    input.UserID,
		})
		if err != nil {
			return errors.NewTableError("products.update", err.Error())
		}

		// ------------------------------------------------------------
		// Step 4: Emit product.info.updated event
		eventPayload := map[string]interface{}{
			"user_id":        input.UserID,
			"product_id":     input.ProductID,
			"updated_fields": updatedFields,
		}

		payloadBytes, err := json.Marshal(eventPayload)
		if err != nil {
			return errors.NewServerError("marshal event payload")
		}

		err = q.InsertEvent(ctx, sqlc.InsertEventParams{
			ID:           uuidutil.New(),
			Userid:       input.UserID,
			EventType:    "product.info.updated",
			EventPayload: payloadBytes,
			DispatchedAt: sqlnull.TimePtr(nil), // NULL = not yet sent
			CreatedAt:    now,
		})
		if err != nil {
			return errors.NewTableError("events.insert", err.Error())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &UpdateProductInfoResult{
		ProductID:     input.ProductID,
		UpdatedFields: updatedFields,
		Updated:       true,
	}, nil
}
