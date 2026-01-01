// ------------------------------------------------------------
// 📁 File: internal/services/product/archive_product_service.go
// 🧠 Handles soft-archiving (logical deletion) of a product.
//     - Validates seller identity
//     - Confirms product ownership
//     - Performs soft archive
//     - Emits product.archived event
//     - Returns product_id and status

package product

import (
	"context"
	"encoding/json"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_archive"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type ArchiveProductInput struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
}

// ------------------------------------------------------------
// 📤 Output to return
type ArchiveProductResult struct {
	ProductID uuid.UUID `json:"product_id"`
	Status    string    `json:"status"`
}

// ------------------------------------------------------------
// 🧱 Dependencies
type ArchiveProductServiceDeps struct {
	Repo repo.ProductArchiveRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service Definition
type ArchiveProductService struct {
	Deps ArchiveProductServiceDeps
}

// 🚀 Constructor
func NewArchiveProductService(deps ArchiveProductServiceDeps) *ArchiveProductService {
	return &ArchiveProductService{Deps: deps}
}

// 🚀 Entrypoint
func (s *ArchiveProductService) Start(
	ctx context.Context,
	input ArchiveProductInput,
) (*ArchiveProductResult, error) {

	now := timeutil.NowUTC()

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// ------------------------------------------------------------
		// Step 1: Validate seller
		user, err := q.GetUserByID(ctx, input.UserID)
		if err != nil {
			return errors.NewNotFoundError("seller")
		}
		if user.IsArchived {
			return errors.NewAuthError("seller is archived")
		}
		if user.IsBanned {
			return errors.NewAuthError("seller is banned")
		}
		if !user.IsSellerProfileApproved {
			return errors.NewAuthError("seller profile not approved")
		}

		// ------------------------------------------------------------
		// Step 2: Confirm product ownership + moderation check
		product, err := q.GetProductByIDAndSellerID(ctx, sqlc.GetProductByIDAndSellerIDParams{
			ID:       input.ProductID,
			SellerID: input.UserID,
		})
		if err != nil {
			return errors.NewNotFoundError("product")
		}
		if product.IsBanned {
			return errors.NewValidationError("product", "is banned")
		}
		if product.IsArchived {
			return errors.NewValidationError("product", "is already archived")
		}

		// ------------------------------------------------------------
		// Step 3: Archive product
		err = q.ArchiveProduct(ctx, sqlc.ArchiveProductParams{
			ID:         input.ProductID,
			SellerID:   input.UserID,
			IsArchived: true,
			UpdatedAt:  now,
		})
		if err != nil {
			return errors.NewTableError("products.archive", err.Error())
		}

		// ------------------------------------------------------------
		// Step 4: Emit product.archived event
		payload, err := json.Marshal(map[string]interface{}{
			"seller_id":  input.UserID,
			"product_id": input.ProductID,
		})
		if err != nil {
			return errors.NewServerError("marshal event payload")
		}

		err = q.InsertEvent(ctx, sqlc.InsertEventParams{
			ID:           uuidutil.New(),
			Userid:       input.UserID,
			EventType:    "product.archived",
			EventPayload: payload,
			DispatchedAt: sqlnull.TimePtr(nil),
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

	return &ArchiveProductResult{
		ProductID: input.ProductID,
		Status:    "product_archived",
	}, nil
}
