package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/e-commerce/internal/core/domain"
	"github.com/hydr0g3nz/e-commerce/internal/core/services"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) GetCategoryAll(ctx *fiber.Ctx) error {
	category, err := h.service.GetCategoryAll()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(category)
}
func (h *CategoryHandler) GetCategory(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	fmt.Println("id", id)
	category, err := h.service.GetCategory(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(category)
}

func (h *CategoryHandler) CreateCategory(ctx *fiber.Ctx) error {
	category := new(domain.Category)
	err := ctx.BodyParser(category)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	fmt.Println("category", category)
	err = h.service.CreateCategory(category)
	if err != nil {
		fmt.Println("Error creating category:", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusCreated).JSON(category)
}

func (h *CategoryHandler) UpdateCategory(ctx *fiber.Ctx) error {
	category := new(domain.Category)
	err := ctx.BodyParser(category)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	fmt.Println("category", category)
	err = h.service.UpdateCategory(category)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(category)
}

func (h *CategoryHandler) DeleteCategory(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	err := h.service.DeleteCategory(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).SendString("Category deleted")
}
func (h *CategoryHandler) AddProduct(ctx *fiber.Ctx) error {
	payload := new(struct {
		CategoryID string `json:"category_id"`
		ProductID  string `json:"product_id"`
	})
	if err := ctx.BodyParser(payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.service.AddProduct(payload.CategoryID, payload.ProductID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).SendString("Product added to category")
}
func (h *CategoryHandler) RemoveProduct(ctx *fiber.Ctx) error {

	if err := h.service.RemoveProduct(ctx.Params("cat_id"), ctx.Params("prod_id")); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).SendString("Product added to category")
}
