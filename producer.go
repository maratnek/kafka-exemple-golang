package main

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func main() {
	fmt.Println("Producer")

	brokers := []string{"kafka:9092"}
	topic := "topic1"

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 // Flush messages every 500 milliseconds

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder("Hello Kafka!"),
	}

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Message sent: topic=%s, partition=%s, offset=%s", topic, partition, offset)

}
