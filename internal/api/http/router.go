package http

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	dfmiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"

	"tanmore_backend/internal/api/http/middleware" // adjust import path if needed

	repo_googleauth "tanmore_backend/internal/repository/google_auth"
	repo_token_refresh "tanmore_backend/internal/repository/token_refresh"

	google_auth_handlers "tanmore_backend/internal/api/http/handlers/google_auth"
	google_auth_services "tanmore_backend/internal/services/google_auth" // same package

	repo_switchmode "tanmore_backend/internal/repository/user_mode_switch"

	switchmode_handlers "tanmore_backend/internal/api/http/handlers/user_switch_mode"
	switchmode_services "tanmore_backend/internal/services/user_mode_switch"

	// Product creation imports
	product_handlers "tanmore_backend/internal/api/http/handlers/product"
	repo_product "tanmore_backend/internal/repository/product/product_creation"
	product_services "tanmore_backend/internal/services/product"

	// Product variant addition
	variant_handlers "tanmore_backend/internal/api/http/handlers/product/product_variant"
	repo_variant "tanmore_backend/internal/repository/product/product_variant"
	variant_services "tanmore_backend/internal/services/product/product_variant"

	repo_variant_archive "tanmore_backend/internal/repository/product/product_variant/product_variant_archive"
	repo_variant_update_info "tanmore_backend/internal/repository/product/product_variant/product_variant_update_info"
	repo_variant_update_price "tanmore_backend/internal/repository/product/product_variant/product_variant_update_price"

	repo_variant_update_in_stock "tanmore_backend/internal/repository/product/product_variant/product_variant_update_in_stock"

	repo_variant_update_stock "tanmore_backend/internal/repository/product/product_variant/product_variant_update_stock_quantity"

	// ⬇️ Add below existing import blocks
	repo_variant_update_weight "tanmore_backend/internal/repository/product/product_variant/product_variant_update_weight"

	repo_variant_add_discount "tanmore_backend/internal/repository/product/product_variant/product_variant_add_discount"

	repo_variant_update_discount "tanmore_backend/internal/repository/product/product_variant/product_variant_update_discount"

	repo_variant_remove_discount "tanmore_backend/internal/repository/product/product_variant/product_variant_remove_retail_discount"

	repo_variant_enable_wholesale "tanmore_backend/internal/repository/product/product_variant/product_variant_enable_wholesale_mode"

	// ⬇️ Wholesale mode (edit)
	repo_variant_edit_wholesale "tanmore_backend/internal/repository/product/product_variant/product_variant_update_wholesale_mode"

	// ➖ Disable Wholesale Mode
	repo_variant_disable_wholesale "tanmore_backend/internal/repository/product/product_variant/product_variant_disable_wholesale"

	repo_variant_add_wholesale_discount "tanmore_backend/internal/repository/product/product_variant/product_variant_add_wholesale_discount"

	repo_variant_update_wholesale_discount "tanmore_backend/internal/repository/product/product_variant/product_variant_update_wholesale_discount"

	repo_variant_remove_wholesale_discount "tanmore_backend/internal/repository/product/product_variant/product_variant_disable_wholesale_discount"

	// Feed or Search endpoint (public)
	feed_search_handlers "tanmore_backend/internal/api/http/handlers/feed_query"
	repo_feed_search "tanmore_backend/internal/repository/product/product_variant/product_variant_index_feed_or_search"
	feed_search_services "tanmore_backend/internal/services/product/product_variant/product_variant_index_feed_or_search"

	// ✅ Seller profile creation
	seller_profile_handlers "tanmore_backend/internal/api/http/handlers/seller_profile"
	repo_seller_profile "tanmore_backend/internal/repository/seller_profile/seller_profile_metadata"
	seller_profile_services "tanmore_backend/internal/services/seller_profile"

	// 🆕 Update Product Info
	update_product_info_handlers "tanmore_backend/internal/api/http/handlers/product"
	repo_update_product_info "tanmore_backend/internal/repository/product/product_update_info"
	update_product_info_services "tanmore_backend/internal/services/product"

	// 🆕 Add Product Media
	add_media_handlers "tanmore_backend/internal/api/http/handlers/product"
	repo_add_media "tanmore_backend/internal/repository/product/product_add_media"
	add_media_services "tanmore_backend/internal/services/product"

	// 🔽 Add below existing product imports
	archive_media_handlers "tanmore_backend/internal/api/http/handlers/product"
	repo_archive_media "tanmore_backend/internal/repository/product/product_archive_media"
	archive_media_services "tanmore_backend/internal/services/product"

	// 🖼️ Set Primary Image
	set_primary_image_handlers "tanmore_backend/internal/api/http/handlers/product"
	repo_set_primary_image "tanmore_backend/internal/repository/product/product_set_primary_image"
	set_primary_image_services "tanmore_backend/internal/services/product"

	cart_handlers "tanmore_backend/internal/api/http/handlers/cart"
	cart_repo "tanmore_backend/internal/repository/cart/add_to_cart"
	cart_services "tanmore_backend/internal/services/cart"

	// 🛒 Update Cart Quantity
	update_quantity_handlers "tanmore_backend/internal/api/http/handlers/cart"
	update_quantity_repo "tanmore_backend/internal/repository/cart/update_to_cart"
	update_quantity_service "tanmore_backend/internal/services/cart"

	// 🛒 Remove Cart Item
	remove_cart_handlers "tanmore_backend/internal/api/http/handlers/cart"
	remove_cart_repo "tanmore_backend/internal/repository/cart/remove_from_cart"
	remove_cart_service "tanmore_backend/internal/services/cart"

	// 🧹 Clear Cart Items
	clear_cart_handlers "tanmore_backend/internal/api/http/handlers/cart"
	clear_cart_repo "tanmore_backend/internal/repository/cart/clear_cart"
	clear_cart_service "tanmore_backend/internal/services/cart"

	get_all_cart_items_handlers "tanmore_backend/internal/api/http/handlers/cart"
	get_all_cart_items_repo "tanmore_backend/internal/repository/cart/get_all_cart_items"
	get_all_cart_items_service "tanmore_backend/internal/services/cart"

	cart_summary_handlers "tanmore_backend/internal/api/http/handlers/cart"
	cart_summary_repo "tanmore_backend/internal/repository/cart/cart_summary"
	cart_summary_service "tanmore_backend/internal/services/cart"

	// 🧾 Checkout Endpoint
	checkout_handlers "tanmore_backend/internal/api/http/handlers/checkout"
	checkout_repo "tanmore_backend/internal/repository/checkout"
	checkout_services "tanmore_backend/internal/services/checkout"

	// 🧭 Category Tree
	category_tree_handlers "tanmore_backend/internal/api/http/handlers/category"
	category_tree_repo "tanmore_backend/internal/repository/category/category_tree"
	category_tree_services "tanmore_backend/internal/services/category"

	// 🛍️ Fetch products by category
	fetch_by_category_handlers "tanmore_backend/internal/api/http/handlers/product"
	fetch_by_category_repo "tanmore_backend/internal/repository/product/fetch_by_category"
	fetch_by_category_service "tanmore_backend/internal/services/product"

	// 🧹 Archive product
	archive_product_handlers "tanmore_backend/internal/api/http/handlers/product"
	archive_product_repo "tanmore_backend/internal/repository/product/product_archive"
	archive_product_services "tanmore_backend/internal/services/product"

	// Product Review endpoint
	review_handlers "tanmore_backend/internal/api/http/handlers/product"
	repo_add_review "tanmore_backend/internal/repository/product/product_add_review"
	review_services "tanmore_backend/internal/services/product"

	// ✏️ Edit Review
	edit_review_handler "tanmore_backend/internal/api/http/handlers/product"
	repo_edit_review "tanmore_backend/internal/repository/product/product_edit_review"
	edit_review_service "tanmore_backend/internal/services/product"

	// Archive Review
	archive_review_handlers "tanmore_backend/internal/api/http/handlers/product"
	repo_archive_review "tanmore_backend/internal/repository/product/product_archive_review"
	archive_review_services "tanmore_backend/internal/services/product"

	// 🗨️ Reply to Review
	handler_reply_review "tanmore_backend/internal/api/http/handlers/product"
	repo_reply_review "tanmore_backend/internal/repository/product/product_reply_review"
	service_reply_review "tanmore_backend/internal/services/product"

	// ✏️ Edit Review Reply (Seller)
	edit_review_reply_handler "tanmore_backend/internal/api/http/handlers/product"
	repo_edit_review_reply "tanmore_backend/internal/repository/product/product_review_reply_edit"
	edit_review_reply_service "tanmore_backend/internal/services/product"

	// 🗑️ Archive Review Reply (Seller Only)
	archive_reply_handler "tanmore_backend/internal/api/http/handlers/product"
	archive_reply_repo "tanmore_backend/internal/repository/product/product_review_reply_archive"
	archive_reply_service "tanmore_backend/internal/services/product"

	// Reviews: Get All Product Reviews with Replies
	reviews_handlers "tanmore_backend/internal/api/http/handlers/product"
	repo_get_reviews "tanmore_backend/internal/repository/product/product_get_all_reviews"
	reviews_services "tanmore_backend/internal/services/product"

	// 🛑 Logout Session
	logout_handlers "tanmore_backend/internal/api/http/handlers/logout"
	repo_logout "tanmore_backend/internal/repository/logout"
	logout_services "tanmore_backend/internal/services/logout"

	// 🎥 Media Upload Presign
	media_handlers "tanmore_backend/internal/api/http/handlers/media"
	media_services "tanmore_backend/internal/services/media"

	// 🆕 Product Full Detail (Seller Side)
	product_full_detail_handlers "tanmore_backend/internal/api/http/handlers/product"
	repo_product_full_detail "tanmore_backend/internal/repository/product/product_get_full_detail"
	product_full_detail_services "tanmore_backend/internal/services/product"

	// 🆕 Fetch All Products for Seller (Grouped by moderation state)
	get_all_products_handler "tanmore_backend/internal/api/http/handlers/product"
	repo_get_all_products_seller "tanmore_backend/internal/repository/product/product_get_all_grouped"
	get_all_products_service "tanmore_backend/internal/services/product"

	// 🆕 Update Product Category
	update_category_handlers "tanmore_backend/internal/api/http/handlers/product"
	repo_update_category "tanmore_backend/internal/repository/product/product_update_category"
	update_category_services "tanmore_backend/internal/services/product"

	// 🌐 Public Product Detail Endpoint
	public_detail_handlers "tanmore_backend/internal/api/http/handlers/product"
	public_detail_repo "tanmore_backend/internal/repository/product/product_get_public_detail"
	public_detail_service "tanmore_backend/internal/services/product"

	// 🧾 Checkout session detail
	// checkout_handlers "tanmore_backend/internal/api/http/handlers/checkout"
	repo_checkout_session "tanmore_backend/internal/repository/checkout/session_details"
	// checkout_services "tanmore_backend/internal/services/checkout"

	// 🆕 Merge Guest Cart
	merge_cart_handlers "tanmore_backend/internal/api/http/handlers/cart"
	merge_cart_repo "tanmore_backend/internal/repository/cart/merge_guest_cart"
	merge_cart_services "tanmore_backend/internal/services/cart"

	// 🆕 Add Shipping Address to Checkout Session
	// checkout_handlers "tanmore_backend/internal/api/http/handlers/checkout"
	checkout_shipping_address_add_repo "tanmore_backend/internal/repository/checkout/add_shipping_address"
	checkout_service "tanmore_backend/internal/services/checkout"

	// 🧾 Checkout – Shipping Address
	// checkout_handlers "tanmore_backend/internal/api/http/handlers/checkout"
	checkout_edit_repo "tanmore_backend/internal/repository/checkout/edit_shipping_address"
	// checkout_service "tanmore_backend/internal/services/checkout"

	// checkout_handlers "tanmore_backend/internal/api/http/handlers/checkout"
	checkout_review_repo "tanmore_backend/internal/repository/checkout/review"
	// checkout_services "tanmore_backend/internal/services/checkout"

	// Checkout Confirm Order
	confirm_order_handlers "tanmore_backend/internal/api/http/handlers/checkout"
	confirm_order_repo "tanmore_backend/internal/repository/checkout/confirm_order"
	confirm_order_services "tanmore_backend/internal/services/checkout"

	"tanmore_backend/pkg/token"
)

