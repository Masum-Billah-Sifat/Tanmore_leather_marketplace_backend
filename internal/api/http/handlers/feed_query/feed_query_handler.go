// ------------------------------------------------------------
// 📁 File: internal/api/http/handlers/feed_query/feed_query_handler.go
// 🧠 Handles GET /api/public/feed and /api/public/search
//     - Parses query parameters (q, filters, sort, page, per_page)
//     - Builds service input
//     - Calls service layer
//     - Returns a paginated list of product-variant cards

package feed_query

import (
	"errors"
	"net/http"
	"strconv"

	service "tanmore_backend/internal/services/product/product_variant/product_variant_index_feed_or_search"
	"tanmore_backend/pkg/response"
	uuidutil "tanmore_backend/pkg/uuid"
)

// 📦 Handler struct
type FeedQueryHandler struct {
	Service *service.FeedQueryService
}

// 🏗️ Constructor
func NewFeedQueryHandler(service *service.FeedQueryService) *FeedQueryHandler {
	return &FeedQueryHandler{Service: service}
}

// 🔁 GET /api/public/feed
func (h *FeedQueryHandler) HandleFeed(w http.ResponseWriter, r *http.Request) {
	h.handleFeedOrSearch(w, r, false)
}

// 🔍 GET /api/public/search
func (h *FeedQueryHandler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	h.handleFeedOrSearch(w, r, true)
}

func (h *FeedQueryHandler) handleFeedOrSearch(w http.ResponseWriter, r *http.Request, isSearch bool) {
	ctx := r.Context()

	// 1️⃣ Parse common filters
	page := parseInt(r.URL.Query().Get("page"), 1)
	perPage := parseInt(r.URL.Query().Get("per_page"), 10)

	query := r.URL.Query().Get("q")
	if isSearch && (query == "" || len(query) < 1) {
		response.BadRequest(w, errors.New("Missing search query. Please provide a non-empty 'q' parameter."))
		return
	}

	categoryID := uuidutil.ParsePtr(r.URL.Query().Get("category_id"))
	minPrice := parseInt64Ptr(r.URL.Query().Get("min_price"))
	maxPrice := parseInt64Ptr(r.URL.Query().Get("max_price"))
	minWeight := parseIntPtr(r.URL.Query().Get("min_weight"))
	maxWeight := parseIntPtr(r.URL.Query().Get("max_weight"))

	color := strPtr(r.URL.Query().Get("color"))
	size := strPtr(r.URL.Query().Get("size"))

	if categoryID == nil && (color != nil || size != nil) {
		response.BadRequest(w, errors.New("Filters 'color' and 'size' are only allowed when 'category_id' is provided."))
		return
	}

	inStock := boolPtr(r.URL.Query().Get("in_stock"))
	hasDiscount := boolPtr(r.URL.Query().Get("has_retail_discount"))
	onlyWholesale := boolPtr(r.URL.Query().Get("only_wholesale"))
	sort := strPtr(r.URL.Query().Get("sort"))

	// 2️⃣ Build service input
	input := service.FeedQueryInput{
		Page:              page,
		PerPage:           perPage,
		Query:             strPtr(query),
		CategoryID:        categoryID,
		MinPrice:          minPrice,
		MaxPrice:          maxPrice,
		MinWeight:         minWeight,
		MaxWeight:         maxWeight,
		Color:             color,
		Size:              size,
		InStock:           inStock,
		HasRetailDiscount: hasDiscount,
		OnlyWholesale:     onlyWholesale,
		Sort:              sort,
	}

	// 3️⃣ Call service
	result, err := h.Service.Start(ctx, input)
	if err != nil {
		response.ServerError(w, err)
		return
	}

	// 4️⃣ Return response
	response.OK(w, "Products fetched successfully", result)
}

func parseInt(s string, fallback int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return i
}

func parseInt64Ptr(s string) *int64 {
	if s == "" {
		return nil
	}
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil
	}
	return &val
}

func parseIntPtr(s string) *int {
	if s == "" {
		return nil
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &val
}

func boolPtr(s string) *bool {
	if s == "" {
		return nil
	}
	if s == "true" {
		b := true
		return &b
	}
	if s == "false" {
		b := false
		return &b
	}
	return nil
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
