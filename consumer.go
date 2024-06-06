package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/IBM/sarama"
)

func main() {
	fmt.Println("consumer")
	brokers := []string{"kafka:9092"}
	topic := "topic1"

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(partitions))
	fmt.Println(len(partitions))

	for _, partition := range partitions {
		go func(partition int32) {
			defer wg.Done()

			partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
			if err != nil {
				log.Fatal(err)
			}
			for msg := range partitionConsumer.Messages() {
				fmt.Printf("Received message:topic=%s, partition=%s, offset=%s, key=%s, value=%s",
					msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
		}(partition)
	}

	<-signals
	wg.Wait()
}
