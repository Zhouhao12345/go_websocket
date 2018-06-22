// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package views

import (
	"go_ws/models"
	"strings"
	"log"
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	theatre *Theatre
	// Registered clients.
	clients map[*Client]bool
	room_id string
	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub(theatre *Theatre) *Hub {
	return &Hub{
		theatre:theatre,
		room_id: "None",
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			current := time.Now()
			messageArray := strings.SplitN(string(message), "&", 4)
			userId := messageArray[0]
			content := messageArray[3]
			messageFullByte := []byte(h.room_id+"&"+current.Format("2006-01-02 15:04:05.000000")+"&"+string(message))
			m := &models.Models{}
			err := m.InsertQuery(
				"INSERT INTO web_chatmessage ( create_uid, create_date, " +
					"update_uid, update_date, content, unread, room_id, user_id ) VALUES" +
					"("+userId+", NOW() + INTERVAL 8 HOUR , "+userId+", " +
					"NOW() + INTERVAL 8 HOUR, '"+content+"' , 1, "+h.room_id+", "+userId+")")
			if err != nil {
				log.Printf("error: %v", err)
			}
			h.theatre.wakeHub <- h
			for client := range h.clients {
				select {
					case client.send <- messageFullByte:
					default:
						close(client.send)
						delete(h.clients, client)
					}
			}
		}
	}
}
