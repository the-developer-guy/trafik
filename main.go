package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [server|client] [options]")
		return
	}

	mode := os.Args[1]

	switch mode {
	case "server":
		if len(os.Args) != 3 {
			fmt.Println("Usage: go run main.go server <address>")
			return
		}
		address := os.Args[2]
		UDPServer(address)

	case "client":
		if len(os.Args) != 5 {
			fmt.Println("Usage: go run main.go client <address> <rate> <message>")
			return
		}
		address := os.Args[2]
		rate := 0
		fmt.Sscanf(os.Args[3], "%d", &rate)
		message := os.Args[4]
		RateLimitedUDPClient(address, rate, message)

	default:
		fmt.Println("Unknown mode. Use 'server' or 'client'.")
	}
}
