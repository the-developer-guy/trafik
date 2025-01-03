package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "trafik"}

	var serverCmd = &cobra.Command{
		Use:   "server <udp_address> <tcp_address>",
		Short: "Run as a server",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			udpAddress := args[0]
			tcpAddress := args[1]
			UDPServer(udpAddress, tcpAddress)
		},
	}

	var clientCmd = &cobra.Command{
		Use:   "client <udp_address> <rate> <message> <tcp_address>",
		Short: "Run as a client",
		Args:  cobra.ExactArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			udpAddress := args[0]
			rate, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Invalid rate value:", err)
				os.Exit(1)
			}
			message := args[2]
			tcpAddress := args[3]
			RateLimitedUDPClient(udpAddress, rate, message, tcpAddress)
		},
	}

	rootCmd.AddCommand(serverCmd, clientCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
