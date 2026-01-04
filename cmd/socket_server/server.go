package main

import (
	"fmt"
	"net"
)

func handleRequest(conn net.Conn) {
	defer conn.Close() // 處理完記得關閉這條連線

	// 在這裡寫讀取資料的邏輯...
	// 例如：
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("連線斷開或錯誤")
			fmt.Println(err)
			return
		}
		fmt.Printf("收到訊息: %s\n", string(buf[:n]))
	}
}

func main() {
	var err error
	conn, err := net.Listen("tcp", ":5000")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer conn.Close()
	fmt.Println("Listening on :5000")
	for {
		conn, err := conn.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			return
		}
		fmt.Println("收到一個新連線！")
		go handleRequest(conn)
	}

}
