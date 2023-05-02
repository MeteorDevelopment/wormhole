package http

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"log"
	"wormhole/pkg/database"
	"wormhole/pkg/jwt"
)

func HandleLogout(ctx *fiber.Ctx) error {
	all := ctx.Query("all") == "true"

	if all {
		accId := ctx.Locals("accId").(uint64)
		err := invalidateAll(ctx, accId)
		if err != nil {
			return err
		}
	} else {
		token := ctx.Locals("token").(string)
		err := invalidateToken(ctx, token)
		if err != nil {
			return err
		}
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func invalidateAll(ctx *fiber.Ctx, accId uint64) error {
	_, err := database.Get().Exec(
		context.Background(),
		"DELETE FROM sessions WHERE account_id = $1;",
		accId,
	)
	if err != nil {
		log.Printf("Error deleting all sessions: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return nil
}

func invalidateToken(ctx *fiber.Ctx, token string) error {
	var claims SessionClaims
	err := jwt.ParseToken(token, &claims)
	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if jwt.IsTokenExpired(claims.CreatedAt) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "token expired"})
	}

	_, err = database.Get().Exec(
		context.Background(),
		"DELETE FROM sessions WHERE id = $1;",
		claims.Id,
	)
	if err != nil {
		log.Printf("Error deleting session: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return nil
}
