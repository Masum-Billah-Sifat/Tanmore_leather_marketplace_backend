// 🧠 Amazon-style Dynamic Feed/Search Query Builder
// package feedquery

// import (
// 	"fmt"
// 	"strings"
// )

// // 📦 FeedQueryParams defines the optional filters from the handler/service
// type FeedQueryParams struct {
// 	Q                 *string
// 	CategoryID        *string
// 	MinPrice          *int64
// 	MaxPrice          *int64
// 	MinWeight         *int
// 	MaxWeight         *int
// 	Color             *string
// 	Size              *string
// 	InStock           *bool
// 	HasRetailDiscount *bool
// 	OnlyWholesale     *bool
// 	Sort              *string
// 	Limit             int
// 	Offset            int
// }

// // 🏗️ Builds the final dynamic query + args
// func BuildDynamicFeedQuery(q FeedQueryParams) (string, []interface{}) {
// 	base := `
// SELECT *, (
//   ts_rank_cd(search_vector, websearch_to_tsquery('simple', $1)) * 1.5 +
//   log(views + 1) * 1.2 +
//   CASE WHEN has_retail_discount THEN 1 ELSE 0 END +
//   CASE WHEN isvariantinstock THEN 0.5 ELSE 0 END +
//   CASE WHEN now() - createdat < interval '7 days' THEN 0.3 ELSE 0 END
// ) AS relevance_score
// FROM product_variant_indexes
// WHERE
//   isproductapproved = true AND isproductarchived = false AND isproductbanned = false AND
//   isvariantarchived = false AND issellerapproved = true AND issellerarchived = false AND
//   issellerbanned = false AND iscategoryarchived = false`

// 	clauses := []string{}
// 	args := []interface{}{defaultOr(q.Q, "")}
// 	argID := 2

// 	add := func(condition string, value interface{}) {
// 		clauses = append(clauses, fmt.Sprintf(condition, argID))
// 		args = append(args, value)
// 		argID++
// 	}

// 	if q.CategoryID != nil {
// 		add("categoryid = $%d", *q.CategoryID)
// 	}
// 	if q.MinPrice != nil {
// 		add("retailprice >= $%d", *q.MinPrice)
// 	}
// 	if q.MaxPrice != nil {
// 		add("retailprice <= $%d", *q.MaxPrice)
// 	}
// 	if q.MinWeight != nil {
// 		add("weight_grams >= $%d", *q.MinWeight)
// 	}
// 	if q.MaxWeight != nil {
// 		add("weight_grams <= $%d", *q.MaxWeight)
// 	}
// 	if q.Color != nil {
// 		add("color = $%d", *q.Color)
// 	}
// 	if q.Size != nil {
// 		add("size = $%d", *q.Size)
// 	}
// 	if q.InStock != nil {
// 		add("isvariantinstock = $%d", *q.InStock)
// 	}
// 	if q.HasRetailDiscount != nil {
// 		add("has_retail_discount = $%d", *q.HasRetailDiscount)
// 	}
// 	if q.OnlyWholesale != nil {
// 		add("haswholesaleenabled = $%d", *q.OnlyWholesale)
// 	}

// 	// 🔀 Sort logic
// 	sortClause := "ORDER BY createdat DESC"
// 	if q.Sort != nil {
// 		switch *q.Sort {
// 		case "price_asc":
// 			sortClause = "ORDER BY retailprice ASC"
// 		case "price_desc":
// 			sortClause = "ORDER BY retailprice DESC"
// 		case "most_viewed":
// 			sortClause = "ORDER BY views DESC"
// 		case "newest":
// 			sortClause = "ORDER BY createdat DESC"
// 		case "relevance":
// 			if q.Q != nil && strings.TrimSpace(*q.Q) != "" {
// 				sortClause = "ORDER BY relevance_score DESC"
// 			}
// 		}
// 	}

// 	// 🔗 Final WHERE clause
// 	where := ""
// 	if len(clauses) > 0 {
// 		where = " AND " + strings.Join(clauses, " AND ")
// 	}

// 	// 📦 Final full query
// 	final := fmt.Sprintf("%s%s %s LIMIT $%d OFFSET $%d", base, where, sortClause, argID, argID+1)
// 	args = append(args, q.Limit, q.Offset)

// 	return final, args
// }

// func defaultOr(val *string, fallback string) string {
// 	if val != nil {
// 		return *val
// 	}
// 	return fallback
// }

package feedquery

import (
	"fmt"
	"strings"

	"github.com/google/uuid" // ✅ Required for UUID type
)

