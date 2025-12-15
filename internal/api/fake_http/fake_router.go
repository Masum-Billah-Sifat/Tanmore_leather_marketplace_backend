package fake_http

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	dfmiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"

	repo_googleauth "tanmore_backend/internal/repository/google_auth"

	google_auth_handlers "tanmore_backend/internal/api/http/handlers/google_auth"
	google_auth_services "tanmore_backend/internal/services/google_auth" // same package

	// Product creation imports
	product_handlers "tanmore_backend/internal/api/http/handlers/product"
	repo_product "tanmore_backend/internal/repository/product/product_creation"
	product_services "tanmore_backend/internal/services/product"

	// Product variant addition
	variant_handlers "tanmore_backend/internal/api/http/handlers/product/product_variant"
	variant_services "tanmore_backend/internal/services/product/product_variant"

	repo_variant_update_price "tanmore_backend/internal/repository/product/product_variant/product_variant_update_price"

	// ⬇️ Add below existing import blocks

	repo_variant_add_discount "tanmore_backend/internal/repository/product/product_variant/product_variant_add_discount"

	repo_variant_update_discount "tanmore_backend/internal/repository/product/product_variant/product_variant_update_discount"

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

	// Login Handler
	googleAuthService := google_auth_services.NewGoogleAuthService(google_auth_services.GoogleAuthServiceDeps{
		Repo: googleAuthRepo,
	})
	googleAuthHandler := google_auth_handlers.NewHandler(googleAuthService)

	// 🛍️ Product Creation Setup
	productRepo := repo_product.NewProductRepository(db)
	productService := product_services.NewCreateProductService(product_services.CreateProductServiceDeps{
		Repo: productRepo,
	})
	productHandler := product_handlers.NewCreateProductHandler(productService)

	// 💵 Update Variant Retail Price Setup
	productVariantUpdatePriceRepo := repo_variant_update_price.NewProductVariantUpdatePriceRepository(db)
	updateRetailPriceService := variant_services.NewUpdateVariantRetailPriceService(variant_services.UpdateVariantRetailPriceServiceDeps{
		Repo: productVariantUpdatePriceRepo,
	})
	updateRetailPriceHandler := variant_handlers.NewUpdateVariantRetailPriceHandler(updateRetailPriceService)

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

	// 📦 Routes
	r.Route("/api/auth/google", func(r chi.Router) {
		r.Post("/", googleAuthHandler.Handle)
	})

	r.Route("/api/seller", func(r chi.Router) {
		r.Use(token.AttachAccessToken)

		// ✅ Create Product
		r.Post("/products", productHandler.Handle)

		// 💵 Update Variant Retail Price
		r.Put("/products/{product_id}/variants/{variant_id}/retail-price", updateRetailPriceHandler.Handle)

		// 💸 Add Retail Discount to Variant
		r.Post("/products/{product_id}/variants/{variant_id}/retail-discount", addRetailDiscountHandler.Handle)

		// update retail discount for variant
		r.Put("/products/{product_id}/variants/{variant_id}/retail-discount", updateRetailDiscountHandler.Handle)

	})

	return r
}
