package views

import (
	"net/http"
)

//func Ws(w http.ResponseWriter, r *http.Request, theatre *Theatre) {
//	room := r.FormValue("room_id")
//	if _, ok := theatre.hubs[room]; ok {
//		serveWs(theatre.hubs[room], w, r)
//	} else {
//		empty_hub := newHub(theatre)
//		empty_hub.room_id = room
//		theatre.register <- empty_hub
//		serveWs(empty_hub, w, r)
//	}
//}

func Ws(w http.ResponseWriter, r *http.Request, theatre *Theatre)  {
	ServeWs(w, r, theatre)
}
