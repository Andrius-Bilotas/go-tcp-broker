package main

import (
	"pubsub-broker/cache"
)

var connectionCache cache.Cache

func main() {
	connectionCache = cache.NewConnectionCache()
	go StartPublisherListener(":8000")
	go StartSubscriberListener(":8001")
	select {}
}
