// ------------------------------------------------------------
// 📁 File: internal/services/seller_profile/create_seller_profile_service.go
// 🧠 Handles creation of seller profile metadata.
//     - Validates user (must be unarchived, unbanned, not yet seller)
//     - Inserts seller profile metadata
//     - Updates user flag `is_seller_profile_created = true`

package seller_profile

import (
	"context"

	"tanmore_backend/internal/db/sqlc"
	repo "tanmore_backend/internal/repository/seller_profile/seller_profile_metadata"

	"tanmore_backend/pkg/errors"
	"tanmore_backend/pkg/sqlnull"
	"tanmore_backend/pkg/timeutil"

	"github.com/google/uuid"
)

// ------------------------------------------------------------
// 📥 Input from handler
type CreateSellerProfileInput struct {
	UserID                  uuid.UUID
	SellerStoreName         string
	SellerContactNo         string
	SellerWhatsappContactNo string
	SellerWebsiteLink       *string
	SellerFacebookPageName  *string
	SellerEmail             *string
	SellerPhysicalLocation  string
}

// ------------------------------------------------------------
// 📤 Result to return
type CreateSellerProfileResult struct {
	SellerProfileID uuid.UUID
	Status          string
}

// ------------------------------------------------------------
// 🧱 Dependencies
type CreateSellerProfileServiceDeps struct {
	Repo repo.SellerProfileMetadataRepoInterface
}

// ------------------------------------------------------------
// 🛠️ Service Definition
type CreateSellerProfileService struct {
	Deps CreateSellerProfileServiceDeps
}

// 🚀 Constructor
func NewCreateSellerProfileService(deps CreateSellerProfileServiceDeps) *CreateSellerProfileService {
	return &CreateSellerProfileService{Deps: deps}
}

// 🚀 Entrypoint
func (s *CreateSellerProfileService) Start(
	ctx context.Context,
	input CreateSellerProfileInput,
) (*CreateSellerProfileResult, error) {

	now := timeutil.NowUTC()

	var sellerProfileID uuid.UUID

	err := s.Deps.Repo.WithTx(ctx, func(q *sqlc.Queries) error {
		// ------------------------------------------------------------
		// Step 1: Validate user
		user, err := q.GetUserByID(ctx, input.UserID)
		if err != nil {
			return errors.NewNotFoundError("user")
		}
		if user.IsArchived {
			return errors.NewAuthError("user is archived")
		}
		if user.IsBanned {
			return errors.NewAuthError("user is banned")
		}
		if user.IsSellerProfileCreated {
			return errors.NewValidationError("seller_profile", "already created")
		}
		if user.IsSellerProfileApproved {
			return errors.NewValidationError("seller_profile", "already approved")
		}

		// ------------------------------------------------------------
		// Step 2: Insert seller profile metadata
		sellerProfileID, err = q.InsertSellerProfileMetadata(ctx, sqlc.InsertSellerProfileMetadataParams{
			SellerID:                input.UserID,
			Sellerstorename:         input.SellerStoreName,
			Sellercontactno:         input.SellerContactNo,
			Sellerwhatsappcontactno: input.SellerWhatsappContactNo,
			Sellerwebsitelink:       sqlnull.StringPtr(input.SellerWebsiteLink),
			Sellerfacebookpagename:  sqlnull.StringPtr(input.SellerFacebookPageName),
			Selleremail:             sqlnull.StringPtr(input.SellerEmail),
			Sellerphysicallocation:  input.SellerPhysicalLocation,
			CreatedAt:               now,
			UpdatedAt:               now,
		})
		if err != nil {
			return errors.NewTableError("seller_profile_metadata.insert", err.Error())
		}

		// ------------------------------------------------------------
		// Step 3: Mark user as profile created
		err = q.UpdateSellerProfileCreated(ctx, sqlc.UpdateSellerProfileCreatedParams{
			IsSellerProfileCreated: true,
			ID:                     input.UserID,
			UpdatedAt:              now,
		})
		if err != nil {
			return errors.NewTableError("users.update", err.Error())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &CreateSellerProfileResult{
		SellerProfileID: sellerProfileID,
		Status:          "seller_profile_created",
	}, nil
}
