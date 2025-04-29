package main

import (
	"net"
	"sync"
)

var (
	subscribers []net.Conn
	publishers  []net.Conn
	mu          sync.Mutex
)

func main() {
	go StartPublisherListener(":8000")
	go StartSubscriberListener(":8001")
	select {}
}
