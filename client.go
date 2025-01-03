package main

import (
	"fmt"
	"net"
	"time"
)

// RateLimitedUDPClient sends packets at a limited rate
func RateLimitedUDPClient(address string, rate int, message string) {
	conn, err := net.Dial("udp", address)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	interval := time.Second / time.Duration(rate)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for i := 0; i < 100; i++ { // Send 100 packets
		<-ticker.C
		_, err := conn.Write([]byte(fmt.Sprintf("%s %d", message, i)))
		if err != nil {
			fmt.Println("Error sending packet:", err)
			return
		}
		fmt.Println("Sent packet", i)
	}
}
