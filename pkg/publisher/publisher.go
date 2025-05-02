package publisher

import (
	"fmt"
	"net"
	"pubsub-broker/pkg/cache"
)

func StartPublisherListener(port string, connectionCache cache.Cache) {
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

		go handlePublisherConnect(conn, connectionCache)
	}
}

func handlePublisherConnect(conn net.Conn, connectionCache cache.Cache) {
	fmt.Printf("Publisher == %s == connected\n", conn.RemoteAddr())

	connectionCache.AddConnection(cache.Publisher, conn)
	subscribers := connectionCache.GetConnections(cache.Subscriber)
	numSubs := len(subscribers)

	conn.Write([]byte(fmt.Sprintf("Active subscribers: %d\n", numSubs)))

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("Publisher == %s == disconnected\n", conn.RemoteAddr())
			connectionCache.RemoveConnection(cache.Publisher, conn)
			return
		}

		subscribers = connectionCache.GetConnections(cache.Subscriber)
		for _, sub := range subscribers {
			_, err := sub.Write(buf[:n])
			if err != nil {
				fmt.Printf("Failed to send message to subscriber %s\n", sub.RemoteAddr())
			}
		}
	}
}
