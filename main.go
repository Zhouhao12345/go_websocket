// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"net/http"
	"go_ws/views"
	"log"
	"time"
	"go_ws/config"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	local, err1 := time.LoadLocation(config.TIMEZONE)
	if err1 != nil {
		log.Fatalln(err1)
	}
	threate := views.NewTheatre(local)
	go threate.Run()

	// view
	http.HandleFunc("/home", views.ServeHome)
	http.HandleFunc("/login", views.ServeLogin)

	// web socket
	http.HandleFunc("/ws_message", func(w http.ResponseWriter, r *http.Request) {
		views.Ws(w,r,threate)
	})
	http.HandleFunc("/ws_unread", func(w http.ResponseWriter, r *http.Request) {
		views.RoomWs(w,r,threate)
	})

	// api
	http.HandleFunc("/api/room/list", views.APIRoom)
	http.HandleFunc("/api/user", views.APIUser)
	http.HandleFunc("/api/room/message/list", views.APIMessage)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
