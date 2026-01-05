// ------------------------------------------------------------
// 📁 File: internal/services/product/update_product_category_service.go
// 🧠 Handles updating the category of a product.
//     - Validates seller identity and moderation
//     - Confirms product ownership
//     - Confirms new category is leaf & not archived
//     - Updates category_id
//     - Emits product_category_updated event

package product

import (
	"context"
	"encoding/json"

	sqlc "tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_update_category"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input

type UpdateProductCategoryInput struct {
	UserID     uuid.UUID
	ProductID  uuid.UUID
	CategoryID uuid.UUID
}

// ------------------------------------------------------------
// 📤 Output

type UpdateProductCategoryResult struct {
	ProductID           uuid.UUID `json:"product_id"`
	UpdatedCategoryID   uuid.UUID `json:"updated_category_id"`
	UpdatedCategoryName string    `json:"updated_category_name"`
	Updated             bool      `json:"updated"`
}

// ------------------------------------------------------------
// 🧱 Dependencies

type UpdateProductCategoryServiceDeps struct {
	Repo repo.ProductUpdateCategoryRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service Definition

type UpdateProductCategoryService struct {
	Deps UpdateProductCategoryServiceDeps
}

// 🚀 Constructor

func NewUpdateProductCategoryService(deps UpdateProductCategoryServiceDeps) *UpdateProductCategoryService {
	return &UpdateProductCategoryService{Deps: deps}
}

// ------------------------------------------------------------
// 🚀 Entrypoint

func (s *UpdateProductCategoryService) Start(
	ctx context.Context,
	input UpdateProductCategoryInput,
) (*UpdateProductCategoryResult, error) {

	now := timeutil.NowUTC()

	var result *UpdateProductCategoryResult

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// Step 1: Validate user
		user, err := q.GetUserByID(ctx, input.UserID)
		if err != nil {
			return errors.NewNotFoundError("user")
		}
		if user.IsBanned || user.IsArchived || !user.IsSellerProfileApproved || !user.IsSellerProfileCreated {
			return errors.NewAuthError("invalid seller access")
		}

		// Step 2: Validate product
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

		// Step 3: Validate new category
		category, err := q.GetCategoryByID(ctx, input.CategoryID)
		if err != nil {
			return errors.NewNotFoundError("category")
		}
		if category.IsArchived || !category.IsLeaf {
			return errors.NewValidationError("category", "must be a valid, active leaf node")
		}

		// Step 4: Update product's category
		err = q.UpdateProductCategory(ctx, sqlc.UpdateProductCategoryParams{
			ProductID:  input.ProductID,
			SellerID:   input.UserID,
			CategoryID: input.CategoryID,
			UpdatedAt:  now,
		})
		if err != nil {
			return errors.NewTableError("products.update_category", err.Error())
		}

		// Step 5: Insert event
		eventPayload := map[string]interface{}{
			"product_id":        input.ProductID,
			"seller_id":         input.UserID,
			"new_category_id":   category.ID,
			"new_category_name": category.Name,
		}
		payloadBytes, err := json.Marshal(eventPayload)
		if err != nil {
			return errors.NewServerError("marshal event payload")
		}

		err = q.InsertEvent(ctx, sqlc.InsertEventParams{
			ID:           uuidutil.New(),
			Userid:       input.UserID,
			EventType:    "product_category_updated",
			EventPayload: payloadBytes,
			DispatchedAt: sqlnull.TimePtr(nil),
			CreatedAt:    now,
		})
		if err != nil {
			return errors.NewTableError("events.insert", err.Error())
		}

		// Prepare result
		result = &UpdateProductCategoryResult{
			ProductID:           input.ProductID,
			UpdatedCategoryID:   category.ID,
			UpdatedCategoryName: category.Name,
			Updated:             true,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
