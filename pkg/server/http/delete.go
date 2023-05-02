package http

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"log"
	"time"
	"wormhole/pkg/database"
	"wormhole/pkg/email"
	"wormhole/pkg/jwt"
)

type DeleteAccountClaims struct {
	CreatedAt int64  `json:"time"`
	AccountId uint64 `json:"account"`
}

func HandleDeleteAccount(ctx *fiber.Ctx) error {
	accId := ctx.Locals("accId").(uint64)

	var username, emailAddress string
	query := "SELECT username, email FROM accounts WHERE id = $1;"
	err := database.Get().QueryRow(context.Background(), query, accId).Scan(&username, &emailAddress)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("No account found for %d after auth.", accId)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no account found"})
		}
		log.Printf("Error selecting accounts: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error fetching account details"})
	}

	claims := &DeleteAccountClaims{
		CreatedAt: time.Now().Unix(),
		AccountId: accId,
	}

	token, err := jwt.MakeToken(claims)
	if err != nil {
		log.Printf("Error creating deletion confirm token: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error creating token"})
	}

	err = email.Send(emailAddress, email.ConfirmDeleteAccount, map[string]any{"Name": username, "Token": token})
	if err != nil {
		log.Printf("Failed to send password change confirm email: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error sending confirmation email"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func HandleConfirmDeleteAccount(ctx *fiber.Ctx) error {
	token := ctx.Query("token")

	var claims DeleteAccountClaims
	err := jwt.ParseToken(token, &claims)
	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid token"})
	}

	query := "DELETE FROM accounts WHERE id = $1;"
	_, err = database.Get().Exec(context.Background(), query, claims.AccountId)
	if err != nil {
		log.Printf("Error deleting account: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error deleting account"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}
