package main

import (
	"fmt"
	"net"
)

func StartPublisherListener(port string) {
	ln, err := net.Listen("tcp", port)

	if err != nil {
		fmt.Println("Error starting publisher listener:", err)
		return
	}

	defer ln.Close()
	fmt.Println("Listening for Publisher connections on port", port)

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println("Error handling publisher connection", err)
			continue
		}

		go handlePublisherConnect(conn)
	}
}

func handlePublisherConnect(conn net.Conn) {
	fmt.Printf("Publisher == %s == connected\n", conn.RemoteAddr())
	mu.Lock()
	publishers = append(publishers, conn)
	numSubs := len(subscribers)
	mu.Unlock()

	conn.Write([]byte(fmt.Sprintf("Active subscribers: %d\n", numSubs)))

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("Publisher == %s == disconnected\n", conn.RemoteAddr())
			return
		}

		mu.Lock()
		for _, sub := range subscribers {
			_, err := sub.Write(buf[:n])
			if err != nil {
				fmt.Printf("Failed to send message to subscriber %s\n", sub.RemoteAddr())
			}
		}
		mu.Unlock()
	}
}
