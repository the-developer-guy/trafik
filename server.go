package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

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
	for {
		n, remoteAddr, err := udpConn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading UDP packet:", err)
			continue
		}
		fmt.Printf("Received packet from %s: %s\nSize: %d bytes\n", remoteAddr, string(buffer[:n]), n)
	}
}

func handleControlConnection(conn net.Conn) {
	defer conn.Close()
	writer := bufio.NewWriter(conn)
	for i := 0; i < 5; i++ {
		_, err := writer.WriteString(fmt.Sprintf("Control message %d\n", i))
		if err != nil {
			fmt.Println("Error writing to control connection:", err)
			return
		}
		writer.Flush()
		time.Sleep(2 * time.Second)
	}
}
