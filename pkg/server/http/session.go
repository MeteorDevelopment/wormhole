package http

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"log"
	"wormhole/pkg/database"
	"wormhole/pkg/jwt"
)

type SessionClaims struct {
	Id        string `json:"id"`
	CreatedAt int64  `json:"time"`
	AccountId uint64 `json:"account"`
}

func Authorized(handler fiber.Handler) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token := ctx.Get("Authorization")
		if token == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		var claims SessionClaims
		err := jwt.ParseToken(token, &claims)
		if err != nil {
			log.Printf("Failed parsing token: %v", err)
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		var accId string
		query := "SELECT id FROM sessions WHERE id = $1 AND account_id = $2;"
		err = database.Get().QueryRow(context.Background(), query, claims.Id, claims.AccountId).Scan(&accId)
		if err != nil {
			if err == pgx.ErrNoRows {
				log.Printf("No session found for %s.", claims.Id)
				return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
			}
			log.Printf("Error selecting sessions: %v", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error fetching session details"})
		}

		ctx.Locals("token", token)
		ctx.Locals("accId", accId)

		return handler(ctx)
	}
}
