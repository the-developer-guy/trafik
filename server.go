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
var packetsLastSecond int
var bytesLastSecond int

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

	buffer := make([]byte, 1024)
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
		n, remoteAddr, err := udpConn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading UDP packet:", err)
			continue
		}
		mu.Lock()
		packetsLastSecond++
		bytesLastSecond += n
		mu.Unlock()
		fmt.Printf("Received packet from %s: %s\nSize: %d bytes\n", remoteAddr, string(buffer[:n]), n)
	}
}

func handleControlConnection(conn net.Conn) {
	defer conn.Close()
	writer := bufio.NewWriter(conn)
	for {
		mu.Lock()
		packetsPerSecond := packetsLastSecond
		bitsPerSecond := bytesLastSecond * 8
		stats := fmt.Sprintf("Packets/s: %d, Bits/s: %d\n", packetsPerSecond, bitsPerSecond)
		mu.Unlock()
		_, err := writer.WriteString(stats)
		if err != nil {
			fmt.Println("Error writing to control connection:", err)
			return
		}
		writer.Flush()
		time.Sleep(1 * time.Second)
	}
}
