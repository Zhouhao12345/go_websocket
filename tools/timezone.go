package tools

import (
	"time"
	"go_ws/config"
	"log"
)

func Now() time.Time  {
	local, err1 := time.LoadLocation(config.TIMEZONE)
	if err1 != nil {
		log.Printf("error: %v", err1)
	}
	current := time.Now()
	return current.In(local)
}
