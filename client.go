package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

// RateLimitedUDPClient sends packets at a limited rate
func RateLimitedUDPClient(address string, rate int, message string, controlAddress string) {
	udpConn, err := net.Dial("udp", address)
	if err != nil {
		fmt.Println("Error connecting to server (UDP):", err)
		return
	}
	defer udpConn.Close()

	tcpConn, err := net.Dial("tcp", controlAddress)
	if err != nil {
		fmt.Println("Error connecting to server (TCP):", err)
		return
	}
	defer tcpConn.Close()

	fmt.Println("Connected to control server")
	controlReader := bufio.NewReader(tcpConn)
	go func() {
		for {
			response, err := controlReader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading from control server:", err)
				return
			}
			fmt.Println("Control server message:", response)
		}
	}()

	interval := time.Second / time.Duration(rate)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for i := 0; i < 100; i++ { // Send 100 packets
		<-ticker.C
		packet := []byte(fmt.Sprintf("%s %d", message, i))
		if len(packet) > 1024 {
			fmt.Println("Error: Packet size exceeds 1024 bytes")
			return
		}
		_, err := udpConn.Write(packet)
		if err != nil {
			fmt.Println("Error sending packet:", err)
			return
		}
		fmt.Printf("Sent packet %d, size: %d bytes\n", i, len(packet))
	}
}
