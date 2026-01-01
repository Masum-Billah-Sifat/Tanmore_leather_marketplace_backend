// 🧠 Repository implementation for product_variant_indexes dynamic querying.
//
//	Handles both feed and search using dynamic query builder.
package product_variant_index_feed_or_search

import (
	"context"
	"database/sql"

	feedquery "tanmore_backend/pkg/feed_query"
)

type ProductVariantIndexFeedOrSearchRepository struct {
	db *sql.DB
}

// 🚀 Constructor
func NewProductVariantIndexFeedOrSearchRepository(db *sql.DB) *ProductVariantIndexFeedOrSearchRepository {
	return &ProductVariantIndexFeedOrSearchRepository{db: db}
}

// 🔍 Run dynamic feed/search query
func (r *ProductVariantIndexFeedOrSearchRepository) RunFeedOrSearchQuery(
	ctx context.Context,
	params feedquery.FeedQueryParams,
) ([]feedquery.ProductVariantIndexFeedOrSearchRow, error) {
	query, args := feedquery.BuildDynamicFeedQuery(params)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []feedquery.ProductVariantIndexFeedOrSearchRow
	for rows.Next() {
		var row feedquery.ProductVariantIndexFeedOrSearchRow
		if err := rows.Scan(
			// ─ Primary Key
			&row.ID,

			// ─ Category
			&row.CategoryID,
			&row.IsCategoryArchived,
			&row.CategoryName,

			// ─ Seller
			&row.SellerID,
			&row.IsSellerApproved,
			&row.IsSellerArchived,
			&row.IsSellerBanned,
			&row.SellerStoreName,

			// ─ Product
			&row.ProductID,
			&row.IsProductApproved,
			&row.IsProductArchived,
			&row.IsProductBanned,
			&row.ProductTitle,
			&row.ProductDescription,
			&row.ProductImages,
			&row.ProductPromoVideoURL,

			// ─ Variant
			&row.VariantID,
			&row.IsVariantArchived,
			&row.IsVariantInStock,
			&row.StockAmount,
			&row.Color,
			&row.Size,
			&row.RetailPrice,
			&row.RetailDiscountType,
			&row.RetailDiscount,
			&row.HasRetailDiscount,
			&row.HasWholesaleEnabled,
			&row.WholesalePrice,
			&row.WholesaleMinQuantity,
			&row.WholesaleDiscountType,
			&row.WholesaleDiscount,
			&row.WeightGrams,

			// ─ Metadata
			&row.Views,
			&row.RelevanceScore,

			// ─ Timestamps
			&row.CreatedAt,
			&row.UpdatedAt,
		); err != nil {
			return nil, err
		}

		results = append(results, row)
	}

	return results, nil
}
