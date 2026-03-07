package server

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	Player *Player
}

func (c Client) ReceiveMessages(matchmaker *Matchmaker) {
	for {
		var message WSMessage
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			expectedCodes := []int{websocket.CloseNormalClosure, websocket.CloseGoingAway}
			if websocket.IsUnexpectedCloseError(err, expectedCodes...) {
				log.Println("WebSocket error:", err)
			}
			return
		}

		switch message.Type {
		case MessageTypeJoinMatch:
			err := c.handleJoinMatch(message.Data, matchmaker)
			if err != nil {
				log.Println("Error handling join match:", err)
			}
		case MessageTypeMove:
			err := c.handleMove(message.Data, matchmaker)
			if err != nil {
				log.Println("Error handling move:", err)
			}
		case MessageTypeResign:
			err := c.handleResign(matchmaker)
			if err != nil {
				log.Println("Error handling resign:", err)
			}
		case MessageTypeOfferDraw:
			err := c.handleOfferDraw(matchmaker)
			if err != nil {
				log.Println("Error handling offer draw:", err)
			}
		case MessageTypeRespondDraw:
			err := c.handleRespondDraw(message.Data, matchmaker)
			if err != nil {
				log.Println("Error handling respond draw:", err)
			}
		default:
			log.Println("Unsupported message type:", message.Type)
		}
	}
}

func (c Client) SendMessages() {
	for message := range c.Player.SendChan {
		err := c.Conn.WriteJSON(message)
		if err != nil {
			log.Println("Error sending message:", err)
			break
		}
	}
}
