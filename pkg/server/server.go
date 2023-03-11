package server

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"

	"wormhole/pkg/auth"
	"wormhole/pkg/database"
	"wormhole/pkg/protocol"
)

func Init() {
	app := fiber.New()
	app.Use(cors.New())

	handler := protocol.NewHandler()

	app.Post("/register", handleRegister)
	app.Post("/login", handleLogin)

	app.Use("/connect", handleSocketUpgrade)
	app.Get("/connect", websocket.New(func(c *websocket.Conn) {
		handleSocketConnection(c, handler)
	}))

	port := os.Getenv("PORT")
	log.Printf("Server listening on port %s...", port)
	log.Fatal(app.Listen(":" + port))
}

func handleRegister(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	err := auth.Register(username, password)
	if err != nil {
		log.Printf("Registration failed: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("User '%s' registered successfully", username)

	token, err := auth.Login(username, password)
	if err != nil {
		log.Printf("Login failed: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("User '%s' logged in successfully", username)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"token": token})
}

func handleLogin(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	token, err := auth.Login(username, password)
	if err != nil {
		log.Printf("Login failed: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("User '%s' logged in successfully", username)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"token": token})
}

func handleSocketUpgrade(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		log.Println("Unauthorized access attempt: Missing token")
		return fiber.ErrUnauthorized
	}

	_, err := auth.IsTokenValid(token)
	if err != nil {
		log.Printf("Token validation failed: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	if !websocket.IsWebSocketUpgrade(c) {
		log.Println("Invalid access attempt: Non-WebSocket upgrade request")
		return fiber.ErrUpgradeRequired
	}

	log.Println("WebSocket connection established")
	c.Locals("token", token)
	return c.Next()
}

func handleSocketConnection(conn *websocket.Conn, h *protocol.Handler) {
	token := conn.Locals("token").(string)
	accID, err := auth.IsTokenValid(token)
	if err != nil {
		log.Printf("Token validation failed: %v", err)
		closeConn(conn, err)
		return
	}

	acc, err := database.GetAccount(accID)
	if err != nil {
		log.Printf("Account retrieval failed: %v", err)
		closeConn(conn, err)
		return
	}

	log.Printf("Connection established for user '%s'", acc.Username)
	h.HandleConnection(acc, conn)
}

func closeConn(conn *websocket.Conn, err error) {
	err = conn.WriteJSON(fiber.Map{"error": err.Error()})
	if err != nil {
		log.Printf("Failed to send error: %v", err)
	}

	err = conn.Close()
	if err != nil {
		log.Printf("Failed to close connection: %v", err)
	}

	log.Printf("Connection closed: %v", err)
}
