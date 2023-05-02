package socket

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/pkg/errors"
	"log"
	"sync"
	"wormhole/pkg/models/message"
	"wormhole/pkg/server/socket/protocol"
)

var (
	clients = make(map[*websocket.Conn]*Client)
	Members = make(map[string][]*websocket.Conn)
	mutex   sync.Mutex
)

func HandleUpgrade(ctx *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(ctx) {
		log.Printf("Connection from %s didn't request websocket upgrade.", ctx.IP())
		return fiber.ErrUpgradeRequired
	}

	return ctx.Next()
}

func HandleConnection(con *websocket.Conn) {
	for {
		msgType, msg, err := con.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				sendErr(con, errors.Wrap(err, "error reading message"))
			}

			closeCon(con)
			break
		}

		if msgType != websocket.TextMessage {
			sendErr(con, errors.New("unsupported message type"))
			continue
		}

		mutex.Lock()
		_, authenticated := clients[con]
		mutex.Unlock()

		packet, err := protocol.Decode(msg)
		if err != nil {
			sendErr(con, errors.Wrap(err, "error decoding packet"))
			continue
		}

		if packet.Version != protocol.Version {
			sendErr(con, errors.New("unsupported protocol version"))
			continue
		}

		if !authenticated && packet.Type != protocol.Authenticate {
			errAndClose(con, errors.New("not authenticated"))
			continue
		}

		switch packet.Type {
		case protocol.Authenticate:
			if authenticated {
				sendErr(con, errors.New("already authenticated"))
				continue
			}

			HandleAuth(con, packet)
		case protocol.Message:
			HandleMessage(con, packet)

		default:
			sendErr(con, errors.New("unknown packet type"))
		}
	}
}

func HandleAuth(con *websocket.Conn, packet *protocol.Packet) {
	//var tokenData struct {
	//	Token string
	//}
	//err := json.Unmarshal(packet.Data, &tokenData)
	//
	//accId, err := jwt.ParseToken(tokenData.Token)
	//if err != nil {
	//	errAndClose(con, errors.Wrap(err, "error validating json"))
	//	return
	//}

	//acc, err := database.GetAccount(accId)
	//if err != nil {
	//	errAndClose(con, errors.Wrap(err, "error fetching account"))
	//	return
	//}
	//
	//mutex.Lock()
	//clients[con] = &Client{
	//	account:    acc,
	//	connection: con,
	//}
	//mutex.Unlock()
}

func HandleMessage(con *websocket.Conn, packet *protocol.Packet) {
	msg, err := message.Decode(packet.Data)
	if err != nil {
		sendErr(con, errors.Wrap(err, "error decoding message"))
		return
	}

	if msg.Content == "" || msg.GroupId == 0 || msg.ChannelId == 0 {
		sendErr(con, errors.New("message is missing required fields"))
		return
	}

	msg.SenderId = clients[con].account.Id

	//if !database.IsMemberOfGroup(msg.GroupId, msg.SenderId) {
	//	sendErr(con, errors.New("unauthorized message sender"))
	//	return
	//}
	//
	//err = database.StoreMessage(msg)
	//if err != nil {
	//	log.Printf("Error storing message: %v", err)
	//}

	relayPacket := protocol.Outbound(protocol.Message, packet.Data)
	relayEncoded, err := relayPacket.Encode()
	if err != nil {
		sendErr(con, errors.Wrap(err, "error encoding message relay packet"))
		return
	}

	mutex.Lock()
	for connection, _ := range clients {
		//if !database.IsMemberOfGroup(msg.GroupId, client.account.Id) {
		//	continue
		//}

		err = connection.WriteMessage(websocket.TextMessage, relayEncoded)
		if err != nil {
			closeCon(connection)
		}
	}
	mutex.Unlock()
}

func errAndClose(c *websocket.Conn, err error) {
	sendErr(c, err)
	closeCon(c)
}

func sendErr(c *websocket.Conn, err error) {
	_ = c.WriteJSON(fiber.Map{"error": err.Error()})
	log.Println(err.Error())
}

func closeCon(c *websocket.Conn) {
	err := c.Close()
	delete(clients, c)
	if err != nil {
		log.Printf("Error closing connection: %v", err)
	}
}
