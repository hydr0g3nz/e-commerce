package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	handlers "github.com/hydr0g3nz/e-commerce/internal/adapters/handler"
	adapters "github.com/hydr0g3nz/e-commerce/internal/adapters/repository"
	"github.com/hydr0g3nz/e-commerce/internal/config"
	"github.com/hydr0g3nz/e-commerce/internal/core/services"
	mongoDb "github.com/hydr0g3nz/e-commerce/pkg/mongo"
)

func main() {
	cfg, err := config.LoadConfig("./config.yml")
	if err != nil {
		panic(err)
	}

	mongo := mongoDb.DBConn(cfg)
	categoryRepository := adapters.NewCategoryRepository(mongo)
	categoryService := services.NewCategoryService(categoryRepository)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	productRepository := adapters.NewProductRepository(cfg, mongo)
	productService := services.NewProductService(productRepository)
	productHandler := handlers.NewProductHandler(productService)

	app := fiber.New(fiber.Config{
		BodyLimit: 16 * 1024 * 1024,
	})
	// Middleware for logging requests
	app.Use(logger.New())

	// Middleware to recover from panics
	app.Use(recover.New())

	// Middleware for CORS (Cross-Origin Resource Sharing)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000", // Customize allowed origins
		AllowMethods: "GET,POST,PUT,DELETE",   // Customize allowed methods
	}))

	api := app.Group(cfg.Server.Path)
	v1 := api.Group("/v1")
	//category
	v1.Get("/category", categoryHandler.GetCategoryAll)
	v1.Get("/category/:id", categoryHandler.GetCategory)
	v1.Post("/category", categoryHandler.CreateCategory)
	v1.Post("/category/product", categoryHandler.AddProduct)
	v1.Put("/category", categoryHandler.UpdateCategory)
	v1.Delete("/category/:cat_id/product/:prod_id", categoryHandler.RemoveProduct)
	v1.Delete("/category/:id", categoryHandler.DeleteCategory)
	//products
	v1.Get("/product", productHandler.GetAllProducts)
	v1.Get("/product/:id", productHandler.GetProductByID)
	v1.Post("/product", productHandler.CreateProduct)
	v1.Put("/product", productHandler.UpdateProduct)
	v1.Delete("/product/:prod_id/variant/:var_id", productHandler.RemoveVariation)
	v1.Delete("/product/:prod_id", productHandler.DeleteProduct)
	v1.Post("/product/variant/:prod_id", productHandler.AddVariation)
	v1.Post("/product/image", productHandler.UploadImage)
	v1.Delete("/product/image/:filename", productHandler.DeleteImage)
	v1.Static("/images", cfg.Upload.ServerPath)
	app.Listen(fmt.Sprintf("127.0.0.1:%d", cfg.Server.Port))

}
