package handler

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func (h *handler) GetPriceSnapshot(c *fiber.Ctx) error {
	symbol := c.Params("symbol")
	if symbol == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Symbol required"})
	}

	snapshot, err := h.service.GetPriceSnapshot(c.UserContext(), symbol)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(snapshot)
}

func (h *handler) GetTrades(c *fiber.Ctx) error {
	symbol := c.Params("symbol")

	trades, err := h.service.GetTrades(c.UserContext(), symbol)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"symbol": symbol,
		"count":  len(trades),
		"data":   trades,
	})
}

func (h *handler) HandleWebSocket(c *websocket.Conn) {
	h.hub.Register(c)
	defer h.hub.Unregister(c)

	for {
		if _, _, err := c.ReadMessage(); err != nil {
			break
		}
	}
}
