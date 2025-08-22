package main

import (
	"fmt"
	"log"
	"market/internal/message_broker"
	"sync"
)

var address = []string{"localhost:9091", "localhost:9092", "localhost:9093"}

func main() {
	wg := &sync.WaitGroup{}
	wg.Go(func() {
		producer, err := message_broker.NewProducer(address)
		if err != nil {
			log.Fatalf("Failed to create producer: %v", err)
		}

		defer producer.Close()
		for i := 0; i < 10; i++ {
			producer.ProduceAsync(fmt.Sprintf("Message %d", i), message_broker.TopicMy)
		}
	})
	wg.Wait()
}
