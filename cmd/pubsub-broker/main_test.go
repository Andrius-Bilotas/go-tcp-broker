package main

import (
	"bufio"
	"net"
	"regexp"
	"testing"
	"time"
)

func TestPubSubTcpBrokerIntegrationTest(t *testing.T) {
	go main()

	// Wait for server startup
	time.Sleep(200 * time.Millisecond)

	pubConn, errPub := net.Dial("tcp", "localhost:8000")

	pubReader := bufio.NewReader(pubConn)

	notification, err := pubReader.ReadString('\n')

	if err != nil {
		t.Error("Error reading notification from server", err)
	}

	expectedNotification := "Active subscribers: 0\n"

	if notification != expectedNotification {
		t.Errorf("Expected %q, but got %q", expectedNotification, notification)
	}

	subConn, errSub := net.Dial("tcp", "localhost:8001")

	notification, err = pubReader.ReadString('\n')

	if err != nil {
		t.Error("Error reading notification from server", err)
	}

	subConnectPattern := `Subscriber == \[::1\]:\d{5} == connected`

	matched, _ := regexp.MatchString(subConnectPattern, notification)

	if !matched {
		t.Errorf("Expected notification to match %q, but got %q", subConnectPattern, notification)
	}

	defer pubConn.Close()

	if errPub != nil {
		t.Error(errPub)
	}

	if errSub != nil {
		t.Error(errSub)
	}

	subReader := bufio.NewReader(subConn)

	testMessage := "test publisher message\n"

	_, err = pubConn.Write([]byte(testMessage))

	if err != nil {
		t.Error("Failed to send message")
	}

	received, err := subReader.ReadString('\n')

	if err != nil {
		t.Error("Failed to read message")
	}

	if received != testMessage {
		t.Errorf("Expected %q but got %q", testMessage, received)
	}

	subConn.Close()

	notification, err = pubReader.ReadString('\n')

	if err != nil {
		t.Error("Error reading notification from server", err)
	}

	if notification != expectedNotification {
		t.Errorf("Expected %q, but got %q", expectedNotification, notification)
	}
}
