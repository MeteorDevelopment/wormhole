package protocol

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/pkg/errors"
	"log"
	"sync"
	"wormhole/pkg/database"
)

type Client struct {
	account *database.Account
	conn    *websocket.Conn
}

type Handler struct {
	clients map[*websocket.Conn]*Client
	mutex   sync.Mutex
}

func NewHandler() *Handler {
	return &Handler{
		clients: make(map[*websocket.Conn]*Client),
	}
}

func (h *Handler) HandleConnection(acc *database.Account, conn *websocket.Conn) {
	client := &Client{account: acc, conn: conn}

	h.Join(client)
	defer h.Leave(client)

	for {
		messageType, message, err := client.conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(errors.Wrap(err, "error reading message from "+client.account.Username))
			}
			break
		}

		switch messageType {
		case websocket.TextMessage:
			msg, err := DecodeMessage(message, client)
			if err != nil {
				client.SendError(err, fmt.Sprintf("error decoding message from %s", client.account.Username))
				continue
			}

			h.Broadcast(msg)

		case websocket.CloseMessage:
			h.Leave(client)
			return

		default:
			client.SendError(err, fmt.Sprintf("unhandled message type %d", messageType))
			continue
		}
	}
}

func (h *Handler) Join(client *Client) {
	h.Broadcast(SystemMessage(fmt.Sprintf("%s joined the chat.", client.account.Username)))

	h.mutex.Lock()
	h.clients[client.conn] = client
	h.mutex.Unlock()
}

func (h *Handler) Leave(client *Client) {
	h.Broadcast(SystemMessage(fmt.Sprintf("%s left the chat.", client.account.Username)))

	h.mutex.Lock()
	err := client.conn.Close()
	if err != nil {
		log.Println(errors.Wrap(err, "error closing connection for "+client.account.Username))
	}
	delete(h.clients, client.conn)
	h.mutex.Unlock()
}

func (h *Handler) Broadcast(message *Message) {
	log.Println(message)

	encoded, err := message.Encode()
	if err != nil {
		log.Println(errors.Wrap(err, "error encoding message"))
		return
	}

	h.mutex.Lock()
	for connection, client := range h.clients {
		err = connection.WriteMessage(websocket.TextMessage, encoded)
		if err != nil {
			h.Leave(client)
		}
	}
	h.mutex.Unlock()
}

func (c *Client) SendError(err error, text string) {
	err = errors.Wrap(err, text)
	c.conn.WriteJSON(fiber.Map{"error": err.Error()})
	log.Println(err)
}
