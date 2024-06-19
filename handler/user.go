package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/promptpay-core-service/service/user"
)

type userHandler struct {
	userService user.UserService
}

func NewUserHandler(userService user.UserService) *userHandler {
	return &userHandler{
		userService: userService,
	}
}

func (h *userHandler) GetUser(c *fiber.Ctx) error {
	result, err := h.userService.GetUser()
	if err != nil {
		return handleError(c, err)
	}
	return c.Status(200).JSON(result)
}
