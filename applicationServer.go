package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func runServer(repoUsers, repoInboundSMS, repoOutboundSMS *MongoDBRepository) {
	ln, err := net.Listen("tcp", ":25080")
	if err != nil {
		fmt.Println("Error setting up server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server is listening on port 25080...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, repoUsers, repoInboundSMS, repoOutboundSMS) // Handle each connection in a separate goroutine
	}
}

// Function to handle incoming TCP connections
func handleConnection(conn net.Conn, repoUsers, repoInboundSMS, repoOutboundSMS *MongoDBRepository) {
	defer conn.Close()

	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Received message:", message)
	if strings.Contains(message, "A new SMS arrived") {

		InboundLoadUnprocessedSMSes(repoUsers, repoInboundSMS, repoOutboundSMS)

	}

}
