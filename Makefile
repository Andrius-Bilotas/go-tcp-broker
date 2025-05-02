BINARY_NAME=bin/tcp_broker

build:
	go build -o ${BINARY_NAME} .

run:
	go build -o ${BINARY_NAME} .
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}
