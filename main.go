// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request)  {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cookie_age := time.Hour * 24 / time.Second
	userid_cookie:=&http.Cookie{
		Name:   "user_id",
		Value:    r.FormValue("user_id"),
		Path:     "/",
		HttpOnly: false,
		MaxAge:  int(cookie_age),
	}
	http.SetCookie(w, userid_cookie)
	http.ServeFile(w, r, "/home/zhouhao/go/src/awesomeProject1/go_ws/home.html")
}

func main() {
	flag.Parse()
	threate := newTheatre()
	go threate.run()
	http.HandleFunc("/home", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		room := r.FormValue("room_id")
		if _, ok := threate.hubs[room]; ok {
			serveWs(threate.hubs[room], w, r)
		} else {
			empty_hub := newHub()
			empty_hub.room_id = room
			threate.register <- empty_hub
			serveWs(empty_hub, w, r)
		}
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
