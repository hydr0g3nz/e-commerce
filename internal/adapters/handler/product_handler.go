package handlers

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/e-commerce/internal/adapters/dto"
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
		log.Println("Error creating product:", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.service.SetProductHeroList(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.service.SetProductList([]dto.ProductListPage{}); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).SendString("Product created")
}

func (h *ProductHandler) GetAllProducts(ctx *fiber.Ctx) error {
	category := ctx.Query("category")
	var products []dto.ProductListPage
	var err error

	if category != "" {
		products, err = h.service.GetByCategory(ctx.Context(), category)
	} else {
		products, err = h.service.GetProductList(ctx.Context())
	}

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
	go h.service.SetProductList([]dto.ProductListPage{})
	go h.service.SetProductHeroList()
	return ctx.Status(fiber.StatusOK).SendString("Product updated")
}

func (h *ProductHandler) DeleteProduct(ctx *fiber.Ctx) error {
	id := ctx.Params("prod_id")
	fmt.Println("id", id)
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
	variationID, _ = url.PathUnescape(variationID)
	fmt.Println("variationID", variationID, "productID", productID)
	if err := h.service.RemoveVariation(productID, variationID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).SendString("Variation removed")
}

func (h *ProductHandler) UploadImage(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("image")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	var fp string
	if fp, err = h.service.UploadImage(ctx, file); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(fiber.Map{"filename": fp})
}

func (h *ProductHandler) DeleteImage(ctx *fiber.Ctx) error {
	filename := ctx.Params("filename")
	if filename == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "filename is required"})
	}
	if err := h.service.DeleteImage(filename); err != nil {
		if err.Error() == "file not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).SendString("Image deleted")
}

func (h *ProductHandler) GetProductHeroList(ctx *fiber.Ctx) error {
	products, err := h.service.GetProductHeroList(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(products)
}

func (h *ProductHandler) GetProductsCategoryDelegate(ctx *fiber.Ctx) error {
	products, err := h.service.GetCacheProductsCategoryDelegate(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to get products by category: %v", err),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"categories": products,
	})
}