// 📦 FeedQueryParams with safer types
type FeedQueryParams struct {
	Q                 *string
	CategoryID        *uuid.UUID // ✅ FIX: was *string → now UUID-safe
	MinPrice          *int64
	MaxPrice          *int64
	MinWeight         *int64 // ✅ FIX: was *int → now *int64 to match BIGINT
	MaxWeight         *int64 // ✅ FIX: was *int → now *int64 to match BIGINT
	Color             *string
	Size              *string
	InStock           *bool
	HasRetailDiscount *bool
	OnlyWholesale     *bool
	Sort              *string
	Limit             int
	Offset            int
}

// // 🧠 Dynamic Feed/Search Query Builder
// func BuildDynamicFeedQuery(q FeedQueryParams) (string, []interface{}) {
// 	// ✅ FIX: Only include FTS if query is non-empty
// 	useFTS := q.Q != nil && strings.TrimSpace(*q.Q) != ""

// 	args := []interface{}{}
// 	argID := 1

// 	// 🔨 Building base SELECT and dynamic scoring logic
// 	base := `
// SELECT *, (
// `
// 	if useFTS {
// 		base += fmt.Sprintf(`
//   ts_rank_cd(search_vector, websearch_to_tsquery('simple', $%d)) * 1.5 +`, argID)
// 		args = append(args, *q.Q)
// 		argID++
// 	}

// 	base += `
//   log(views + 1) * 1.2 +
//   CASE WHEN has_retail_discount THEN 1 ELSE 0 END +
//   CASE WHEN isvariantinstock THEN 0.5 ELSE 0 END +
//   CASE WHEN now() - createdat < interval '7 days' THEN 0.3 ELSE 0 END
// ) AS relevance_score
// FROM product_variant_indexes
// WHERE
//   isproductapproved = true AND
//   isproductarchived = false AND
//   isproductbanned = false AND
//   isvariantarchived = false AND
//   issellerapproved = true AND
//   issellerarchived = false AND
//   issellerbanned = false AND
//   iscategoryarchived = false`

// 	// 📌 Additional WHERE clauses (filters)
// 	clauses := []string{}

// 	add := func(condition string, value interface{}) {
// 		clauses = append(clauses, fmt.Sprintf(condition, argID))
// 		args = append(args, value)
// 		argID++
// 	}

// 	// ✅ FIX: Proper UUID handling
// 	if q.CategoryID != nil {
// 		add("categoryid = $%d", *q.CategoryID)
// 	}

// 	if q.MinPrice != nil {
// 		add("retailprice >= $%d", *q.MinPrice)
// 	}
// 	if q.MaxPrice != nil {
// 		add("retailprice <= $%d", *q.MaxPrice)
// 	}
// 	if q.MinWeight != nil {
// 		add("weight_grams >= $%d", *q.MinWeight)
// 	}
// 	if q.MaxWeight != nil {
// 		add("weight_grams <= $%d", *q.MaxWeight)
// 	}

// 	// ✅ FIX: Case-insensitive matching
// 	if q.Color != nil {
// 		add("LOWER(color) = LOWER($%d)", *q.Color)
// 	}
// 	if q.Size != nil {
// 		add("LOWER(size) = LOWER($%d)", *q.Size)
// 	}

// 	if q.InStock != nil {
// 		add("isvariantinstock = $%d", *q.InStock)
// 	}
// 	if q.HasRetailDiscount != nil {
// 		add("has_retail_discount = $%d", *q.HasRetailDiscount)
// 	}

// 	// ✅ FIX: Only filter when true, ignore false
// 	if q.OnlyWholesale != nil && *q.OnlyWholesale {
// 		clauses = append(clauses, "haswholesaleenabled = true")
// 	}

// 	// 🔀 Sort logic with FTS fallback
// 	sortClause := "ORDER BY createdat DESC"
// 	if q.Sort != nil {
// 		switch *q.Sort {
// 		case "price_asc":
// 			sortClause = "ORDER BY retailprice ASC"
// 		case "price_desc":
// 			sortClause = "ORDER BY retailprice DESC"
// 		case "most_viewed":
// 			sortClause = "ORDER BY views DESC"
// 		case "newest":
// 			sortClause = "ORDER BY createdat DESC"
// 		case "relevance":
// 			if useFTS {
// 				sortClause = "ORDER BY relevance_score DESC"
// 			}
// 		}
// 	}

// 	// 🔗 Final WHERE clause
// 	where := ""
// 	if len(clauses) > 0 {
// 		where = " AND " + strings.Join(clauses, " AND ")
// 	}

