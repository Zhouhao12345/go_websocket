package views

import (
	"net/http"
	"go_ws/tools"
	"time"
	"github.com/gorilla/websocket"
	"log"
	"bytes"
	"go_ws/models"
	"encoding/json"
	"go_ws/error_ws"
	"strconv"
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
	world *World
	mp *Map

	pos Position

	username string
	user string
	image string

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	mapEnter chan *Member
	mapLeave chan *Member
	move chan map[string]string
	test_connect chan []byte
	receive_error chan []byte
	mpInit chan *Map
}

func (m *Member) readPump() {
	defer func() {
		m.mp.unregister <- m
		m.world.unregisterMember <- m
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
				case "move":
					if m.mp.name != "None" {
						position := event["data"].(map[string]interface{})
						positionX := position["x"].(string)
						positionY := position["y"].(string)
						m.pos.x , _ = strconv.Atoi(positionX)
						m.pos.y , _ = strconv.Atoi(positionY)
						m.mp.move <- map[string]string{
							"user": m.user,
							"x": positionX,
							"y": positionY,
						}
					}
				case "enter_map":
					if m.mp.name != "None" {
						m.mp.unregister <- m
					}
					//randNum := tools.GenerateRandomNumber(100)
					//v := randNum % 2; v==0 ||
					maps := m.world.maps
					if len(maps) == 0{
						empty_map := newMap(m.world)
						empty_map.name = "Hello World"
						m.world.register <- empty_map
						m.mp = empty_map
						m.mp.register <- m
					} else {
						randNumLen := tools.GenerateRandomNumber(len(maps))
						for mp , _ := range maps {
							if randNumLen == 0 {
								m.mp = mp
								m.mp.register <- m
								break
							}
							randNumLen = randNumLen - 1
						}
					}
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
				messageArrays := make([]string, 0)
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

			case moveDate, ok := <-m.move:
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
				messageArrays := make([]map[string]string, 0)
				messageArrays = append(messageArrays, map[string]string{
					"user": moveDate["user"],
					"x": moveDate["x"],
					"y": moveDate["y"],
				})
				// Add queued chat messages to the current websocket message.
				n := len(m.move)
				for i := 0; i < n; i++ {
					moveDate = <-m.move
					messageArrays = append(messageArrays, map[string]string{
						"user": moveDate["user"],
						"x": moveDate["x"],
						"y": moveDate["y"],
					})
				}
				messageData["method"] = "move"
				messageData["data"] = messageArrays
				encoder := json.NewEncoder(w)
				encoder.SetEscapeHTML(false)
				encoder.Encode(messageData)

				if err := w.Close(); err != nil {
					return
				}
			case member, ok := <-m.mapEnter:
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
				messageArrays := make([]map[string]string, 0)
				messageArrays = append(messageArrays, map[string]string{
					"user": member.user,
					"image": member.image,
					"username": member.username,
				})
				// Add queued chat messages to the current websocket message.
				n := len(m.mapEnter)
				for i := 0; i < n; i++ {
					member = <-m.mapEnter
					messageArrays = append(messageArrays, map[string]string{
						"user": member.user,
						"image": member.image,
						"username": member.username,
					})
				}
				messageData["method"] = "mapEnter"
				messageData["data"] = messageArrays
				encoder := json.NewEncoder(w)
				encoder.SetEscapeHTML(false)
				encoder.Encode(messageData)
				if err := w.Close(); err != nil {
					return
				}
			case member, ok := <-m.mapLeave:
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
				messageArrays := make([]map[string]string, 0)
				messageArrays = append(messageArrays, map[string]string{
					"user": member.user,
					"username": member.username,
				})
				// Add queued chat messages to the current websocket message.
				n := len(m.mapLeave)
				for i := 0; i < n; i++ {
					member = <-m.mapLeave
					messageArrays = append(messageArrays, map[string]string{
						"user": member.user,
						"username": member.username,
					})
				}
				messageData["method"] = "mapLeave"
				messageData["data"] = messageArrays
				encoder := json.NewEncoder(w)
				encoder.SetEscapeHTML(false)
				encoder.Encode(messageData)
				if err := w.Close(); err != nil {
					return
				}
			case mp, ok := <-m.mpInit:
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
				messageArrays := make([]map[string]interface{}, 0)
				items := make([]string, 0)
				buildings := make([]string, 0)
				members := make([]map[string]interface{}, 0)
				for _ ,item := range mp.items {
					items = append(items, item.name)
				}
				for _ ,building := range mp.buildings {
					buildings = append(buildings, building.name)
				}
				for _ ,member := range mp.members {
					members = append(members, map[string]interface{}{
						"name": member.username,
						"id": member.user,
						"image": member.image,
						"positionX": member.pos.x,
						"positionY": member.pos.y,
					})
				}
				messageArrays = append(messageArrays, map[string]interface{}{
					"mapName": mp.name,
					"items": items,
					"members": members,
					"buildings": buildings,
				})
				// Add queued chat messages to the current websocket message.
				n := len(m.mpInit)
				for i := 0; i < n; i++ {
					mp = <-m.mpInit
					items = make([]string, 0)
					buildings = make([]string, 0)
					members = make([]map[string]interface{}, 0)
					for _ ,item := range mp.items {
						items = append(items, item.name)
					}
					for _ ,building := range mp.buildings {
						buildings = append(buildings, building.name)
					}
					for _ ,member := range mp.members {
						members = append(members, map[string]interface{}{
							"name": member.username,
							"id": member.user,
							"image": member.image,
							"positionX": member.pos.x,
							"positionY": member.pos.y,
						})
					}
					messageArrays = append(messageArrays, map[string]interface{}{
						"mapName": mp.name,
						"items": items,
						"members": members,
						"buildings": buildings,
					})
				}
				messageData["method"] = "mapInit"
				messageData["data"] = messageArrays
				encoder := json.NewEncoder(w)
				encoder.SetEscapeHTML(false)
				encoder.Encode(messageData)
				if err := w.Close(); err != nil {
					return
				}
			case error_code, ok := <-m.receive_error:
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
				messageArrays := make([]map[string]string, 0)
				messageArrays = append(messageArrays, map[string]string{
					"error": error_ws.Errormessagegenerate(string(error_code)),
					"code": string(error_code),
				})
				// Add queued chat messages to the current websocket message.
				n := len(m.receive_error)
				for i := 0; i < n; i++ {
					messageArrays = append(messageArrays, map[string]string{
						"error": error_ws.Errormessagegenerate(string(error_code)),
						"code": string(error_code),
					})
				}
				messageData["method"] = "error_received"
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

func ServeWs(w http.ResponseWriter, r *http.Request, world *World)  {
	signed, userId := tools.SingleSign(r)
	if signed == false {
		http.Error(w, "Please Sign in!", http.StatusOK)
		return
	}
	if member_existed, ok := world.members[userId];ok{
		member_existed.receive_error <- []byte("0003")
		delete(world.members, member_existed.user)
		delete(member_existed.mp.members, member_existed.user)
		member_existed.conn.Close()
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Connect Forbidden!", http.StatusForbidden)
		return
	}

	userRow, err := models.SelectQuery(
		"select username, avatar_image as image from users where id = ?", userId)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "DB ERROR", http.StatusInternalServerError)
		return
	}

	mp := newMap(world)
	member := &Member{
		world: world,
		conn: conn,
		mpInit:make(chan *Map),
		mp:mp,
		pos:Position{
			x:50,
			y:50},
		user: userId,
		username:userRow[0]["username"],
		image:userRow[0]["image"],
		test_connect:make(chan []byte,256),
		move:make(chan map[string]string),
		mapEnter:make(chan *Member),
		mapLeave:make(chan *Member),
		receive_error:make(chan []byte,256),
	}
	member.world.registerMember <- member
	go member.writePump()
	go member.readPump()
}
