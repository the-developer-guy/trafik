package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

var mu sync.Mutex
var packetsLastSecond int64
var bytesLastSecond int64

// UDPServer listens for UDP packets
func UDPServer(address string, controlAddress string) {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		os.Exit(1)
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error starting UDP server:", err)
		os.Exit(1)
	}
	defer udpConn.Close()

	tcpListener, err := net.Listen("tcp", controlAddress)
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		os.Exit(1)
	}
	defer tcpListener.Close()

	fmt.Println("Server listening on UDP", address, "and TCP", controlAddress)

	go func() {
		for {
			tcpConn, err := tcpListener.Accept()
			if err != nil {
				fmt.Println("Error accepting TCP connection:", err)
				continue
			}
			fmt.Println("Control client connected")
			go handleControlConnection(tcpConn)
		}
	}()

	buffer := make([]byte, 2048)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			mu.Lock()
			packetsLastSecond = 0
			bytesLastSecond = 0
			mu.Unlock()
		}
	}()

	for {
		n, _, err := udpConn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading UDP packet:", err)
			continue
		}
		mu.Lock()
		packetsLastSecond++
		bytesLastSecond += int64(n)
		mu.Unlock()
	}
}

func handleControlConnection(conn net.Conn) {
	defer conn.Close()
	writer := bufio.NewWriter(conn)
	for {
		mu.Lock()
		packetsPerSecond := packetsLastSecond
		bitsPerSecond := bytesLastSecond * 8
		mu.Unlock()
		stats := fmt.Sprintf("%spps, %sbps\n", magnitude(packetsPerSecond), magnitudeWithPrecision(bitsPerSecond))

		_, err := writer.WriteString(stats)
		if err != nil {
			fmt.Println("Error writing to control connection:", err)
			return
		}
		writer.Flush()
		time.Sleep(1 * time.Second)
	}
}
