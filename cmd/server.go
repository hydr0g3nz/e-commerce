package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
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

	app := fiber.New()
	api := app.Group(cfg.Server.Path)
	v1 := api.Group("/v1")
	fmt.Println("server path", cfg.Server.Path)
	v1.Get("/category/:id", categoryHandler.GetCategory)
	v1.Post("/category", categoryHandler.CreateCategory)
	v1.Put("/category", categoryHandler.UpdateCategory)
	v1.Delete("/category/:id", categoryHandler.DeleteCategory)

	app.Listen("127.0.0.1:3000")

	fmt.Println("mongo connected", mongo)
}