// 	// 🧾 Final query with LIMIT/OFFSET
// 	final := fmt.Sprintf("%s%s %s LIMIT $%d OFFSET $%d", base, where, sortClause, argID, argID+1)
// 	args = append(args, q.Limit, q.Offset)

// 	// ✅ Optional: Debug logging
// 	fmt.Println("FINAL QUERY:")
// 	fmt.Println(final)
// 	fmt.Println("ARGS:")
// 	for i, a := range args {
// 		fmt.Printf("  $%d = %#v\n", i+1, a)
// 	}

// 	return final, args
// }

// func BuildDynamicFeedQuery(q FeedQueryParams) (string, []interface{}) {
// 	useFTS := q.Q != nil && strings.TrimSpace(*q.Q) != ""

// 	args := []interface{}{}
// 	argID := 1

// 	base := `
// SELECT *, (
// `
// 	if useFTS {
// 		base += fmt.Sprintf(`
//   ts_rank_cd(search_vector, websearch_to_tsquery('simple', $%d)) * 1.5 +`, argID)
// 		args = append(args, *q.Q)
// 		argID++
// 	}

// 	base += `
//   log(views + 1) * 1.2 +
//   CASE WHEN has_retail_discount THEN 1 ELSE 0 END +
//   CASE WHEN isvariantinstock THEN 0.5 ELSE 0 END +
//   CASE WHEN now() - createdat < interval '7 days' THEN 0.3 ELSE 0 END
// ) AS relevance_score
// FROM product_variant_indexes
// WHERE
//   1=1
// -- 🔍 Temporarily commenting out all base filters for debug
// --   isproductapproved = true AND
// --   isproductarchived = false AND
// --   isproductbanned = false AND
// --   isvariantarchived = false AND
// --   issellerapproved = true AND
// --   issellerarchived = false AND
// --   issellerbanned = false AND
// --   iscategoryarchived = false
// `

// 	clauses := []string{}
// 	add := func(condition string, value interface{}) {
// 		clauses = append(clauses, fmt.Sprintf(condition, argID))
// 		args = append(args, value)
// 		argID++
// 	}

// 	// All dynamic filters stay (we're testing only base clause here)
// 	if q.CategoryID != nil {
// 		add("categoryid = $%d", *q.CategoryID)
// 	}
// 	if q.MinPrice != nil {
// 		add("retailprice >= $%d", *q.MinPrice)
// 	}
// 	if q.MaxPrice != nil {
// 		add("retailprice <= $%d", *q.MaxPrice)
// 	}
// 	if q.MinWeight != nil {
// 		add("weight_grams >= $%d", *q.MinWeight)
// 	}
// 	if q.MaxWeight != nil {
// 		add("weight_grams <= $%d", *q.MaxWeight)
// 	}
// 	if q.Color != nil {
// 		add("LOWER(color) = LOWER($%d)", *q.Color)
// 	}
// 	if q.Size != nil {
// 		add("LOWER(size) = LOWER($%d)", *q.Size)
// 	}
// 	if q.InStock != nil {
// 		add("isvariantinstock = $%d", *q.InStock)
// 	}
// 	if q.HasRetailDiscount != nil {
// 		add("has_retail_discount = $%d", *q.HasRetailDiscount)
// 	}
// 	if q.OnlyWholesale != nil && *q.OnlyWholesale {
// 		clauses = append(clauses, "haswholesaleenabled = true")
// 	}

// 	// Final conditional WHERE
// 	where := ""
// 	if len(clauses) > 0 {
// 		where = " AND " + strings.Join(clauses, " AND ")
// 	}

// 	// Sorting
// 	sortClause := "ORDER BY createdat DESC"
// 	if q.Sort != nil {
// 		switch *q.Sort {
// 		case "price_asc":
// 			sortClause = "ORDER BY retailprice ASC"
// 		case "price_desc":
// 			sortClause = "ORDER BY retailprice DESC"
// 		case "most_viewed":
// 			sortClause = "ORDER BY views DESC"
// 		case "newest":
// 			sortClause = "ORDER BY createdat DESC"
// 		case "relevance":
// 			if useFTS {
// 				sortClause = "ORDER BY relevance_score DESC"
// 			}
// 		}
// 	}

// 	// Final Query
// 	final := fmt.Sprintf("%s%s %s LIMIT $%d OFFSET $%d", base, where, sortClause, argID, argID+1)
// 	args = append(args, q.Limit, q.Offset)

// 	// Logging
// 	fmt.Println("FINAL QUERY:")
// 	fmt.Println(final)
// 	fmt.Println("ARGS:")
// 	for i, a := range args {
// 		fmt.Printf("  $%d = %#v\n", i+1, a)
// 	}