func NewRouter(db *sql.DB, redisClient *redis.Client) http.Handler {
	r := chi.NewRouter()

	// 🌐 Global middlewares
	r.Use(dfmiddleware.RequestID)
	r.Use(dfmiddleware.RealIP)
	r.Use(dfmiddleware.Logger)
	r.Use(dfmiddleware.Recoverer)

	// ✅ 🔥 Add this line to enable CORS!
	r.Use(middleware.CORSMiddleware)

	// ⚙️ Google Login related stuff
	googleAuthRepo := repo_googleauth.NewGoogleAuthRepository(db)
	tokenRefreshRepo := repo_token_refresh.NewTokenRefreshRepository(db)

	// Login Handler
	googleAuthService := google_auth_services.NewGoogleAuthService(google_auth_services.GoogleAuthServiceDeps{
		Repo: googleAuthRepo,
	})
	googleAuthHandler := google_auth_handlers.NewHandler(googleAuthService)

	// Refresh Token Handler
	refreshTokenService := google_auth_services.NewRefreshTokenService(tokenRefreshRepo)
	refreshTokenHandler := google_auth_handlers.NewRefreshTokenHandler(refreshTokenService)

	// 🛑 Logout Handler Setup
	logoutRepo := repo_logout.NewLogoutRepository(db)
	logoutService := logout_services.NewLogoutService(logout_services.LogoutServiceDeps{
		Repo: logoutRepo,
	})
	logoutHandler := logout_handlers.NewHandler(logoutService)

	// 🔁 Switch Mode Setup
	switchModeRepo := repo_switchmode.NewUserModeSwitchRepository(db)
	switchModeService := switchmode_services.NewSwitchModeService(switchModeRepo)
	switchModeHandler := switchmode_handlers.NewSwitchModeHandler(switchModeService)

	// 🛍️ Product Creation Setup
	productRepo := repo_product.NewProductRepository(db)
	productService := product_services.NewCreateProductService(product_services.CreateProductServiceDeps{
		Repo: productRepo,
	})
	productHandler := product_handlers.NewCreateProductHandler(productService)

	// ➕ Add Product Variant Setup
	productVariantRepo := repo_variant.NewProductVariantRepository(db)
	addVariantService := variant_services.NewAddProductVariantService(variant_services.AddProductVariantServiceDeps{
		Repo: productVariantRepo,
	})
	addVariantHandler := variant_handlers.NewAddProductVariantHandler(addVariantService)

	// ➖ Remove Product Variant Setup
	productVariantArchiveRepo := repo_variant_archive.NewProductVariantArchiveRepository(db)
	removeVariantService := variant_services.NewRemoveProductVariantService(
		variant_services.RemoveProductVariantServiceDeps{
			Repo: productVariantArchiveRepo,
		},
	)
	removeVariantHandler := variant_handlers.NewRemoveProductVariantHandler(removeVariantService)

	// 📝 Update Variant Info Setup
	productVariantUpdateInfoRepo := repo_variant_update_info.NewProductVariantUpdateInfoRepository(db)

	updateVariantInfoService := variant_services.NewUpdateVariantInfoService(variant_services.UpdateVariantInfoServiceDeps{
		Repo: productVariantUpdateInfoRepo,
	})
	updateVariantInfoHandler := variant_handlers.NewUpdateVariantInfoHandler(updateVariantInfoService)

	// 💵 Update Variant Retail Price Setup
	productVariantUpdatePriceRepo := repo_variant_update_price.NewProductVariantUpdatePriceRepository(db)
	updateRetailPriceService := variant_services.NewUpdateVariantRetailPriceService(variant_services.UpdateVariantRetailPriceServiceDeps{
		Repo: productVariantUpdatePriceRepo,
	})
	updateRetailPriceHandler := variant_handlers.NewUpdateVariantRetailPriceHandler(updateRetailPriceService)

	// 📦 Update Variant In-Stock Setup
	productVariantUpdateInStockRepo := repo_variant_update_in_stock.NewProductVariantUpdateInStockRepository(db)
	updateInStockService := variant_services.NewUpdateVariantInStockService(variant_services.UpdateVariantInStockServiceDeps{
		Repo: productVariantUpdateInStockRepo,
	})
	updateInStockHandler := variant_handlers.NewUpdateVariantInStockHandler(updateInStockService)

	// 📦 Update Variant Stock Quantity Setup
	productVariantUpdateStockRepo := repo_variant_update_stock.NewProductVariantUpdateStockQuantityRepository(db)
	updateStockQuantityService := variant_services.NewUpdateVariantStockQuantityService(variant_services.UpdateVariantStockQuantityServiceDeps{
		Repo: productVariantUpdateStockRepo,
	})
	updateStockQuantityHandler := variant_handlers.NewUpdateVariantStockQuantityHandler(updateStockQuantityService)

	// ⚖️ Update Variant Weight Setup
	productVariantUpdateWeightRepo := repo_variant_update_weight.NewProductVariantUpdateWeightRepository(db)
	updateWeightService := variant_services.NewUpdateVariantWeightService(variant_services.UpdateVariantWeightServiceDeps{
		Repo: productVariantUpdateWeightRepo,
	})
	updateWeightHandler := variant_handlers.NewUpdateVariantWeightHandler(updateWeightService)

	// 💸 Add Variant Retail Discount Setup
	productVariantDiscountRepo := repo_variant_add_discount.NewProductVariantAddDiscountRepository(db)
	addRetailDiscountService := variant_services.NewAddVariantRetailDiscountService(variant_services.AddVariantRetailDiscountServiceDeps{
		Repo: productVariantDiscountRepo,
	})
	addRetailDiscountHandler := variant_handlers.NewAddVariantRetailDiscountHandler(addRetailDiscountService)

	// 🔁 Update Variant Retail Discount Setup
	productVariantUpdateDiscountRepo := repo_variant_update_discount.NewProductVariantUpdateDiscountRepository(db)
	updateRetailDiscountService := variant_services.NewUpdateVariantRetailDiscountService(variant_services.UpdateVariantRetailDiscountServiceDeps{
		Repo: productVariantUpdateDiscountRepo,
	})
	updateRetailDiscountHandler := variant_handlers.NewUpdateVariantRetailDiscountHandler(updateRetailDiscountService)

	// ❌ Remove Variant Retail Discount Setup
	productVariantRemoveDiscountRepo := repo_variant_remove_discount.NewProductVariantRemoveDiscountRepository(db)
	removeRetailDiscountService := variant_services.NewRemoveVariantRetailDiscountService(variant_services.RemoveVariantRetailDiscountServiceDeps{
		Repo: productVariantRemoveDiscountRepo,
	})
	removeRetailDiscountHandler := variant_handlers.NewRemoveVariantRetailDiscountHandler(removeRetailDiscountService)

	// 🏷️ Enable Variant Wholesale Mode Setup
	productVariantEnableWholesaleRepo :=
		repo_variant_enable_wholesale.NewProductVariantEnableWholesaleRepository(db)

	enableWholesaleService :=
		variant_services.NewEnableWholesaleModeService(
			variant_services.EnableWholesaleModeServiceDeps{
				Repo: productVariantEnableWholesaleRepo,
			},
		)

	enableWholesaleHandler :=
		variant_handlers.NewEnableVariantWholesaleModeHandler(enableWholesaleService)

		// ✏️ Edit Variant Wholesale Info Setup
	productVariantEditWholesaleRepo :=
		repo_variant_edit_wholesale.NewProductVariantEditWholesaleInfoRepository(db)

	editWholesaleService :=
		variant_services.NewEditWholesaleInfoService(
			variant_services.EditWholesaleInfoServiceDeps{
				Repo: productVariantEditWholesaleRepo,
			},
		)

	editWholesaleHandler :=
		variant_handlers.NewEditVariantWholesaleInfoHandler(editWholesaleService)

		// ➖ Disable Wholesale Mode Setup
	productVariantDisableWholesaleRepo := repo_variant_disable_wholesale.NewProductVariantDisableWholesaleRepository(db)
	disableWholesaleModeService := variant_services.NewDisableWholesaleModeService(variant_services.DisableWholesaleModeServiceDeps{
		Repo: productVariantDisableWholesaleRepo,
	})
	disableWholesaleModeHandler := variant_handlers.NewDisableVariantWholesaleModeHandler(disableWholesaleModeService)

	// ➕ Add Wholesale Discount Setup
	productVariantAddWholesaleDiscountRepo := repo_variant_add_wholesale_discount.NewProductVariantAddWholesaleDiscountRepository(db)
	addWholesaleDiscountService := variant_services.NewAddWholesaleDiscountService(variant_services.AddWholesaleDiscountServiceDeps{
		Repo: productVariantAddWholesaleDiscountRepo,
	})
	addWholesaleDiscountHandler := variant_handlers.NewAddVariantWholesaleDiscountHandler(addWholesaleDiscountService)

	// 🔁 Update Variant Wholesale Discount Setup
	productVariantUpdateWholesaleDiscountRepo := repo_variant_update_wholesale_discount.NewProductVariantUpdateWholesaleDiscountRepository(db)
	updateWholesaleDiscountService := variant_services.NewUpdateWholesaleDiscountService(variant_services.UpdateWholesaleDiscountServiceDeps{
		Repo: productVariantUpdateWholesaleDiscountRepo,
	})
	updateWholesaleDiscountHandler := variant_handlers.NewUpdateVariantWholesaleDiscountHandler(updateWholesaleDiscountService)

	// ❌ Remove Wholesale Discount Setup
	productVariantRemoveWholesaleDiscountRepo := repo_variant_remove_wholesale_discount.NewProductVariantRemoveWholesaleDiscountRepository(db)
	removeWholesaleDiscountService := variant_services.NewRemoveVariantWholesaleDiscountService(variant_services.RemoveVariantWholesaleDiscountServiceDeps{
		Repo: productVariantRemoveWholesaleDiscountRepo,
	})
	removeWholesaleDiscountHandler := variant_handlers.NewRemoveVariantWholesaleDiscountHandler(removeWholesaleDiscountService)

	// 🔍 Public Feed/Search Endpoint Setup
	feedSearchRepo := repo_feed_search.NewProductVariantIndexFeedOrSearchRepository(db)
	feedSearchService := feed_search_services.NewFeedQueryService(feedSearchRepo)
	feedSearchHandler := feed_search_handlers.NewFeedQueryHandler(feedSearchService)

	// 🧾 Seller Profile Metadata Setup
	sellerProfileRepo := repo_seller_profile.NewSellerProfileMetadataRepository(db)
	sellerProfileService := seller_profile_services.NewCreateSellerProfileService(seller_profile_services.CreateSellerProfileServiceDeps{
		Repo: sellerProfileRepo,
	})
	sellerProfileHandler := seller_profile_handlers.NewCreateSellerProfileHandler(sellerProfileService)

	// 🆕 Update Product Info Setup
	productUpdateRepo := repo_update_product_info.NewProductUpdateInfoRepository(db)
	updateProductInfoService := update_product_info_services.NewUpdateProductInfoService(update_product_info_services.UpdateProductInfoServiceDeps{
		Repo: productUpdateRepo,
	})
	updateProductInfoHandler := update_product_info_handlers.NewUpdateProductInfoHandler(updateProductInfoService)

	// 🆕 Add Product Media Setup
	productMediaRepo := repo_add_media.NewProductAddMediaRepository(db)
	addProductMediaService := add_media_services.NewAddProductMediaService(add_media_services.AddProductMediaServiceDeps{
		Repo: productMediaRepo,
	})
	addProductMediaHandler := add_media_handlers.NewAddProductMediaHandler(addProductMediaService)

	// 🏗️ Archive Media - DELETE /api/seller/products/:product_id/media/:media_id
	archiveMediaRepo := repo_archive_media.NewProductArchiveMediaRepository(db)
	archiveMediaService := archive_media_services.NewArchiveProductMediaService(archive_media_services.ArchiveProductMediaServiceDeps{
		Repo: archiveMediaRepo,
	})
	archiveMediaHandler := archive_media_handlers.NewArchiveProductMediaHandler(archiveMediaService)

	// 🖼️ Set Primary Image
	setPrimaryImageRepo := repo_set_primary_image.NewProductSetPrimaryImageRepository(db)
	setPrimaryImageService := set_primary_image_services.NewSetPrimaryImageService(set_primary_image_services.SetPrimaryImageServiceDeps{
		Repo: setPrimaryImageRepo,
	})
	setPrimaryImageHandler := set_primary_image_handlers.NewSetPrimaryImageHandler(setPrimaryImageService)

	// ------------------------------------------------------------
	// 🛒 Add to Cart endpoint wiring
	cartRepo := cart_repo.NewAddToCartRepository(db)
	addToCartService := cart_services.NewAddToCartService(cart_services.AddToCartServiceDeps{
		Repo: cartRepo,
	})
	addToCartHandler := cart_handlers.NewAddToCartHandler(addToCartService)

	// ------------------------------------------------------------
	// 🛒 Update Cart Quantity Endpoint Wiring

	updateCartQuantityRepo := update_quantity_repo.NewUpdateCartQuantityRepository(db)
	updateCartQuantityService := update_quantity_service.NewUpdateCartQuantityService(update_quantity_service.UpdateCartQuantityServiceDeps{
		Repo: updateCartQuantityRepo,
	})
	updateCartQuantityHandler := update_quantity_handlers.NewUpdateCartQuantityHandler(updateCartQuantityService)

	// ------------------------------------------------------------
	// 🗑️ Remove Cart Item Endpoint Wiring
	removeCartRepo := remove_cart_repo.NewRemoveFromCartRepository(db)
	removeCartService := remove_cart_service.NewRemoveCartItemService(remove_cart_service.RemoveCartItemServiceDeps{
		Repo: removeCartRepo,
	})
	removeCartHandler := remove_cart_handlers.NewRemoveCartItemHandler(removeCartService)

	// ------------------------------------------------------------
	// 🧹 Clear Cart Endpoint Wiring
	clearCartRepo := clear_cart_repo.NewClearCartRepository(db)
	clearCartService := clear_cart_service.NewClearCartService(clear_cart_service.ClearCartServiceDeps{
		Repo: clearCartRepo,
	})
	clearCartHandler := clear_cart_handlers.NewClearCartHandler(clearCartService)

	// 🛒 Get All Cart Items Endpoint Wiring
	getAllCartItemsRepo := get_all_cart_items_repo.NewGetAllCartItemsRepository(db)
	getAllCartItemsService := get_all_cart_items_service.NewGetAllCartItemsService(get_all_cart_items_service.GetAllCartItemsServiceDeps{
		Repo: getAllCartItemsRepo,
	})
	getAllCartItemsHandler := get_all_cart_items_handlers.NewGetAllCartItemsHandler(getAllCartItemsService)

	// 🧮 Cart Summary Endpoint Wiring
	cartSummaryRepo := cart_summary_repo.NewCartSummaryRepository(db)
	cartSummaryService := cart_summary_service.NewCartSummaryService(cart_summary_service.CartSummaryServiceDeps{
		Repo: cartSummaryRepo,
	})
	cartSummaryHandler := cart_summary_handlers.NewCartSummaryHandler(cartSummaryService)

	// 🧾 Checkout Endpoint Setup
	checkoutRepo := checkout_repo.NewCheckoutRepository(db)
	checkoutService := checkout_services.NewCheckoutService(checkout_services.CheckoutServiceDeps{
		Repo: checkoutRepo,
	})
	checkoutHandler := checkout_handlers.NewCheckoutHandler(checkoutService)

	// 🧭 Category Tree Setup
	categoryTreeRepo := category_tree_repo.NewCategoryTreeRepository(db)
	categoryTreeService := category_tree_services.NewGetCategoryTreeService(category_tree_services.GetCategoryTreeServiceDeps{
		Repo: categoryTreeRepo,
	})
	categoryTreeHandler := category_tree_handlers.NewGetCategoryTreeHandler(categoryTreeService)

	// 🛍️ Fetch Products by Category Setup
	fetchByCategoryRepo := fetch_by_category_repo.NewFetchByCategoryRepository(db)
	fetchByCategoryService := fetch_by_category_service.NewGetProductsByCategoryService(fetch_by_category_service.GetProductsByCategoryServiceDeps{
		Repo: fetchByCategoryRepo,
	})
	fetchByCategoryHandler := fetch_by_category_handlers.NewGetProductsByCategoryHandler(fetchByCategoryService)

	// 🧹 Archive Product Setup
	archiveProductRepo := archive_product_repo.NewProductArchiveRepository(db)
	archiveProductService := archive_product_services.NewArchiveProductService(archive_product_services.ArchiveProductServiceDeps{
		Repo: archiveProductRepo,
	})
	archiveProductHandler := archive_product_handlers.NewArchiveProductHandler(archiveProductService)

	// 📝 Add Product Review Setup
	reviewRepo := repo_add_review.NewProductAddReviewRepository(db)
	addReviewService := review_services.NewAddProductReviewService(review_services.AddProductReviewServiceDeps{
		Repo: reviewRepo,
	})
	addReviewHandler := review_handlers.NewAddProductReviewHandler(addReviewService)

	// ✏️ Edit Product Review Setup
	editReviewRepo := repo_edit_review.NewProductEditReviewRepository(db)
	editReviewService := edit_review_service.NewEditProductReviewService(edit_review_service.EditProductReviewServiceDeps{
		Repo: editReviewRepo,
	})
	editReviewHandler := edit_review_handler.NewEditProductReviewHandler(editReviewService)

	// 🗑️ Archive Product Review Setup
	archiveReviewRepo := repo_archive_review.NewProductArchiveReviewRepository(db)
	archiveReviewService := archive_review_services.NewArchiveProductReviewService(archive_review_services.ArchiveProductReviewServiceDeps{
		Repo: archiveReviewRepo,
	})
	archiveReviewHandler := archive_review_handlers.NewArchiveProductReviewHandler(archiveReviewService)

	// ------------------------------------------------------------
	// 🗨️ Reply to Product Review Setup
	replyRepo := repo_reply_review.NewProductReplyReviewRepository(db)
	replyService := service_reply_review.NewReplyToReviewService(service_reply_review.ReplyToReviewServiceDeps{
		Repo: replyRepo,
	})
	replyHandler := handler_reply_review.NewReplyToReviewHandler(replyService)

	// ✏️ Edit Review Reply Setup
	editReplyRepo := repo_edit_review_reply.NewProductReviewReplyEditRepository(db)
	editReplyService := edit_review_reply_service.NewEditReviewReplyService(edit_review_reply_service.EditReviewReplyServiceDeps{
		Repo: editReplyRepo,
	})
	editReplyHandler := edit_review_reply_handler.NewEditReviewReplyHandler(editReplyService)

	// 🗑️ Archive Review Reply Setup
	archiveReplyRepo := archive_reply_repo.NewProductReviewReplyArchiveRepository(db)
	archiveReplyService := archive_reply_service.NewArchiveReviewReplyService(archive_reply_service.ArchiveReviewReplyServiceDeps{
		Repo: archiveReplyRepo,
	})
	archiveReplyHandler := archive_reply_handler.NewArchiveReviewReplyHandler(archiveReplyService)

	// 🆕 Get All Product Reviews with Replies Setup
	getReviewsRepo := repo_get_reviews.NewProductGetAllReviewsRepository(db)
	getReviewsService := reviews_services.NewGetAllProductReviewsService(reviews_services.GetAllProductReviewsServiceDeps{
		Repo: getReviewsRepo,
	})
	getReviewsHandler := reviews_handlers.NewGetAllReviewsHandler(getReviewsService)

	// ------------------------------------------------------------
	// 🎥 Presigned Media Upload Endpoint Wiring (No Repo needed)

	mediaService := media_services.NewMediaService()
	mediaHandler := media_handlers.NewHandler(mediaService)

	// 🆕 Product Full Detail (Seller Side) Setup
	productFullDetailRepo := repo_product_full_detail.NewProductGetFullDetailRepository(db)
	productFullDetailService := product_full_detail_services.NewGetProductFullDetailService(
		product_full_detail_services.GetProductFullDetailServiceDeps{
			Repo: productFullDetailRepo,
		},
	)
	productFullDetailHandler := product_full_detail_handlers.NewGetProductFullDetailHandler(productFullDetailService)

	// Inside func NewRouter(...) just before r.Route("/api/media") or wherever the other seller routes are setup:
	getAllProductsBySellerRepo := repo_get_all_products_seller.NewProductGetAllGroupedRepository(db)
	getAllProductsBySellerService := get_all_products_service.NewGetAllProductsBySellerService(
		get_all_products_service.GetAllProductsBySellerServiceDeps{
			Repo: getAllProductsBySellerRepo,
		},
	)
	getAllProductsBySellerHandler := get_all_products_handler.NewGetAllProductsBySellerHandler(getAllProductsBySellerService)

	// 🆕 Update Product Category Setup
	productUpdateCategoryRepo := repo_update_category.NewProductUpdateCategoryRepository(db)
	updateProductCategoryService := update_category_services.NewUpdateProductCategoryService(update_category_services.UpdateProductCategoryServiceDeps{
		Repo: productUpdateCategoryRepo,
	})
	updateProductCategoryHandler := update_category_handlers.NewUpdateProductCategoryHandler(updateProductCategoryService)

	// ------------------------------------------------------------
	// 🌐 Public Product Detail Setup

	publicDetailRepo := public_detail_repo.NewProductGetPublicDetailRepository(db)
	publicDetailService := public_detail_service.NewGetProductPublicDetailService(public_detail_service.GetProductPublicDetailServiceDeps{
		Repo: publicDetailRepo,
	})
	publicDetailHandler := public_detail_handlers.NewGetProductPublicDetailHandler(publicDetailService)

	// ------------------------------------------------------------
	// 🧾 Checkout Session Details Endpoint Wiring
	checkoutSessionRepo := repo_checkout_session.NewCheckoutSessionDetailsRepository(db)
	checkoutSessionService := checkout_services.NewGetCheckoutSessionDetailsService(checkout_services.GetCheckoutSessionDetailsServiceDeps{
		Repo: checkoutSessionRepo,
	})
	checkoutSessionHandler := checkout_handlers.NewGetCheckoutSessionDetailsHandler(checkoutSessionService)

	// ------------------------------------------------------------
	// 🛒 Merge Guest Cart Endpoint Wiring

	mergeGuestCartRepo := merge_cart_repo.NewMergeGuestCartRepository(db)
	mergeGuestCartService := merge_cart_services.NewMergeGuestCartService(
		merge_cart_services.MergeGuestCartServiceDeps{
			Repo: mergeGuestCartRepo,
		},
	)
	mergeGuestCartHandler := merge_cart_handlers.NewMergeGuestCartHandler(mergeGuestCartService)

	// ------------------------------------------------------------
	// 🚚 Add Shipping Address to Checkout Setup
	addShippingRepo := checkout_shipping_address_add_repo.NewAddShippingAddressRepository(db)
	addShippingService := checkout_service.NewAddShippingAddressService(addShippingRepo, "YOUR_PATHAO_TOKEN", 123456) // replace with actual token + store ID
	addShippingHandler := checkout_handlers.NewAddShippingAddressHandler(addShippingService)

	// ------------------------------------------------------------
	// 🧾 Edit Shipping Address Setup

	editShippingAddressRepo := checkout_edit_repo.NewEditShippingAddressRepository(db)

	// editShippingAddressService := checkout_service.NewEditShippingAddressService(
	// 	checkout_service.NewEditShippingAddressServiceDeps{
	// 		Repo: editShippingAddressRepo,
	// 	},
	// )

	editShippingAddressService := checkout_services.NewEditShippingAddressService(
		editShippingAddressRepo,
		"your-pathao-token-here",
		123456,
	)

	editShippingAddressHandler :=
		checkout_handlers.NewEditShippingAddressHandler(editShippingAddressService)

		// ------------------------------------------------------------
	// 🧾 Review Checkout Summary Endpoint Wiring

	reviewCheckoutRepo := checkout_review_repo.NewReviewCheckoutRepository(db)

	reviewCheckoutService := checkout_service.NewReviewCheckoutService(
		reviewCheckoutRepo,
	)

	reviewCheckoutHandler :=
		checkout_handlers.NewReviewCheckoutHandler(reviewCheckoutService)

		// 🧾 Confirm COD Order Setup
	confirmOrderRepo := confirm_order_repo.NewConfirmOrderRepository(db)
	confirmOrderService := confirm_order_services.NewConfirmOrderService(confirm_order_services.ConfirmOrderServiceDeps{
		Repo: confirmOrderRepo,
	})
	confirmOrderHandler := confirm_order_handlers.NewConfirmOrderHandler(confirmOrderService)

	r.Route("/api/media", func(r chi.Router) {
		r.Use(token.AttachAccessToken) // ✅ Require valid access token

		r.Post("/presign-upload", mediaHandler.Handle)
	})

	// 📦 Routes
	r.Route("/api/auth/google", func(r chi.Router) {
		r.Post("/", googleAuthHandler.Handle)
	})

	r.Route("/api/auth/refresh", func(r chi.Router) {
		r.Post("/", refreshTokenHandler.Handle)
	})

	r.Route("/api/auth/logout", func(r chi.Router) {
		r.Use(token.AttachAccessToken) // 🛡️ Requires access token
		r.Post("/", logoutHandler.Handle)
	})

	r.Route("/api/user", func(r chi.Router) {
		// 🛡️ Requires access token
		r.Use(token.AttachAccessToken)

		r.Post("/switch-mode", switchModeHandler.Handle)
	})

	// 🌐 Public Feed/Search Routes
	r.Get("/api/feed", feedSearchHandler.HandleFeed)
	r.Get("/api/search", feedSearchHandler.HandleSearch)

	// r.Route("/api/products", func(r chi.Router) {
	// 	r.Get("/{product_id}", publicDetailHandler.Handle)
	// })

	// 🆕 Public Product Reviews Endpoint
	r.Get("/api/products/{product_id}/reviews", getReviewsHandler.Handle)
	r.Get("/api/products/{product_id}", publicDetailHandler.Handle)

	// 🌲 Public category tree route
	r.Get("/api/categories/tree", categoryTreeHandler.Handle)

	// 🛍️ Public route: Get all products by category
	r.Get("/api/category-products", fetchByCategoryHandler.Handle)

	r.Route("/api/seller", func(r chi.Router) {
		r.Use(token.AttachAccessToken)

		// 🧾 Create Seller Profile Metadata
		r.Post("/profile/metadata", sellerProfileHandler.Handle)

		// ✅ Create Product
		r.Post("/products", productHandler.Handle)

		// 🆕 Update Product Title and/or Description
		r.Put("/products/{product_id}", updateProductInfoHandler.Handle)

		// 🆕 Add Product Media (image or promo_video)
		r.Post("/products/{product_id}/media", addProductMediaHandler.Handle)

		// 🗑️ Unified media archive endpoint
		r.Delete("/api/seller/products/{product_id}/media/{media_id}", archiveMediaHandler.Handle)

		// PUT /api/seller/products/:product_id/images/:media_id/set-primary
		r.Put("/products/{product_id}/images/{media_id}/set-primary", setPrimaryImageHandler.Handle)

		// ✅ Archive product (soft delete)
		r.Put("/products/{product_id}/archive", archiveProductHandler.Handle)

		// ➕ Add Variant to Product
		r.Post("/products/{product_id}/variants", addVariantHandler.Handle)

		// ➖ Remove Variant from Product
		r.Delete("/products/{product_id}/variants/{variant_id}", removeVariantHandler.Handle)

		// 📝 Update Variant Info
		r.Put("/products/{product_id}/variants/{variant_id}/info", updateVariantInfoHandler.Handle)

		// 💵 Update Variant Retail Price
		r.Put("/products/{product_id}/variants/{variant_id}/retail-price", updateRetailPriceHandler.Handle)

		// 📦 Update Variant In-Stock Status
		r.Put("/products/{product_id}/variants/{variant_id}/in-stock", updateInStockHandler.Handle)

		// 📦 Update Variant Stock Quantity
		r.Put("/products/{product_id}/variants/{variant_id}/stock-quantity", updateStockQuantityHandler.Handle)

		// ⚖️ Update Variant Weight
		r.Put("/products/{product_id}/variants/{variant_id}/weight", updateWeightHandler.Handle)

		// 💸 Add Retail Discount to Variant
		r.Post("/products/{product_id}/variants/{variant_id}/retail-discount", addRetailDiscountHandler.Handle)

		// update retail discount for variant
		r.Put("/products/{product_id}/variants/{variant_id}/retail-discount", updateRetailDiscountHandler.Handle)

		// ❌ Remove retail discount from variant
		r.Delete("/products/{product_id}/variants/{variant_id}/retail-discount", removeRetailDiscountHandler.Handle)

		// 🏷️ Enable Wholesale Mode for Variant
		r.Post(
			"/products/{product_id}/variants/{variant_id}/wholesale-mode",
			enableWholesaleHandler.Handle,
		)

		// ✏️ Edit Wholesale Info (price / min qty)
		r.Put(
			"/products/{product_id}/variants/{variant_id}/wholesale-mode",
			editWholesaleHandler.Handle,
		)

		// ➖ Disable Wholesale Mode for Variant
		r.Delete("/products/{product_id}/variants/{variant_id}/wholesale-mode", disableWholesaleModeHandler.Handle)

		// ➕ Add Wholesale Discount to Variant
		r.Post("/products/{product_id}/variants/{variant_id}/wholesale-discount", addWholesaleDiscountHandler.Handle)

		// 🔁 Update Wholesale Discount to Variant
		r.Put("/products/{product_id}/variants/{variant_id}/wholesale-discount", updateWholesaleDiscountHandler.Handle)

		// ❌ Remove Wholesale Discount from Variant
		r.Delete("/products/{product_id}/variants/{variant_id}/wholesale-discount", removeWholesaleDiscountHandler.Handle)

		// 🆕 Fetch full detail for a seller's product
		r.Get("/products/{product_id}", productFullDetailHandler.Handle)

		// Add inside r.Route("/api/seller") block:
		r.Get("/products", getAllProductsBySellerHandler.Handle)

		// 🆕 Update Product Category
		r.Put("/products/{product_id}/category", updateProductCategoryHandler.Handle)

	})

	// r.Route("/api/cart", func(r chi.Router) {
	// 	r.Use(token.AttachAccessToken)

	// 	r.Post("/add", addToCartHandler.Handle)
	// 	r.Put("/update", updateCartQuantityHandler.Handle) // ⬅️ Add here

	// 	r.Delete("/remove/{variant_id}", removeCartHandler.Handle) // ⬅️ Remove specific item
	// 	r.Delete("/clear", clearCartHandler.Handle)                // ⬅️ Clear entire cart

	// 	r.Get("/items", getAllCartItemsHandler.Handle) // 🆕 Get All Cart Items
	// 	r.Post("/summary", cartSummaryHandler.Handle)  // 🆕 Cart Summary

	// 	r.Post("/checkout/initiate", checkoutHandler.Handle) // 🧾 Add this line

	// 	// 🆕 Merge guest cart into authenticated cart
	// 	r.Post("/merge-guest", mergeGuestCartHandler.Handle)

	// })

	// r.Route("/api/cart", func(r chi.Router) {
	// 	// ❌ Public endpoints (guests + logged in users)
	// 	r.Post("/add", addToCartHandler.Handle)
	// 	r.Put("/update", updateCartQuantityHandler.Handle)
	// 	r.Delete("/remove/{variant_id}", removeCartHandler.Handle)
	// 	r.Delete("/clear", clearCartHandler.Handle)
	// 	r.Get("/items", getAllCartItemsHandler.Handle)
	// 	r.Post("/summary", cartSummaryHandler.Handle)
	// })

	// // ✅ Auth-only endpoints (requires access token)
	// r.Route("/api/cart", func(r chi.Router) {
	// 	r.Use(token.AttachAccessToken) // ⬅️ only for the ones below

	// 	r.Post("/checkout/initiate", checkoutHandler.Handle)
	// 	r.Post("/merge-guest", mergeGuestCartHandler.Handle)
	// })

	r.Route("/api/cart", func(r chi.Router) {
		// ❌ Public endpoints
		r.Post("/add", addToCartHandler.Handle)
		r.Put("/update", updateCartQuantityHandler.Handle)
		r.Delete("/remove/{variant_id}", removeCartHandler.Handle)
		r.Delete("/clear", clearCartHandler.Handle)
		r.Get("/items", getAllCartItemsHandler.Handle)
		r.Post("/summary", cartSummaryHandler.Handle)

		// ✅ Auth-only endpoints
		r.Group(func(r chi.Router) {
			r.Use(token.AttachAccessToken)
			r.Post("/checkout/initiate", checkoutHandler.Handle)
			r.Post("/merge-guest", mergeGuestCartHandler.Handle)
		})
	})

	r.Route("/api/checkout", func(r chi.Router) {
		r.Use(token.AttachAccessToken) // Requires access token

		// 🧾 Get Checkout Session Details
		r.Get("/{checkout_session_id}", checkoutSessionHandler.Handle)

		r.Post("/{checkout_session_id}/add-shipping-address", addShippingHandler.Handle)

		// ✏️ Edit Shipping Address
		r.Put("/{checkout_session_id}/shipping-address/{shipping_address_id}/edit", editShippingAddressHandler.Handle)

		// 🧾 Review checkout summary
		r.Get("/{checkout_session_id}/review", reviewCheckoutHandler.Handle)

		// confirm order
		r.Post("/{checkout_session_id}/confirm", confirmOrderHandler.Handle)
	})

	r.Route("/api/products", func(r chi.Router) {
		r.Use(token.AttachAccessToken)

		r.Post("/{product_id}/reviews", addReviewHandler.Handle)

		// ✏️ Edit Review
		r.Put("/{product_id}/reviews/{review_id}", editReviewHandler.Handle)

		// 🗑️ Archive a Review
		r.Put("/{product_id}/reviews/{review_id}/archive", archiveReviewHandler.Handle)

		// 🗨️ Reply to Review (Seller Only)
		r.Post("/{product_id}/reviews/{review_id}/reply", replyHandler.Handle)

		// ✏️ Edit Review Reply (Seller Only)
		r.Put("/{product_id}/reviews/{review_id}/reply", editReplyHandler.Handle)

		r.Put("/{product_id}/reviews/{review_id}/reply/archive", archiveReplyHandler.Handle)

	})

	return r
}
