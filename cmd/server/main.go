package main

import (
	"context"

	"github.com/nmramorov/gowatcher/internal/log"
	"github.com/nmramorov/gowatcher/internal/server"
)

func main() {
	ctx := context.Background()

	server := server.Server{}
	err := server.Run(ctx)
	if err != nil {
		log.ErrorLog.Printf("internal error launching server: %e", err)
	}
}
