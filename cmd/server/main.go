package main

import (
	"github.com/nmramorov/gowatcher/internal/server"
)

func main() {
	server := server.Server{}
	server.Run()
}
