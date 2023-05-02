package http

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
	"wormhole/pkg/database"
	"wormhole/pkg/email"
	"wormhole/pkg/jwt"
	"wormhole/pkg/models/account"
	"wormhole/pkg/snowflake"
)

type RegisterClaims struct {
	CreatedAt int64  `json:"time"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  []byte `json:"password"`
}

func HandleRegister(ctx *fiber.Ctx) error {
	username := ctx.FormValue("username")
	emailAddress := ctx.FormValue("emailAddress")
	password := ctx.FormValue("password")

	if err := account.ValidateLogin(username, emailAddress, password); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM accounts WHERE username = $1 OR email = $2);"
	err := database.Get().QueryRow(context.Background(), query, username, emailAddress).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if account exists: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error during account lookup"})
	}
	if exists {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "username or email already exists"})
	}

	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to encrypt password: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to encrypt password"})
	}

	claims := &RegisterClaims{
		CreatedAt: time.Now().Unix(),
		Username:  username,
		Email:     emailAddress,
		Password:  encryptedPass,
	}

	token, err := jwt.MakeToken(claims)
	if err != nil {
		log.Printf("Failed to create register confirm token: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error creating confirmation token"})
	}

	err = email.Send(emailAddress, email.ConfirmRegister, map[string]any{"Name": username, "Token": token})
	if err != nil {
		log.Printf("Failed to send confirm email: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error sending confirmation email"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func HandleRegisterConfirm(ctx *fiber.Ctx) error {
	token := ctx.Query("token")
	if token == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid token"})
	}

	var claims RegisterClaims
	err := jwt.ParseToken(token, &claims)
	if err != nil {
		log.Printf("Failed to resolve token: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid token"})
	}

	if jwt.IsTokenExpired(claims.CreatedAt) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "token expired"})
	}

	emailHash := md5.Sum([]byte(claims.Email))
	pfp := fmt.Sprintf("https://www.gravatar.com/avatar/%s?d=retro", hex.EncodeToString(emailHash[:]))

	_, err = database.Get().Exec(
		context.Background(),
		"INSERT INTO accounts (id, username, email, password, profile_picture) VALUES ($1, $2, $3, $4, $5);",
		snowflake.NextId(),
		claims.Username,
		claims.Email,
		claims.Password,
		pfp)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create account"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}
