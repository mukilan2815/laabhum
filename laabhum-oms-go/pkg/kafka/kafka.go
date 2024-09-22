package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

func SetupProducer(brokers []string) sarama.SyncProducer {
    config := sarama.NewConfig()
    producer, err := sarama.NewSyncProducer(brokers, config)
    if err != nil {
        log.Fatalf("Error creating Kafka producer: %v", err)
    }
    return producer
}

func SendMessage(producer sarama.SyncProducer, topic string, message string) error {
    msg := &sarama.ProducerMessage{
        Topic: topic,
        Value: sarama.StringEncoder(message),
    }
    _, _, err := producer.SendMessage(msg)
    return err
}

func PublishOrderCreated(orderID string) {
    log.Printf("Published 'Order Created' event for order ID: %s", orderID)
}

func PublishOrderExecuted(orderID string) {
    log.Printf("Published 'Order Executed' event for order ID: %s", orderID)
}

func PublishOrderCanceled(orderID string) {
    log.Printf("Published 'Order Canceled' event for order ID: %s", orderID)
}
