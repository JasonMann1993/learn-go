package main

import (
	"fmt"
	"net"
	"protocal"
)

// socket_stick/client/main.go

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println("dial failed, err", err)
		return
	}
	defer conn.Close()
	for i := 0; i < 20; i++ {
		msg := `Hello, Hello. How are you?`
		b, errorl := protocal.Encode(msg)
		if errorl != nil {
			fmt.Println(err)
		}
		conn.Write(b)
	}
}