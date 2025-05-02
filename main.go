package main

import (
	"pubsub-broker/cache"
)

func main() {
	connectionCache := cache.NewConnectionCache()
	go StartPublisherListener(":8000", connectionCache)
	go StartSubscriberListener(":8001", connectionCache)
	select {}
}
