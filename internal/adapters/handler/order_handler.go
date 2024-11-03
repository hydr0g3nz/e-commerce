package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/e-commerce/internal/core/domain"
	"github.com/hydr0g3nz/e-commerce/internal/core/services"
)

func NewOrderHandler(service *services.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

type OrderHandler struct {
	service *services.OrderService
}

func (h *OrderHandler) CreateOrder(ctx *fiber.Ctx) error {
	order := new(domain.Order)
	err := ctx.BodyParser(order)
	if err != nil {
		fmt.Println("error", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	fmt.Println("order", order)
	if err := order.ValidateCreate(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	err = h.service.CreateOrder(ctx.Context(), order)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(order)
}
