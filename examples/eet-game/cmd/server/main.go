// Command server starts the EET Space Invaders game server.
package main

import (
	"log"

	"eet-game/internal/server"
)

func main() {
	cfg := server.DefaultConfig()
	if err := server.Run(cfg); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
