package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rnikrozoft/promptpay-core-service/pkg/errs"
)

func handleError(ctx *fiber.Ctx, err error) error {
	switch e := err.(type) {
	case errs.AppError:
		return ctx.Status(e.Code).JSON(e)
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(errs.NewUnexpectedError())
	}
}
