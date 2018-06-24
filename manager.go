package main

import (
	"os"
	"log"
	"go_ws/app"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Please Enter One Command: runserver or migrate")
	}
	command := os.Args[1]
	switch command {
	case "runserver":
		app.Runserver()
		break
	case "migrate":
		app.Migrate()
		break
	default:
		log.Fatalln("Please Enter One Command: runserver or migrate")
		break
	}
}