package tcp

import (
	"fmt"
	"net"
)

// Connect returns the connection of a client if connection is successfully established with server running on port
func Connect(port string) (net.Conn, error) {
	conn, err := net.Dial("tcp", "127.0.0.1:"+port)
	//checkError(err, "Connection error")
	fmt.Println("Connection established...")

	return conn, err
}