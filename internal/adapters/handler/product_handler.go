package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/e-commerce/internal/core/domain"
	"github.com/hydr0g3nz/e-commerce/internal/core/services"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) CreateProduct(ctx *fiber.Ctx) error {
	product := new(domain.Product)
	if err := ctx.BodyParser(product); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.service.Create(product); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).SendString("Product created")
}

func (h *ProductHandler) GetAllProducts(ctx *fiber.Ctx) error {
	products, err := h.service.GetAll()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(products)
}

func (h *ProductHandler) GetProductByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	product, err := h.service.GetByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(product)
}

func (h *ProductHandler) UpdateProduct(ctx *fiber.Ctx) error {
	product := new(domain.Product)
	if err := ctx.BodyParser(product); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.service.Update(product); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).SendString("Product updated")
}

func (h *ProductHandler) DeleteProduct(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := h.service.Delete(id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).SendString("Product deleted")
}

func (h *ProductHandler) AddVariation(ctx *fiber.Ctx) error {
	productID := ctx.Params("prod_id")
	variation := new(domain.Variation)
	if err := ctx.BodyParser(variation); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.service.AddVariation(productID, variation); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).SendString("Variation added")
}

func (h *ProductHandler) RemoveVariation(ctx *fiber.Ctx) error {
	productID := ctx.Params("prod_id")
	variationID := ctx.Params("var_id")
	if err := h.service.RemoveVariation(productID, variationID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).SendString("Variation removed")
}
