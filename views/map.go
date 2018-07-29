// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package views

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Map struct {
	world *World
	// Registered clients.
	members map[string]*Member
	// Inbound messages from the clients.
	items []*Item

	buildings []*Building

	// Register requests from the clients.
	register chan *Member

	// Unregister requests from clients.
	unregister chan *Member

	move chan map[string]string

	name string
}

func newMap(world *World) *Map {
	// Need random building and items
	mp := &Map{
		world:world,
		move:  make(chan map[string]string),
		register:   make(chan *Member),
		unregister: make(chan *Member),
		members:    make(map[string]*Member),
		name: "None",
	}
	items := make([]*Item, 0)
	buildings := make([]*Building, 0)
	for i:=0; i< 10; i++ {
		items = append(items, NewItem(mp))
		buildings = append(buildings, NewBuilding(mp))
	}
	mp.items = items
	mp.buildings = buildings
	return mp
}

func (h *Map) run() {
	for {
		select {
		case membernew := <-h.register:
			for memberID, member := range h.members {
				select {
				case member.mapEnter <- membernew:
				default:
					close(member.move)
					close(member.mapEnter)
					close(member.test_connect)
					close(member.receive_error)
					delete(h.members, memberID)
				}
			}
			h.members[membernew.user] = membernew
			membernew.mpInit <- h
		case memberLeave := <-h.unregister:
			if _, ok := h.members[memberLeave.user]; ok {
				delete(h.members, memberLeave.user)
				for memberID, member := range h.members {
					select {
					case member.mapLeave <- memberLeave:
					default:
						close(member.move)
						close(member.mapEnter)
						close(member.test_connect)
						close(member.receive_error)
						delete(h.members, memberID)
					}
				}
			}
		case moveData := <-h.move:
			for memberID, member := range h.members {
				if moveData["user"] != member.user{
					select {
					case member.move <- moveData:
					default:
						close(member.move)
						close(member.mapEnter)
						close(member.test_connect)
						close(member.receive_error)
						delete(h.members, memberID)
					}
				}
			}
		}
	}
}
