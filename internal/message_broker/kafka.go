package message_broker

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

const (
	TopicMy     = "my-topic"
	BufferSize  = 100
	WorkerCount = 10 // Количество воркеров для обработки сообщений
)

type Message struct {
	message string
	topic   string
}

type Producer struct {
	writer *kafka.Writer
	msgCh  chan Message
}

func NewProducer(address []string) (*Producer, error) {
	w := &kafka.Writer{
		Addr:         kafka.TCP(address...),
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	p := &Producer{
		writer: w,
		msgCh:  make(chan Message, BufferSize),
	}

	p.StartWorker()
	return p, nil
}

func (p *Producer) ProduceAsync(message, topic string) {
	p.msgCh <- Message{message: message, topic: topic}
}

func (p *Producer) worker() {
	for {
		select {
		case msg, ok := <-p.msgCh:
			if !ok {
				zap.L().Info("message channel closed, exiting worker")
				return
			}
			err := p.writer.WriteMessages(context.Background(), kafka.Message{
				Topic: msg.topic,
				Value: []byte(msg.message),
				Key:   nil, // Можно добавить ключ, если нужно
			})
			if err != nil {
				zap.L().Error("failed to write message", zap.Error(err))
				p.msgCh <- msg // Возвращаем сообщение обратно в канал для повторной попытки
			}
			fmt.Println("Message sent successfully:", msg.message)
		}
	}
}

func (p *Producer) StartWorker() {
	for i := 0; i < WorkerCount; i++ {
		go p.worker()
	}
}

func (p *Producer) Close() {
	close(p.msgCh)
	_ = p.writer.Close()
}
