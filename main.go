package main

import (
	"fmt"
	"net"
	"slices"
	"sync"
)

var (
	subscribers []net.Conn
	publishers  []net.Conn
	mu          sync.Mutex
)

func main() {
	go startPublisherListener(":8000")
	go startSubscriberListener(":8001")
	select {}
}

func startPublisherListener(port string) {
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

func startSubscriberListener(port string) {
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

	// Block subscriber read to catch disconnects
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
