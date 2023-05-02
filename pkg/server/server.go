package server

import (
	"log"
	"wormhole/pkg/config"
	"wormhole/pkg/server/http"
	"wormhole/pkg/server/socket"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
)

func Init() {
	app := fiber.New()
	app.Use(cors.New())

	app.Route("/auth", func(r fiber.Router) {
		r.Post("/register", http.HandleRegister)
		r.Get("/registerConfirm", http.HandleRegisterConfirm)

		r.Post("/login", http.HandleLogin)
		r.Post("/logout", http.Authorized(http.HandleLogout))

		r.Post("/delete", http.Authorized(http.HandleDeleteAccount))
		r.Get("/deleteConfirm", http.HandleConfirmDeleteAccount)

		r.Route("/change", func(r fiber.Router) {
			r.Post("/username", http.Authorized(http.HandleChangeUsername))

			r.Post("/password", http.Authorized(http.HandleChangePassword))
			r.Get("/passwordConfirm", http.HandleConfirmChangePassword)

			r.Post("/email", http.Authorized(http.HandleChangeEmail))
			r.Get("/emailConfirm", http.HandleConfirmChangeEmail)
		})
	})

	app.Use("/chat", socket.HandleUpgrade)
	app.Get("/chat", websocket.New(func(c *websocket.Conn) {
		go socket.HandleConnection(c)
	}))

	log.Fatal(app.Listen(":" + config.Get().Port))
}
