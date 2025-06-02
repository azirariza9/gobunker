package controller

import (
	"fmt"
	"gobunker/middleware"
	"gobunker/model"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type chatController struct {
	rg              *gin.RouterGroup
	authMidddleware middleware.AuthMiddleware
}

func (cc *chatController) Route() {
	cc.rg.GET("/chat", cc.authMidddleware.RequireToken("admin", "user"), cc.chatHandler)

}

func NewChatController(rg *gin.RouterGroup, authMiddleware middleware.AuthMiddleware) *chatController {
	go broadcaster()
	return &chatController{rg: rg, authMidddleware: authMiddleware}
}

var clientsMap = model.ClientsMap{
	Clients:   make(map[*websocket.Conn]bool),
	Broadcast: make(chan []byte),
	Mutex:     &sync.Mutex{},
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (cc *chatController) chatHandler(c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()

	clientsMap.Mutex.Lock()
	clientsMap.Clients[conn] = true
	clientsMap.Mutex.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			clientsMap.Mutex.Lock()
			delete(clientsMap.Clients, conn)
			clientsMap.Mutex.Unlock()
			break
		}
		clientsMap.Broadcast <- message
	}

}

func broadcaster() {
	for {
		// Grab the next message from the broadcast channel
		message := <-clientsMap.Broadcast

		// Send the message to all connected clients
		clientsMap.Mutex.Lock()
		for client := range clientsMap.Clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client.Close()
				delete(clientsMap.Clients, client)
			}
		}
		clientsMap.Mutex.Unlock()
	}
}
