package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math"
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

	interval := time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	buf := make([]byte, 1024)
	_, err = rand.Read(buf)
	if err != nil {
		fmt.Println("couldn't fill TX buffer")
	}

	for i := 1; i <= 20; i++ {
		<-ticker.C
		packetCount := int64(math.Pow(2, float64(i)))
		var burst int64
		fmt.Printf("sending %d packets\n", packetCount)
		for burst = 0; burst < packetCount; burst++ {
			_, err := udpConn.Write(buf)
			if err != nil {
				fmt.Println("Error sending packet:", err)
				return
			}
		}
	}
}
