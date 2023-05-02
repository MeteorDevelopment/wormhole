package http

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
	"wormhole/pkg/database"
	"wormhole/pkg/jwt"
	"wormhole/pkg/snowflake"
)

func HandleLogin(ctx *fiber.Ctx) error {
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")

	if username == "" || password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid username or password"})
	}

	var accountId uint64
	var accountPassword []byte
	query := "SELECT id, password FROM accounts WHERE username = $1;"
	err := database.Get().QueryRow(context.Background(), query, username).Scan(&accountId, &accountPassword)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no account found"})
		}
		log.Printf("Error selecting accounts: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if bcrypt.CompareHashAndPassword(accountPassword, []byte(password)) != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid password"})
	}

	claims := &SessionClaims{
		Id:        snowflake.NextId(),
		CreatedAt: time.Now().Unix(),
		AccountId: accountId,
	}

	token, err := jwt.MakeToken(claims)
	if err != nil {
		log.Printf("Error creating token: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	_, err = database.Get().Exec(
		context.Background(),
		"INSERT INTO sessions (id, account_id) VALUES ($1, $2);",
		claims.Id,
		claims.AccountId,
	)
	if err != nil {
		log.Printf("Error storing session: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}
