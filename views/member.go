package views

import (
	"net/http"
	"go_ws/tools"
	"time"
	"github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Member struct {
	theatre *Theatre
	user string
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	room chan []byte
}

func (m *Member) SendNof() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		m.theatre.unregisterMember <- m
		m.conn.Close()
	}()
	for {
		select {
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
				w.Write(message)

				// Add queued chat messages to the current websocket message.
				n := len(m.room)
				for i := 0; i < n; i++ {
					w.Write(newline)
					w.Write(<-m.room)
				}
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

func serveRoomWs(w http.ResponseWriter, r *http.Request, theatre *Theatre)  {
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

	member := &Member{
		theatre: theatre,
		conn: conn,
		room: make(chan []byte, 256),
		user: userId,
	}
	member.theatre.registerMember <- member

	go member.SendNof()
}
