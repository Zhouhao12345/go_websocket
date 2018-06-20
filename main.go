// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"net/http"
	"go_ws/views"
	"log"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	threate := views.NewTheatre()
	go threate.Run()
	http.HandleFunc("/home", views.ServeHome)
	http.HandleFunc("/login", views.ServeLogin)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		views.Ws(w,r,threate)
	})
	http.HandleFunc("/api/room/list", views.APIRoom)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
