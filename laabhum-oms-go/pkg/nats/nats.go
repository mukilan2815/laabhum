package nats

import (
	"github.com/nats-io/nats.go"
	"log"
)

type NatsClient struct {
	conn *nats.Conn
	ec   *nats.EncodedConn
}

func NewNatsClient(url string) *NatsClient {
	nc, _ := nats.Connect(url)
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	return &NatsClient{
		conn: nc,
		ec:   ec,
	}
}

func ConnectNATS(url string) *nats.Conn {
	conn, err := nats.Connect(url)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	return conn
}

func PublishMessage(conn *nats.Conn, subject string, message []byte) error {
	return conn.Publish(subject, message)
}
