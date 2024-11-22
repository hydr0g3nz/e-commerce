package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	handlers "github.com/hydr0g3nz/e-commerce/internal/adapters/handler"
	"github.com/hydr0g3nz/e-commerce/internal/adapters/middleware"
	adapters "github.com/hydr0g3nz/e-commerce/internal/adapters/repository"
	"github.com/hydr0g3nz/e-commerce/internal/config"
	"github.com/hydr0g3nz/e-commerce/internal/core/services"
	mongoDb "github.com/hydr0g3nz/e-commerce/pkg/mongo"
	rd "github.com/hydr0g3nz/e-commerce/pkg/redis"
)

func main() {
	cfg, err := config.LoadConfig("./config.yaml")
	if err != nil {
		panic(err)
	}

	redis := rd.NewRedisClient(cfg.Cache)
	mongo := mongoDb.DBConn(cfg)
	categoryRepository := adapters.NewCategoryRepository(mongo)
	categoryService := services.NewCategoryService(categoryRepository)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	productRepository := adapters.NewProductRepository(cfg, mongo, redis)
	productService := services.NewProductService(productRepository)
	productHandler := handlers.NewProductHandler(productService)

	authRepository := adapters.NewAuthRepository(mongo)
	authService := services.NewAuthService(cfg.Key.AccessToken, cfg.Key.RefreshToken, authRepository)
	authHandler := handlers.NewAuthHandler(authService)

	orderRepository := adapters.NewOrderRepository(mongo)
	orderService, err := services.NewOrderService(orderRepository, productRepository, cfg.Amqp.Url)
	if err != nil {
		panic(err)
	}
	// Init product list
	if err := productService.InitProductList(); err != nil {
		panic(err)
	}
	if err := productService.InitProductHeroList(); err != nil {
		panic(err)
	}
	// Start reservation consumer
	orderService.StartReservationConsumer()
	defer orderService.Close()
	orderHandler := handlers.NewOrderHandler(orderService)

	app := fiber.New(fiber.Config{
		BodyLimit: 16 * 1024 * 1024,
	})
	// Middleware for logging requests
	app.Use(logger.New())

	// Middleware to recover from panics
	app.Use(recover.New())

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow all origins
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	m := middleware.NewAuthMiddleware(cfg.Key.AccessToken)

	api := app.Group(cfg.Server.Path)
	// api.Use(middleware.AuthenticateJWT())
	v1 := api.Group("/v1")
	//category
	v1.Get("/category", categoryHandler.GetCategoryAll)
	v1.Get("/category/:id", categoryHandler.GetCategory)
	v1.Post("/category", m.AuthenticateJWT(), m.RequireRole("admin"), categoryHandler.CreateCategory)
	v1.Post("/category/product", m.AuthenticateJWT(), m.RequireRole("admin"), categoryHandler.AddProduct)
	v1.Put("/category", m.AuthenticateJWT(), m.RequireRole("admin"), categoryHandler.UpdateCategory)
	v1.Delete("/category/:cat_id/product/:prod_id", m.AuthenticateJWT(), m.RequireRole("admin"), categoryHandler.RemoveProduct)
	v1.Delete("/category/:id", m.AuthenticateJWT(), m.RequireRole("admin"), categoryHandler.DeleteCategory)
	//products
	v1.Get("/product", productHandler.GetAllProducts)
	v1.Get("/product-hero", productHandler.GetProductHeroList)
	v1.Get("/product/:id", productHandler.GetProductByID)
	v1.Post("/product", m.AuthenticateJWT(), m.RequireRole("admin"), productHandler.CreateProduct)
	v1.Put("/product", m.AuthenticateJWT(), m.RequireRole("admin"), productHandler.UpdateProduct)
	v1.Delete("/product/:prod_id/variant/:var_id", m.AuthenticateJWT(), m.RequireRole("admin"), productHandler.RemoveVariation)
	v1.Delete("/product/:prod_id", m.AuthenticateJWT(), m.RequireRole("admin"), productHandler.DeleteProduct)
	v1.Post("/product/variant/:prod_id", m.AuthenticateJWT(), m.RequireRole("admin"), productHandler.AddVariation)
	v1.Post("/product/image", m.AuthenticateJWT(), m.RequireRole("admin"), productHandler.UploadImage)
	v1.Delete("/product/image/:filename", m.AuthenticateJWT(), m.RequireRole("admin"), productHandler.DeleteImage)
	v1.Static("/images", cfg.Upload.ServerPath)
	//orders
	v1.Post("/order", m.AuthenticateJWT(), orderHandler.CreateOrder)
	//auth
	v1.Post("/auth/login", authHandler.Login)
	v1.Post("/auth/register", authHandler.Register)
	v1.Post("/auth/refresh", authHandler.Refresh)

	app.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))

}
