package main

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	fmt.Println("handle connection")
	defer conn.Close()

}

func main() {
	var err error
	conn, err := net.Listen("tcp", ":5001")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := conn.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("new connection")
		go handleConnection(conn)
	}

}
