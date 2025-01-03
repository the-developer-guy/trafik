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
		Use:   "server <address>",
		Short: "Run as a server",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			address := args[0]
			UDPServer(address)
		},
	}

	var clientCmd = &cobra.Command{
		Use:   "client <address> <rate> <message>",
		Short: "Run as a client",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			address := args[0]
			rate, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Invalid rate value:", err)
				os.Exit(1)
			}
			message := args[2]
			RateLimitedUDPClient(address, rate, message)
		},
	}

	rootCmd.AddCommand(serverCmd, clientCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
