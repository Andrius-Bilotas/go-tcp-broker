BINARY_NAME=bin/tcp_broker

build:
	go build -o ${BINARY_NAME} ./cmd/pubsub-broker

run:
	go build -o ${BINARY_NAME} ./cmd/pubsub-broker
	./${BINARY_NAME}

test:
	go test -v ./cmd/pubsub-broker

clean:
	go clean
	rm ${BINARY_NAME}
