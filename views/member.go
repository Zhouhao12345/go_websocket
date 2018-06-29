package views

import (
	"net/http"
	"go_ws/tools"
	"time"
	"github.com/gorilla/websocket"
	"log"
	"bytes"
	"go_ws/models"
	"strings"
	"encoding/json"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 2048
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Member struct {
	theatre *Theatre
	hub *Hub

	username string
	user string
	image string

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	room chan []byte
	send chan []byte
	test_connect chan []byte
	room_created chan map[string]string
	room_deleted chan []byte
}

func (m *Member) readPump() {
	defer func() {
		if m.hub.room_id != "None" {
			m.hub.unregister <- m
		}
		m.theatre.unregisterMember <- m
		m.conn.Close()
	}()
	m.conn.SetReadLimit(maxMessageSize)
	m.conn.SetReadDeadline(time.Now().Add(pongWait))
	m.conn.SetPongHandler(func(string) error { m.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := m.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		var event map[string]interface{}
		if err:=tools.StringToJson(string(message), &event); err==nil{
			switch event["method"] {
				case "test_connect":
					messageFullByte := []byte(m.user)
					m.test_connect <- messageFullByte
				case "send_message":
					content := event["data"].(string)
					messageFullByte := []byte(string(m.user)+"&"+m.username+"&"+m.image+"&"+content)
					m.hub.broadcast <- messageFullByte
				case "enter_room":
					if m.hub.room_id != "None" {
						m.hub.unregister <- m
					}
					room_id := event["data"].(map[string]interface{})["room_id"].(string)
					if hub, ok := m.theatre.hubs[room_id]; ok {
						m.hub = hub
						m.hub.register <- m
					} else {
						empty_hub := newHub(m.theatre)
						empty_hub.room_id = room_id
						m.theatre.register <- empty_hub
						m.hub = empty_hub
						m.hub.register <- m
					}
				case "delete_room":
					room_id := event["data"].(map[string]interface{})["room_id"].(string)
					if hub, ok := m.theatre.hubs[room_id]; ok {
						m.theatre.unregister <- hub
					}
					m.theatre.deleteHub <- []byte(room_id)
			}
		}
	}
}

func (m *Member) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		m.conn.Close()
	}()
	for {
		select {
			case message, ok:= <- m.test_connect:
				m.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok {
					// The hub closed the channel.
					m.conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}

				w, err := m.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}

				messageData := make(map[string]interface{})
				var messageArrays []string
				messageArrays = append(messageArrays, string(message))
				n := len(m.test_connect)
				for i := 0; i < n; i++ {
					messageArrays = append(messageArrays, string(<-m.test_connect))
				}

				messageData["method"] = "test_connect"
				messageData["data"] = messageArrays
				json.NewEncoder(w).Encode(messageData)

				if err := w.Close(); err != nil {
					return
				}

			case message, ok := <-m.room:
				m.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok {
					// The hub closed the channel.
					m.conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}

				w, err := m.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}

				messageArray := strings.SplitN(string(message), "&", 3)
				messageData := make(map[string]interface{})
				var messageArrays []map[string]string
				messageArrays = append(messageArrays, map[string]string{
					"rid": messageArray[0],
					"create_date": messageArray[1] + ".000000",
					"content": messageArray[2],
				})
				// Add queued chat messages to the current websocket message.
				n := len(m.room)
				for i := 0; i < n; i++ {
					messageArray := strings.SplitN(string(message), "&", 3)
					messageArrays = append(messageArrays, map[string]string{
						"rid": messageArray[0],
						"create_date": messageArray[1] + ".000000",
						"content": messageArray[2],
					})
				}
				messageData["method"] = "unread_room"
				messageData["data"] = messageArrays
				encoder := json.NewEncoder(w)
				encoder.SetEscapeHTML(false)
				encoder.Encode(messageData)

				if err := w.Close(); err != nil {
					return
				}
			case message, ok := <-m.send:
				m.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok {
					// The hub closed the channel.
					m.conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}

				w, err := m.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}

				messageArray := strings.SplitN(string(message), "&", 7)
				messageData := make(map[string]interface{})
				var messageArrays []map[string]string
				messageArrays = append(messageArrays, map[string]string{
					"mid": messageArray[0],
					"rid": messageArray[1],
					"create_date": messageArray[2] + ".000000",
					"from_uid": messageArray[3],
					"from_name": messageArray[4],
					"image": messageArray[5],
					"content": messageArray[6],
				})
				// Add queued chat messages to the current websocket message.
				n := len(m.send)
				for i := 0; i < n; i++ {
					messageArray := strings.SplitN(string(<-m.send), "&", 7)
					messageArrays = append(messageArrays, map[string]string{
						"mid": messageArray[0],
						"rid": messageArray[1],
						"create_date": messageArray[2] + ".000000",
						"from_uid": messageArray[3],
						"from_name": messageArray[4],
						"image": messageArray[5],
						"content": messageArray[6],
					})
				}
				messageData["method"] = "message_send"
				messageData["data"] = messageArrays
				encoder := json.NewEncoder(w)
				encoder.SetEscapeHTML(false)
				encoder.Encode(messageData)
				if err := w.Close(); err != nil {
					return
				}
			case room, ok := <-m.room_created:
				m.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok {
					// The hub closed the channel.
					m.conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}

				w, err := m.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}
				messageData := make(map[string]interface{})
				var messageArrays []map[string]string
				messageArrays = append(messageArrays, room)
				// Add queued chat messages to the current websocket message.
				n := len(m.room_created)
				for i := 0; i < n; i++ {
					messageArrays = append(messageArrays, <-m.room_created)
				}
				messageData["method"] = "room_created"
				messageData["data"] = messageArrays
				json.NewEncoder(w).Encode(messageData)
				if err := w.Close(); err != nil {
					return
				}
		case room_deleted, ok := <-m.room_deleted:
			m.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				m.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := m.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			messageData := make(map[string]interface{})
			var messageArrays []map[string]string
			messageArrays = append(messageArrays, map[string]string{
				"room_id": string(room_deleted),
			})
			// Add queued chat messages to the current websocket message.
			n := len(m.room_deleted)
			for i := 0; i < n; i++ {
				messageArrays = append(messageArrays, map[string]string{
					"room_id": string(<-m.room_deleted),
				})
			}
			messageData["method"] = "room_deleted"
			messageData["data"] = messageArrays
			json.NewEncoder(w).Encode(messageData)
			if err := w.Close(); err != nil {
				return
			}
			case <-ticker.C:
				m.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := m.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
	}
}

func ServeWs(w http.ResponseWriter, r *http.Request, theatre *Theatre)  {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Connect Forbidden!", http.StatusForbidden)
		return
	}
	signed, userId := tools.SingleSign(r)
	if signed == false {
		http.Error(w, "Please Sign in!", http.StatusOK)
		return
	}

	m := &models.Models{}
	userRow, err := m.SelectQuery(
		"select user.username, guser.avatar_image as image from auth_user as user " +
			"inner join web_ggacuser as guser on guser.user_ptr_id = user.id where user.id = ?", userId)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "DB ERROR", http.StatusInternalServerError)
		return
	}

	hub := newHub(theatre)
	member := &Member{
		theatre: theatre,
		conn: conn,
		hub:hub,
		room: make(chan []byte, 256),
		user: userId,
		send: make(chan []byte, 1024),
		username:userRow[0]["username"],
		image:userRow[0]["image"],
		test_connect:make(chan []byte,256),
		room_created:make(chan map[string]string),
		room_deleted:make(chan []byte,256),
	}
	member.theatre.registerMember <- member

	go member.writePump()
	go member.readPump()
}
