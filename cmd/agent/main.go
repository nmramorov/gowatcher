package main

import (
	"fmt"

	"github.com/nmramorov/gowatcher/internal/client"
)

var (
	BuildVersion string = "N/A"
	BuildDate    string = "N/A"
	BuildCommit  string = "N/A"
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
