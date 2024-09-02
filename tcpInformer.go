package main

import (
	"fmt"
	"net"
)

func TcpInform(address string, port string, message string) error {
	conn, err := net.Dial("tcp", net.JoinHostPort(address, port))
	if err != nil {
		return fmt.Errorf("error connecting to server: %v", err)
	}
	defer conn.Close()

	_, err = fmt.Fprintln(conn, message)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	return nil
}
