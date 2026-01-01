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

	// ✅ Seller profile creation
	seller_profile_handlers "tanmore_backend/internal/api/http/handlers/seller_profile"
	repo_seller_profile "tanmore_backend/internal/repository/seller_profile/seller_profile_metadata"
	seller_profile_services "tanmore_backend/internal/services/seller_profile"

	// 🆕 Update Product Info
	update_product_info_handlers "tanmore_backend/internal/api/http/handlers/product"
	repo_update_product_info "tanmore_backend/internal/repository/product/product_update_info"
	update_product_info_services "tanmore_backend/internal/services/product"

	cart_handlers "tanmore_backend/internal/api/http/handlers/cart"
	cart_repo "tanmore_backend/internal/repository/cart/add_to_cart"
	cart_services "tanmore_backend/internal/services/cart"

	// 🛒 Update Cart Quantity
	update_quantity_handlers "tanmore_backend/internal/api/http/handlers/cart"
	update_quantity_repo "tanmore_backend/internal/repository/cart/update_to_cart"
	update_quantity_service "tanmore_backend/internal/services/cart"

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

		// 🧾 Create Seller Profile Metadata
		r.Post("/profile/metadata", sellerProfileHandler.Handle)

		// ✅ Create Product
		r.Post("/products", productHandler.Handle)

		// 🆕 Update Product Title and/or Description
		r.Put("/products/{product_id}", updateProductInfoHandler.Handle)

	})

	r.Route("/api/cart", func(r chi.Router) {
		r.Use(token.AttachAccessToken)

		r.Post("/add", addToCartHandler.Handle)
		r.Put("/update", updateCartQuantityHandler.Handle) // ⬅️ Add here
	})

	return r
}
