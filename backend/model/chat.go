package model

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ClientsMap struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan []byte
	Mutex     *sync.Mutex
}
