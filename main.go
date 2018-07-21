// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"net/http"
	"go_ws/views"
	"log"
	"go_ws/config"
	"runtime"
)

var addr = flag.String("addr", config.HOSTNAME+":"+config.PORT, "http service address")

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(4)
	world := views.NewWorld()
	go world.Run()

	// view
	http.HandleFunc("/home", views.ServeHome)
	http.HandleFunc("/login", views.ServeLogin)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		views.Ws(w,r,world)
	})

	// api
	http.HandleFunc("/api/user", views.APIUser)
	http.HandleFunc("/api/user/login", views.APILogin)
	http.HandleFunc("/api/user/logout", func(w http.ResponseWriter, r *http.Request) {
		views.APILogout(w, r, world)
	})
	http.HandleFunc("/api/user/register", views.APIRegister)
	http.HandleFunc("/api/user/detail", views.APIUserDetail)
	//http.HandleFunc("/api/user/characters/list", views.APIUserCharacterList)
	//http.HandleFunc("/api/user/characters/detail", views.APIUserCharacterDetail)
	//http.HandleFunc("/api/user/friends/list", views.APIFriendList)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
