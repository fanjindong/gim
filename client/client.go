package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	conn.Write([]byte("ping"))
	pong := make([]byte, 5)
	conn.Read(pong)
	fmt.Println(string(pong))
}
