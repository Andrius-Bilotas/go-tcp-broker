package main

import (
	"pubsub-broker/pkg/cache"
	"pubsub-broker/pkg/publisher"
	"pubsub-broker/pkg/subscriber"
)

func main() {
	connectionCache := cache.NewConnectionCache()
	go publisher.StartPublisherListener(":8000", connectionCache)
	go subscriber.StartSubscriberListener(":8001", connectionCache)
	select {}
}
