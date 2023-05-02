package http

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"log"
	"wormhole/pkg/database"
	"wormhole/pkg/models/account"
)

func HandleChangeUsername(ctx *fiber.Ctx) error {
	accId := ctx.Locals("accId").(uint64)

	oldUsername := ctx.FormValue("old")
	newUsername := ctx.FormValue("new")

	err := account.ValidateUsername(newUsername)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var username string
	query := "SELECT username FROM accounts WHERE id = $1;"
	err = database.Get().QueryRow(context.Background(), query, accId).Scan(&username)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("No account found for %d after auth.", accId)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no account found"})
		}
		log.Printf("Error selecting accounts: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error fetching account details"})
	}

	if oldUsername != username {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "current username incorrect"})
	}

	query = "UPDATE accounts SET username = $1 WHERE id = $2;"
	_, err = database.Get().Exec(context.Background(), query, newUsername, accId)
	if err != nil {
		log.Printf("Error updating username: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error updating username"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}