// 	return final, args
// }

// func BuildDynamicFeedQuery(q FeedQueryParams) (string, []interface{}) {
// 	useFTS := q.Q != nil && strings.TrimSpace(*q.Q) != ""

// 	args := []interface{}{}
// 	argID := 1

// 	// ✅ Explicit select of all fields (matching schema + struct field order)
// 	base := `
// SELECT
// 	id,
// 	categoryid,
// 	iscategoryarchived,
// 	categoryname,

// 	sellerid,
// 	issellerapproved,
// 	issellerarchived,
// 	issellerbanned,
// 	sellerstorename,

// 	productid,
// 	isproductapproved,
// 	isproductarchived,
// 	isproductbanned,
// 	producttitle,
// 	productdescription,
// 	productimages,
// 	productpromovideourl,

// 	variantid,
// 	isvariantarchived,
// 	isvariantinstock,
// 	stockamount,
// 	color,
// 	size,
// 	retailprice,
// 	retaildiscounttype,
// 	retaildiscount,
// 	has_retail_discount,

// 	haswholesaleenabled,
// 	wholesaleprice,
// 	wholesaleminquantity,
// 	wholesalediscounttype,
// 	wholesalediscount,

// 	weight_grams,
// 	views,
// `

// 	// ✅ Add FTS relevance score if needed
// 	if useFTS {
// 		base += fmt.Sprintf(`
// 	ts_rank_cd(search_vector, websearch_to_tsquery('simple', $%d)) AS relevance_score,`, argID)
// 		args = append(args, *q.Q)
// 		argID++
// 	} else {
// 		base += `NULL AS relevance_score,`
// 	}

// 	base += `
// 	createdat,
// 	updatedat
// FROM product_variant_indexes
// WHERE 1=1
// -- ✅ All base filters temporarily commented out for debug
// --  AND isproductapproved = true
// --  AND isproductarchived = false
// --  AND isproductbanned = false
// --  AND isvariantarchived = false
// --  AND issellerapproved = true
// --  AND issellerarchived = false
// --  AND issellerbanned = false
// --  AND iscategoryarchived = false
// `

// 	// 📌 Dynamic filter building
// 	clauses := []string{}
// 	add := func(condition string, value interface{}) {
// 		clauses = append(clauses, fmt.Sprintf(condition, argID))
// 		args = append(args, value)
// 		argID++
// 	}

// 	if q.CategoryID != nil {
// 		add("categoryid = $%d", *q.CategoryID)
// 	}
// 	if q.MinPrice != nil {
// 		add("retailprice >= $%d", *q.MinPrice)
// 	}
// 	if q.MaxPrice != nil {
// 		add("retailprice <= $%d", *q.MaxPrice)
// 	}
// 	if q.MinWeight != nil {
// 		add("weight_grams >= $%d", *q.MinWeight)
// 	}
// 	if q.MaxWeight != nil {
// 		add("weight_grams <= $%d", *q.MaxWeight)
// 	}
// 	if q.Color != nil {
// 		add("LOWER(color) = LOWER($%d)", *q.Color)
// 	}
// 	if q.Size != nil {
// 		add("LOWER(size) = LOWER($%d)", *q.Size)
// 	}
// 	if q.InStock != nil {
// 		add("isvariantinstock = $%d", *q.InStock)
// 	}
// 	if q.HasRetailDiscount != nil {
// 		add("has_retail_discount = $%d", *q.HasRetailDiscount)
// 	}
// 	if q.OnlyWholesale != nil && *q.OnlyWholesale {
// 		clauses = append(clauses, "haswholesaleenabled = true")
// 	}

// 	where := ""
// 	if len(clauses) > 0 {
// 		where = " AND " + strings.Join(clauses, " AND ")
// 	}

// 	// 🔀 Sort logic
// 	sortClause := "ORDER BY createdat DESC"
// 	if q.Sort != nil {
// 		switch *q.Sort {
// 		case "price_asc":
// 			sortClause = "ORDER BY retailprice ASC"
// 		case "price_desc":
// 			sortClause = "ORDER BY retailprice DESC"
// 		case "most_viewed":
// 			sortClause = "ORDER BY views DESC"
// 		case "newest":
// 			sortClause = "ORDER BY createdat DESC"
// 		case "relevance":
// 			if useFTS {
// 				sortClause = "ORDER BY relevance_score DESC"
// 			}
// 		}
// 	}

