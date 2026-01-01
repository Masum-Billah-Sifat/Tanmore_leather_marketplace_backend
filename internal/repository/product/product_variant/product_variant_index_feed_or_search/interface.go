// 🧠 Repository interface for dynamic product variant index feed/search queries.
package product_variant_index_feed_or_search

import (
	"context"

	feedquery "tanmore_backend/pkg/feed_query"
)

type ProductVariantIndexFeedOrSearchRepo interface {
	// 🔍 Run dynamic feed/search query and return rows
	RunFeedOrSearchQuery(
		ctx context.Context,
		params feedquery.FeedQueryParams,
	) ([]feedquery.ProductVariantIndexFeedOrSearchRow, error)
}
