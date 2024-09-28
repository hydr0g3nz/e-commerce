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

func (h *CategoryHandler) GetCategory(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
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
