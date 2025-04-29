package main

import (
	"fmt"
	"net"
	"slices"
)

func StartSubscriberListener(port string) {
	ln, err := net.Listen("tcp", port)

	if err != nil {
		fmt.Println("Error starting subscriber listener:", err)
		return
	}

	defer ln.Close()

	fmt.Println("Listening for Subscriber connections on port", port)

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println("Error handling subscriber connection", err)
			continue
		}

		go handleSubscriberConnect(conn)
	}
}

func handleSubscriberConnect(conn net.Conn) {
	fmt.Printf("Subscriber == %s == connected\n", conn.RemoteAddr())

	notify := false

	mu.Lock()
	wasEmpty := len(subscribers) == 0
	subscribers = append(subscribers, conn)
	mu.Unlock()

	if wasEmpty {
		notify = true
	}

	if notify {
		mu.Lock()
		for _, pub := range publishers {
			_, err := pub.Write([]byte(fmt.Sprintf("Subscriber == %s == connected\n", conn.RemoteAddr())))

			if err != nil {
				fmt.Println("Failed to inform publisher about a new connection")
			}
		}
		mu.Unlock()

	}

	// Block subscriber read to catch and inform about disconnects
	buf := make([]byte, 1)
	_, err := conn.Read(buf)

	if err != nil {
		fmt.Printf("Subscriber == %s == disconnected\n", conn.RemoteAddr())
		removeSubscriber(conn)
	}

}

func removeSubscriber(conn net.Conn) {
	for i, sub := range subscribers {
		if sub == conn {
			subscribers = slices.Delete(subscribers, i, i+1)
			break
		}
	}

	subCount := len(subscribers)

	if subCount == 0 {
		for _, pub := range publishers {
			_, err := pub.Write([]byte(fmt.Sprintf("Active subscribers: %d\n", subCount)))

			if err != nil {
				fmt.Println("Failed to notify publisher", err)
			}
		}
	}
}
