// ------------------------------------------------------------
// 📁 File: internal/services/product_variant/remove_product_variant_service.go
// 🧠 Handles soft-deleting (archiving) a product variant.
//     - Validates seller & ownership via snapshot
//     - Performs soft delete
//     - Emits variant.archived event
//     - Returns variant_id and product_id

package product_variant

import (
	"context"
	"encoding/json"
	"time"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/product/product_variant/product_variant_archive"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"
	uuidutil "tanmore_backend/pkg/uuid"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type RemoveProductVariantInput struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	VariantID uuid.UUID
}

// ------------------------------------------------------------
// 📤 Result to return
type RemoveProductVariantResult struct {
	ProductID uuid.UUID
	VariantID uuid.UUID
	Status    string
}

// ------------------------------------------------------------
// 🧱 Dependencies
type RemoveProductVariantServiceDeps struct {
	Repo repo.ProductVariantArchiveRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service Definition
type RemoveProductVariantService struct {
	Deps RemoveProductVariantServiceDeps
}

// 🚀 Constructor
func NewRemoveProductVariantService(deps RemoveProductVariantServiceDeps) *RemoveProductVariantService {
	return &RemoveProductVariantService{Deps: deps}
}

// 🚀 Entrypoint
func (s *RemoveProductVariantService) Start(
	ctx context.Context,
	input RemoveProductVariantInput,
) (*RemoveProductVariantResult, error) {

	now := timeutil.NowUTC()

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// ------------------------------------------------------------
		// Step 1: Fetch snapshot (seller + product + variant + category)
		snapshot, err := q.GetVariantSnapshot(ctx, sqlc.GetVariantSnapshotParams{
			Sellerid:  input.UserID,
			Productid: input.ProductID,
			Variantid: input.VariantID,
		})
		if err != nil {
			return errors.NewNotFoundError("variant snapshot")
		}

		// ------------------------------------------------------------
		// Step 2: Check moderation rules
		if snapshot.Issellerarchived || snapshot.Issellerbanned {
			return errors.NewAuthError("seller is not allowed to modify products")
		}
		if snapshot.Isproductarchived || snapshot.Isproductbanned {
			return errors.NewValidationError("product", "cannot modify archived or banned product")
		}
		if snapshot.Iscategoryarchived {
			return errors.NewValidationError("category", "product's category is archived")
		}
		if snapshot.Isvariantarchived {
			return errors.NewValidationError("variant", "variant is already archived")
		}

		// ------------------------------------------------------------
		// Step 3: Perform soft-delete
		err = q.ArchiveProductVariant(ctx, sqlc.ArchiveProductVariantParams{
			ID:         input.VariantID,
			ProductID:  input.ProductID,
			IsArchived: true,
			UpdatedAt:  now,
		})
		if err != nil {
			return errors.NewTableError("product_variants.archive", err.Error())
		}

		// ------------------------------------------------------------
		// Step 4: Emit variant.archived event
		payload, err := json.Marshal(map[string]interface{}{
			"user_id":    input.UserID,
			"product_id": input.ProductID,
			"variant_id": input.VariantID,
		})
		if err != nil {
			return errors.NewServerError("marshal event payload")
		}

		err = q.InsertEvent(ctx, sqlc.InsertEventParams{
			ID:           uuidutil.New(),
			Userid:       input.UserID,
			EventType:    "variant.archived",
			EventPayload: payload,
			DispatchedAt: sqlnull.Time(time.Time{}),
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

	return &RemoveProductVariantResult{
		ProductID: input.ProductID,
		VariantID: input.VariantID,
		Status:    "variant_archived",
	}, nil
}
