package http

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
	"wormhole/pkg/database"
	"wormhole/pkg/email"
	"wormhole/pkg/jwt"
	"wormhole/pkg/models/account"
)

type PasswordChangeClaims struct {
	CreatedAt   int64  `json:"time"`
	AccountId   uint64 `json:"account"`
	NewPassword []byte `json:"new_password"`
}

func HandleChangePassword(ctx *fiber.Ctx) error {
	accId := ctx.Locals("accId").(uint64)

	oldPassword := ctx.FormValue("old")
	newPassword := ctx.FormValue("new")

	err := account.ValidatePassword(newPassword)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var username, emailAddress string
	var oldHash []byte
	query := "SELECT username, email, password FROM accounts WHERE id = $1;"
	err = database.Get().QueryRow(context.Background(), query, accId).Scan(&username, &emailAddress, &oldHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("No account found for %d after auth.", accId)
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no account found"})
		}
		log.Printf("Error selecting accounts: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error fetching account details"})
	}

	if bcrypt.CompareHashAndPassword(oldHash, []byte(oldPassword)) != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "incorrect password"})
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "new password is invalid"})
	}

	claims := &PasswordChangeClaims{
		CreatedAt:   time.Now().Unix(),
		AccountId:   accId,
		NewPassword: newHash,
	}

	token, err := jwt.MakeToken(claims)
	if err != nil {
		log.Printf("Error creating token: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error creating confirmation token"})
	}

	err = email.Send(emailAddress, email.ConfirmChangePassword, map[string]any{"Name": username, "Token": token})
	if err != nil {
		log.Printf("Failed to send password change confirm email: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error sending confirmation email"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func HandleConfirmChangePassword(ctx *fiber.Ctx) error {
	token := ctx.Query("token")
	if token == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid token"})
	}

	var claims PasswordChangeClaims
	err := jwt.ParseToken(token, &claims)
	if err != nil {
		log.Printf("Error parsing password change token: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error parsing token"})
	}

	if jwt.IsTokenExpired(claims.CreatedAt) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "token expired"})
	}

	query := "UPDATE accounts SET password = $1 WHERE id = $2;"
	_, err = database.Get().Exec(context.Background(), query, claims.NewPassword, claims.AccountId)
	if err != nil {
		log.Printf("Error updating account password: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update password"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}
