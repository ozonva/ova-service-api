package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
)

type Producer interface {
	SendMessage(message string) (int32, int64, error)
	SendMessages(messages []string) error
}

type SyncProducer struct {
	topic    string
	producer sarama.SyncProducer
}

func NewSyncProducer(topic string, brokers []string) (*SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)

	if err != nil {
		return nil, err
	}

	ksp := SyncProducer{
		topic:    topic,
		producer: producer,
	}

	return &ksp, err
}

func (ksp SyncProducer) SendMessage(message string) (int32, int64, error) {
	if len(message) == 0 {
		return -1, -1, fmt.Errorf("empty message is not allowed")
	}

	msg := prepareMessage(ksp.topic, message)
	partition, offset, err := ksp.producer.SendMessage(msg)
	return partition, offset, err
}

func (ksp SyncProducer) SendMessages(messages []string) error {
	msgs := make([]*sarama.ProducerMessage, len(messages))

	for i, message := range messages {
		if len(message) == 0 {
			return fmt.Errorf("some of the messages are empty")
		}

		msg := prepareMessage(ksp.topic, message)
		msgs[i] = msg
	}

	err := ksp.producer.SendMessages(msgs)
	return err
}

func prepareMessage(topic string, message string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.StringEncoder(message),
	}
	return msg
}
