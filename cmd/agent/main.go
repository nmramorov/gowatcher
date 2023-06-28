package main

import (
	"fmt"

	"github.com/nmramorov/gowatcher/internal/client"
)

var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

func printBuildInfo() {
	fmt.Printf("Build version: %s\n", BuildVersion)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Build commit: %s\n", BuildCommit)
}

func main() {
	printBuildInfo()
	client := client.Client{}
	client.Run()
}
