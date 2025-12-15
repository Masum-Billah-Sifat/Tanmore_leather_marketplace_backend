package http

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	dfmiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"

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

	"tanmore_backend/pkg/token"
)

func NewRouter(db *sql.DB, redisClient *redis.Client) http.Handler {
	r := chi.NewRouter()

	// 🌐 Global middlewares
	r.Use(dfmiddleware.RequestID)
	r.Use(dfmiddleware.RealIP)
	r.Use(dfmiddleware.Logger)
	r.Use(dfmiddleware.Recoverer)

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

	// 📦 Routes
	r.Route("/api/auth/google", func(r chi.Router) {
		r.Post("/", googleAuthHandler.Handle)
	})

	r.Route("/api/auth/refresh", func(r chi.Router) {
		r.Post("/", refreshTokenHandler.Handle)
	})

	r.Route("/api/user", func(r chi.Router) {
		// 🛡️ Requires access token
		r.Use(token.AttachAccessToken)

		r.Post("/switch-mode", switchModeHandler.Handle)
	})

	r.Route("/api/seller", func(r chi.Router) {
		r.Use(token.AttachAccessToken)

		// ✅ Create Product
		r.Post("/products", productHandler.Handle)

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

	})

	return r
}
