package main

import (
	"fmt"
	"net"
	"os"
)

// UDPServer listens for UDP packets
func UDPServer(address string) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Server listening on", address)

	buffer := make([]byte, 4096)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading packet:", err)
			continue
		}
		fmt.Printf("Received packet from %s: %s\nSize: %d bytes\n", remoteAddr, string(buffer[:n]), n)
	}
}
