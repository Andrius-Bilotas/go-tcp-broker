package main

import (
	"fmt"
	"net"
	"pubsub-broker/cache"
)

func StartSubscriberListener(port string, connectionCache cache.Cache) {
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

		go handleSubscriberConnect(conn, connectionCache)
	}
}

func handleSubscriberConnect(conn net.Conn, connectionCache cache.Cache) {
	fmt.Printf("Subscriber == %s == connected\n", conn.RemoteAddr())

	publishers := connectionCache.GetConnections(cache.Publisher)

	connectionCache.AddConnection("subscribers", conn)

	for _, pub := range publishers {
		_, err := pub.Write([]byte(fmt.Sprintf("Subscriber == %s == connected\n", conn.RemoteAddr())))

		if err != nil {
			fmt.Println("Failed to inform publisher about a new connection")
		}
	}

	// Block subscriber read to catch and inform about disconnects
	buf := make([]byte, 1)
	_, err := conn.Read(buf)

	if err != nil {
		fmt.Printf("Subscriber == %s == disconnected\n", conn.RemoteAddr())
		subscribers := connectionCache.RemoveConnection(cache.Subscriber, conn)

		if len(subscribers) == 0 {
			publishers = connectionCache.GetConnections(cache.Publisher)

			for _, pub := range publishers {
				_, err := pub.Write([]byte(fmt.Sprintln("Active subscribers: 0")))

				if err != nil {
					fmt.Println("Failed to send message to publisher")
				}
			}
		}
	}

}
