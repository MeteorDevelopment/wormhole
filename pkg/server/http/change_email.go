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
	"wormhole/pkg/models/account"
)

type EmailChangeClaims struct {
	CreatedAt int64  `json:"time"`
	AccountId uint64 `json:"account"`
	NewEmail  string `json:"new_email"`
}

func HandleChangeEmail(ctx *fiber.Ctx) error {
	accId := ctx.Locals("accId").(uint64)

	oldEmail := ctx.FormValue("old")
	newEmail := ctx.FormValue("new")

	err := account.ValidateEmail(newEmail)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var username, emailAddress string
	query := "SELECT username, email FROM accounts WHERE id = $1;"
	err = database.Get().QueryRow(context.Background(), query, accId).Scan(&username, &emailAddress)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("No account found for %d after auth.", accId)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no account found"})
		}
		log.Printf("Error selecting accounts: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error fetching account details"})
	}

	if emailAddress != oldEmail {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "incorrect email"})
	}

	claims := &EmailChangeClaims{
		CreatedAt: time.Now().Unix(),
		AccountId: accId,
		NewEmail:  newEmail,
	}

	token, err := jwt.MakeToken(claims)
	if err != nil {
		log.Printf("Error creating token: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error creating confirmation token"})
	}

	err = email.Send(emailAddress, email.ConfirmChangeEmail, map[string]any{"Name": username, "Token": token, "NewEmail": newEmail})
	if err != nil {
		log.Printf("Failed to send email change confirm email: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error sending confirmation email"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func HandleConfirmChangeEmail(ctx *fiber.Ctx) error {
	token := ctx.Query("token")
	if token == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid token"})
	}

	var claims EmailChangeClaims
	err := jwt.ParseToken(token, &claims)
	if err != nil {
		log.Printf("Error parsing email change token: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error parsing token"})
	}

	if jwt.IsTokenExpired(claims.CreatedAt) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "token expired"})
	}

	query := "UPDATE accounts SET email = $1 WHERE id = $2;"
	_, err = database.Get().Exec(context.Background(), query, claims.NewEmail, claims.AccountId)
	if err != nil {
		log.Printf("Error updating account email: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update email"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}
