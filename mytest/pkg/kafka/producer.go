/*
  封装kafka的 同步s生产者 syncProducer
*/

package kafka

import (
	"github.com/Shopify/sarama"
)

type Producer struct {
	client sarama.SyncProducer
}

type Message struct {
	Topic string
	Value []byte
}

func NewProducer(addr []string) (producer Producer, err error) {
	p, err := sarama.NewSyncProducer(ParseAddrs(addr), nil)
	if err != nil {
		return
	}
	producer.client = p
	return
}

func (p Producer) Producer(msg Message) (err error) {
	_, _, err = p.client.SendMessage(&sarama.ProducerMessage{
		Topic: msg.Topic,
		Value: sarama.ByteEncoder(msg.Value),
	})
	return
}

func (p Producer) Close() error {
	return p.client.Close()
}