// 	// 🔚 Final query
// 	final := fmt.Sprintf("%s%s %s LIMIT $%d OFFSET $%d", base, where, sortClause, argID, argID+1)
// 	args = append(args, q.Limit, q.Offset)

// 	// 🪵 Debug logging
// 	fmt.Println("FINAL QUERY:")
// 	fmt.Println(final)
// 	fmt.Println("ARGS:")
// 	for i, a := range args {
// 		fmt.Printf("  $%d = %#v\n", i+1, a)
// 	}

// 	return final, args
// }

func BuildDynamicFeedQuery(q FeedQueryParams) (string, []interface{}) {
	useFTS := q.Q != nil && strings.TrimSpace(*q.Q) != ""

	args := []interface{}{}
	argID := 1

	// ✅ Explicit SELECT fields matching schema and struct
	base := `
SELECT
	id,
	categoryid,
	iscategoryarchived,
	categoryname,

	sellerid,
	issellerapproved,
	issellerarchived,
	issellerbanned,
	sellerstorename,

	productid,
	isproductapproved,
	isproductarchived,
	isproductbanned,
	producttitle,
	productdescription,
	productimages,
	productpromovideourl,

	variantid,
	isvariantarchived,
	isvariantinstock,
	stockamount,
	color,
	size,
	retailprice,
	retaildiscounttype,
	retaildiscount,
	has_retail_discount,

	haswholesaleenabled,
	wholesaleprice,
	wholesaleminquantity,
	wholesalediscounttype,
	wholesalediscount,

	weight_grams,
	views,
`

	// ✅ FTS relevance score if search query exists
	if useFTS {
		base += fmt.Sprintf(`
	ts_rank_cd(search_vector, websearch_to_tsquery('simple', $%d)) AS relevance_score,`, argID)
		args = append(args, *q.Q)
		argID++
	} else {
		base += `NULL AS relevance_score,`
	}

	base += `
	createdat,
	updatedat
FROM product_variant_indexes
WHERE
  isproductapproved = true AND
  isproductarchived = false AND
  isproductbanned = false AND
  isvariantarchived = false AND
  issellerapproved = true AND
  issellerarchived = false AND
  issellerbanned = false AND
  iscategoryarchived = false
`

	// 📌 Dynamic filters
	clauses := []string{}
	add := func(condition string, value interface{}) {
		clauses = append(clauses, fmt.Sprintf(condition, argID))
		args = append(args, value)
		argID++
	}

	if q.CategoryID != nil {
		add("categoryid = $%d", *q.CategoryID)
	}
	if q.MinPrice != nil {
		add("retailprice >= $%d", *q.MinPrice)
	}
	if q.MaxPrice != nil {
		add("retailprice <= $%d", *q.MaxPrice)
	}
	if q.MinWeight != nil {
		add("weight_grams >= $%d", *q.MinWeight)
	}
	if q.MaxWeight != nil {
		add("weight_grams <= $%d", *q.MaxWeight)
	}
	if q.Color != nil {
		add("LOWER(color) = LOWER($%d)", *q.Color)
	}
	if q.Size != nil {
		add("LOWER(size) = LOWER($%d)", *q.Size)
	}
	if q.InStock != nil {
		add("isvariantinstock = $%d", *q.InStock)
	}
	if q.HasRetailDiscount != nil {
		add("has_retail_discount = $%d", *q.HasRetailDiscount)
	}
	if q.OnlyWholesale != nil && *q.OnlyWholesale {
		clauses = append(clauses, "haswholesaleenabled = true")
	}

	where := ""
	if len(clauses) > 0 {
		where = " AND " + strings.Join(clauses, " AND ")
	}

	// 🔀 Sorting logic
	sortClause := "ORDER BY createdat DESC"
	if q.Sort != nil {
		switch *q.Sort {
		case "price_asc":
			sortClause = "ORDER BY retailprice ASC"
		case "price_desc":
			sortClause = "ORDER BY retailprice DESC"
		case "most_viewed":
			sortClause = "ORDER BY views DESC"
		case "newest":
			sortClause = "ORDER BY createdat DESC"
		case "relevance":
			if useFTS {
				sortClause = "ORDER BY relevance_score DESC"
			}
		}
	}

	// Final query
	final := fmt.Sprintf("%s%s %s LIMIT $%d OFFSET $%d", base, where, sortClause, argID, argID+1)
	args = append(args, q.Limit, q.Offset)

	// 🪵 Debug log
	fmt.Println("FINAL QUERY:")
	fmt.Println(final)
	fmt.Println("ARGS:")
	for i, a := range args {
		fmt.Printf("  $%d = %#v\n", i+1, a)
	}

	return final, args
}
