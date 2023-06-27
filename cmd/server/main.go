package main

import (
	"context"

	"github.com/nmramorov/gowatcher/internal/server"
)

func main() {
	ctx := context.Background()

	server := server.Server{}
	server.Run(ctx)
}
