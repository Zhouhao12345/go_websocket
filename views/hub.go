// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package views

import (
	"go_ws/models"
	"strings"
	"log"
	"go_ws/tools"
	"strconv"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	theatre *Theatre
	// Registered clients.
	members map[*Member]bool
	room_id string
	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Member

	// Unregister requests from clients.
	unregister chan *Member
}

func newHub(theatre *Theatre) *Hub {
	return &Hub{
		theatre:theatre,
		room_id: "None",
		broadcast:  make(chan []byte),
		register:   make(chan *Member),
		unregister: make(chan *Member),
		members:    make(map[*Member]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case member := <-h.register:
			h.members[member] = true
		case member := <-h.unregister:
			if _, ok := h.members[member]; ok {
				delete(h.members, member)
			}
		case message := <-h.broadcast:
			current_time := tools.Now().Format("2006-01-02 15:04:05")
			messageArray := strings.SplitN(string(message), "&", 4)
			userId := messageArray[0]
			content := messageArray[3]

			// todo fixme unread should in improve
			var unread string
			if len(h.members) > 1 {
				unread = "0"
			} else {
				unread = "1"
			}

			id, err := models.InsertQuery(
				"INSERT INTO web_chatmessage ( create_uid, create_date, " +
					"update_uid, update_date, content, unread, room_id, user_id ) VALUES" +
					"(?, ?, ?, ?, ?, ?, ?, ?)",
				userId, current_time , userId, current_time ,tools.RemoveHtmlTags(content), unread, h.room_id, userId)
			messageFullByte := []byte(strconv.FormatInt(id, 10)+"&"+h.room_id+"&"+ current_time+"&"+tools.RemoveHtmlTags(string(message)))
			if err != nil {
				log.Printf("error: %v", err)
				h.theatre.members[userId].receive_error <- []byte("0001")
				continue
			}
			h.theatre.wakeHub <- h
			for member := range h.members {
				select {
				case member.send <- messageFullByte:
				default:
					close(member.send)
					delete(h.members, member)
				}
			}
		}
	}
}
